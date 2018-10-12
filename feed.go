// Copyright ContentMine Ltd 2018

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const PMCAPIURL string = "https://www.ebi.ac.uk/europepmc/webservices/rest"

type Header struct {
	Vars []string `json:"vars"`
}

type DataValue struct {
	DataType *string `json:"datatype"`
	Type     string  `json:"type"`
	Value    string  `json:"value"`
	Language *string `json:"xml:lang"`
}

type Paper struct {
	Date             DataValue `json:"date"`
	Item             DataValue `json:"item"`
	ItemLabel        DataValue `json:"itemLabel"`
	JournalLabel     DataValue `json:"journalLabel"`
	LicenseLabel     DataValue `json:"licenseLabel"`
	MainSubjectLabel DataValue `json:"mainsubjectLabel"`
	PMCID            DataValue `json:"pmcid"`
	Title            DataValue `json:"title"`
}

type Results struct {
	Papers []Paper `json:"bindings"`
}

type PaperFeed struct {
	Header  Header  `json:"head"`
	Results Results `json:"results"`
}

// Parsing

func LoadFeedFromFile(path string) (PaperFeed, error) {
	var feed PaperFeed

	f, err := os.Open(path)
	if err != nil {
		return PaperFeed{}, err
	}

	err = json.NewDecoder(f).Decode(&feed)
	return feed, err
}

// Convenience functions

func (paper Paper) ID() string {
	return paper.PMCID.Value
}

func (paper Paper) String() string {
	return fmt.Sprintf("<Paper %s: %s>", paper.ID(), paper.Title.Value)
}

func (paper Paper) FullTextURL() string {
	return fmt.Sprintf("%s/PMC%s/fullTextXML", PMCAPIURL, paper.ID())
}

func (paper Paper) SupplementaryFilesURL() string {
	return fmt.Sprintf("%s/PMC%s/fullTextXML", PMCAPIURL, paper.ID())

}
