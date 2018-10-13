// Copyright ContentMine Ltd 2018

package main

import (
	"flag"
	//    "fmt"
	"log"
	"sync"
)

// This will be set by the build script to something meaningful
var Version string

// We could fire off 100 requests at once, but that's not being nice to
// either the local machine or PMC's API, so we limite the number of
// concurrent paper requests here
const concurrencyLimit int = 5

func main() {

	var feed_path string
	var target_path string
	var dictionaries_path string
	flag.StringVar(&feed_path, "feed", "", "JSON feed of papers, required")
	flag.StringVar(&target_path, "output", ".", "Directory to store the results, required")
	flag.StringVar(&dictionaries_path, "dictionaries", "", "Directory of dictionaries to load.")
	flag.Parse()

	log.Printf("Feed to parse: %s", feed_path)

	feed, err := LoadFeedFromFile(feed_path)
	if err != nil {
		panic(err)
	}

	// the SPARQL seems to have duplicates in, so let's check
	library := make(map[string]Paper)
	for _, paper := range feed.Results.Papers[0:10] {
		if _, prs := library[paper.ID()]; prs == true {
			log.Printf("Found a duplicate paper: %v", paper.ID())
		} else {
			library[paper.ID()] = paper
		}
	}
	log.Printf("We have %d papers to process", len(library))

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
			err := processor.ProcessPaper()
			if err != nil {
				log.Printf("Failed to process paper %s: %v", to_process.ID(), err)
			}
		}()
	}
	wg.Wait()

	// Load the dictionaries of terms we want to create annotations for
	dictionaries, err := LoadDictionariesFromDirectory(dictionaries_path)
	if err != nil {
		panic(err)
	}
	log.Printf("We have loaded %d dictionaries", len(dictionaries))
	for _, dict := range dictionaries {
		log.Printf("Dict %s has %d entries", dict.Identifier, len(dict.Entries))
	}
}
