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

	var wg sync.WaitGroup

	for _, paper := range feed.Results.Papers {
		wg.Add(1)
		go func() {
			defer wg.Done()

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
