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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mdales/ahocorasick"
)

type DictionaryLog struct {
	Type  string `json:"type"`
	Query string `json:"type"`
	Time  int    `json:"time"`
}

type DictionaryEntryIdentifiers struct {
	ContentMine string `json:"contentmine"`
	WikiData    string `json:"wikidata"`
}

type DictionaryEntry struct {
	Name        string                     `json:"name"`
	Term        string                     `json:"term"`
	Identifiers DictionaryEntryIdentifiers `json:"identifiers"`
}

type Dictionary struct {
	Identifier string            `json:"id"`
	log        []DictionaryLog   `json:"log"`
	Entries    []DictionaryEntry `json:"entries"`

	Matcher *ahocorasick.Matcher
}

type DictionaryMatch struct {
	Offset     int
	Entry      DictionaryEntry
	Dictionary *Dictionary
}

// Sorting interface for hits using the ahocorasick Matcher

type DictionaryMatchesByOffset []DictionaryMatch

func (a DictionaryMatchesByOffset) Len() int           { return len(a) }
func (a DictionaryMatchesByOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DictionaryMatchesByOffset) Less(i, j int) bool { return a[i].Offset < a[j].Offset }

// Parsing

func LoadDictionaryFromFile(path string) (Dictionary, error) {
	var dict Dictionary

	f, err := os.Open(path)
	if err != nil {
		return Dictionary{}, err
	}

	err = json.NewDecoder(f).Decode(&dict)
	if err != nil {
		return Dictionary{}, err
	}

	raw := make([]string, len(dict.Entries))

	for idx, entry := range dict.Entries {
		raw[idx] = entry.Term
	}

	dict.Matcher = ahocorasick.NewStringMatcher(raw)

	return dict, nil
}

func LoadDictionariesFromDirectory(directory_path string) ([]Dictionary, error) {

	files, err := ioutil.ReadDir(directory_path)
	if err != nil {
		return nil, err
	}

	res := []Dictionary{}

	for _, f := range files {
		p := path.Join(directory_path, f.Name())
		if filepath.Ext(p) == ".json" {
			dict, err := LoadDictionaryFromFile(p)
			if err != nil {
				return nil, err
			}
			res = append(res, dict)
		}
	}

	return res, nil
}

// Helper functions

func (d Dictionary) FindMatches(prose []byte) []DictionaryMatch {

	hits := d.Matcher.Match(prose)

	res := make([]DictionaryMatch, len(hits))
	for i := 0; i < len(hits); i++ {
		hit := hits[i]
		res[i] = DictionaryMatch{
			Offset:     hit.Position,
			Entry:      d.Entries[hit.Key],
			Dictionary: &d,
		}
	}

	return res
}
