// Copyright ContentMine Ltd 2018

package main

import (
    "flag"
//    "fmt"
    "log"
)

// This will be set by the build script to something meaningful
var Version string


func main() {

	var feed_path string
	flag.StringVar(&feed_path, "feed", "", "JSON feed of papers, required")
	flag.Parse()

    log.Printf("Feed to parse: %s", feed_path)

    feed, err := LoadFeedFromFile(feed_path)
    if err != nil {
        panic(err)
    }

    log.Printf("We have %d papers to process", len(feed.Results.Papers))
}
