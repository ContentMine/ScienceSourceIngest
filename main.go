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
	flag.StringVar(&feed_path, "feed", "", "JSON feed of papers, required")
	flag.StringVar(&target_path, "output", ".", "Directory to store the results, required")
	flag.Parse()

	log.Printf("Feed to parse: %s", feed_path)

	feed, err := LoadFeedFromFile(feed_path)
	if err != nil {
		panic(err)
	}

	log.Printf("We have %d papers to process", len(feed.Results.Papers))

    // Here I use a traditional wait group to wait for everyone to be done,
    // and I use a channel to control the number of concurrent operations allowed.
    // In theory I can use the channel also to wait at the end, but it's not as
    // easy to read the code here, so I've chosen to use both mechanisms for
    // the sake of code clarity
    var wg sync.WaitGroup
	sem := make(chan bool, concurrencyLimit)
	for _, paper := range feed.Results.Papers {
		sem <- true
		wg.Add(1)
		go func() {
            defer func() {
                wg.Done()
                <-sem
            }()
            log.Printf("Process paper %s", paper.ID())

			var processor = PaperProcessor{Paper: paper, TargetDirectory: target_path}
			err := processor.ProcessPaper()
			if err != nil {
				log.Printf("Failed to process paper %s: %v", paper.ID(), err)
			}
		}()

		// just process one whilst testing...
		break
	}
	wg.Wait()
}
