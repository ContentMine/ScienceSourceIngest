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
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	europmc "github.com/ContentMine/go-europmc"
)

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
	return europmc.FullTextURL(paper.ID())
}

func (paper Paper) SupplementaryFilesURL() string {
	return europmc.SupplementaryFilesURL(paper.ID())
}

func (paper Paper) WikiDataID() string {

	parts := strings.Split(paper.Item.Value, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

func (paper Paper) PublicationDate() (time.Time, error) {
	return time.Parse(time.RFC3339, paper.Date.Value)
}
