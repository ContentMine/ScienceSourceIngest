//   Copyright 2018 Content Mine Ltd
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"flag"
	//    "fmt"
	"log"
	"sync"

	"github.com/mdales/wikibase"
)

// These will be set by the build script to something meaningful
var Remote string
var Version string

// We could fire off 100 requests at once, but that's not being nice to
// either the local machine or PMC's API, so we limite the number of
// concurrent paper requests here
const concurrencyLimit int = 5

func main() {

	var feed_path string
	var target_path string
	var dictionaries_path string
	var url_base string
	var oauth_tokens_path string
	flag.StringVar(&feed_path, "feed", "", "JSON feed of papers, required")
	flag.StringVar(&target_path, "output", ".", "Directory to store the results, required")
	flag.StringVar(&dictionaries_path, "dictionaries", "", "Directory of dictionaries to load.")
	flag.StringVar(&url_base, "urlbase", "http://localhost:8181", "Base URL for science source.")
	flag.StringVar(&oauth_tokens_path, "oauth", "oauth.json", "JSON file with oauth credentials in.")
	flag.Parse()

	log.Printf("Feed to parse: %s", feed_path)

	feed, err := LoadFeedFromFile(feed_path)
	if err != nil {
		panic(err)
	}

	// the SPARQL seems to have duplicates in, so let's check
	library := make(map[string]Paper)
	for _, paper := range feed.Results.Papers[0:1] {
		if _, prs := library[paper.ID()]; prs == true {
			log.Printf("Found a duplicate paper: %v", paper.ID())
		} else {
			library[paper.ID()] = paper
		}
	}
	log.Printf("We have %d papers to process", len(library))

	// Load the dictionaries of terms we want to create annotations for
	dictionaries, err := LoadDictionariesFromDirectory(dictionaries_path)
	if err != nil {
		panic(err)
	}
	log.Printf("We have loaded %d dictionaries", len(dictionaries))
	for _, dict := range dictionaries {
		log.Printf("Dict %s has %d entries", dict.Identifier, len(dict.Entries))
	}

	// Connect to Science Source instance and get any information we need
	oauthInfo, load_err := wikibase.LoadOauthInformation(oauth_tokens_path)
	if load_err != nil {
		panic(load_err)
	}
	sciSourceClient := NewScienceSourceClient(oauthInfo, url_base)
	err = sciSourceClient.GetPropertyAndItemConfigurationFromServer()
	if err != nil {
		panic(err)
	}

	// Here I use a traditional wait group to wait for everyone to be done,
	// and I use a channel to control the number of concurrent operations allowed.
	// In theory I can use the channel also to wait at the end, but it's not as
	// easy to read the code here, so I've chosen to use both mechanisms for
	// the sake of code clarity
	var wg sync.WaitGroup
	sem := make(chan bool, concurrencyLimit)
	for _, paper := range library {
		to_process := paper
		sem <- true
		wg.Add(1)

		go func() {
			defer func() {
				wg.Done()
				<-sem
			}()
			log.Printf("Process paper %s", to_process.ID())

			var processor = PaperProcessor{Paper: to_process, TargetDirectory: target_path}
			err := processor.ProcessPaper(dictionaries, sciSourceClient)
			if err != nil {
				log.Printf("Failed to process paper %s: %v", to_process.ID(), err)
			}
		}()
	}
	wg.Wait()
}
