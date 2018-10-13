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
}

// Parsing

func LoadDictionaryFromFile(path string) (Dictionary, error) {
	var dict Dictionary

	f, err := os.Open(path)
	if err != nil {
		return Dictionary{}, err
	}

	err = json.NewDecoder(f).Decode(&dict)
	return dict, err
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
