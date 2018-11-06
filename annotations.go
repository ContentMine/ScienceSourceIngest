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
	"reflect"
)

// Encoding of structures in json comes from data schema found here:
//      https://sciencesource.wmflabs.org/wiki/Data_schema
type Annotation struct {
	// P0-P1 ??
	WikiDataItemCode string `json:"P2,omitempty"`
	InstanceOf       string `json:"P3,omitempty"`
	SubclassOf       string `json:"P4,omitempty"`
	// P5 ??
	PrecedingAnchorPoint      *string `json:"P6,omitempty"`
	FollowingAnchorPoint      *string `json:"P7,omitempty"`
	DistanceToPreceding       int     `json:"P8,omitempty"`
	DistanceToFollowing       int     `json:"P9,omitempty"`
	CharacterNumber           int     `json:"P10,omitempty"`
	ArticleTextTitle          string  `json:"P11,omitempty"`
	AnchorPoint               string  `json:"P12,omitempty"`
	PrecedingPhrase           string  `json:"P13,omitempty"`
	FollowingPhrase           string  `json:"P14,omitempty"`
	TermFound                 string  `json:"P15,omitempty"`
	DictionaryName            string  `json:"P16,omitempty"`
	PublicationName           string  `json:"P17,omitempty"`
	LengthOfTermFound         int     `json:"P18,omitempty"`
	BasedOn                   string  `json:"P19,omitempty"`
	ScienceSourceArticleTitle string  `json:"P20,omitempty"`
	// P21 ??
	TimeCode string `json:"P22,omitempty"`
	// P23 ??
	Anchors      string `json:"P24,omitempty"`
	FormatterURL string `json:"P25,omitempty"`
}

type ItemType string

type AnnotationInstance struct {
	Item              ItemType `item:"annotation"`
	InstanceOf        string   `json:"P3" property:"instance of"`
	BasedOn           string   `json:"P19,omitempty" property:"based on"` // Ref to anchor point
	TermFound         string   `json:"P15,omitempty" property:"term found"`
	WikiDataItemCode  string   `json:"P2,omitempty" property:"Wikidata item code"`
	DictionaryName    string   `json:"P16,omitempty" property:"dictionary name"`
	LengthOfTermFound int      `json:"P18,omitempty" property:"length of term found"`
	TimeCode          string   `json:"P22,omitempty" property:"time code1"`
}

type AnchorPoint struct {
	Item                      ItemType `item:"anchor point"`
	InstanceOf                string   `json:"P3"`
	PrecedingAnchorPoint      *string  `json:"P6,omitempty" property:"preceding anchor point"` // Ref to anchor point/article
	FollowingAnchorPoint      *string  `json:"P7,omitempty" property:"following anchor point"` // Ref to anchor point/terminus
	DistanceToPreceding       int      `json:"P8,omitempty" property:"distance to preceding"`
	DistanceToFollowing       int      `json:"P9,omitempty" property:"distance to following"`
	Anchors                   string   `json:"P24,omitempty" property:"anchors"`
	AnchorPoint               string   `json:"P12,omitempty" property:"anchor point in"` // Ref to article
	PrecedingPhrase           string   `json:"P13,omitempty" property:"preceding phrase"`
	FollowingPhrase           string   `json:"P14,omitempty" property:"following phrase"`
	CharacterNumber           int      `json:"P10,omitempty" property:"character number"`
	TimeCode                  string   `json:"P22,omitempty" property:"time code1"`
	ScienceSourceArticleTitle string   `json:"P20,omitempty" property:"ScienceSource article title"`
}

type Article struct {
	Item                      ItemType `item:"article"`
	InstanceOf                string   `json:"P3" property:"instance of"`
	WikiDataItemCode          string   `json:"P2,omitempty" property:"Wikidata item code"`
	ArticleTextTitle          string   `json:"P11,omitempty" property:"article text title"`
	PublicationDate           string   `json:"P17,omitempty" property:"publication date"`
	TimeCode                  string   `json:"P22,omitempty" property:"time code1"`
	ScienceSourceArticleTitle string   `json:"P20,omitempty" property:"ScienceSource article title"`
	FollowingAnchorPoint      *string  `json:"P7,omitempty" property:"following anchor point"`
	CharacterNumber           int      `json:"P10,omitempty" property:"character number"` // always 0
	PrecedingPhrase           string   `json:"P13,omitempty" property:"preceding phrase"`
	FollowingPhrase           string   `json:"P14,omitempty" property:"following phrase"`
}

// terminus needs looking up too


func getValuesForTags(tagname string) []string {
	tagset := make(map[string]bool, 0)

	types := [3]reflect.Type{
		reflect.TypeOf(AnnotationInstance{}),
		reflect.TypeOf(AnchorPoint{}),
		reflect.TypeOf(Article{}),
	}

	for _, t := range types {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get(tagname)
			if len(tag) > 0 {
				tagset[tag] = true
			}
		}
	}

	list := []string{}
	for k := range tagset {
		list = append(list, k)
	}

	return list
}

func GetPropertyAndItemLists() ([]string, []string) {
    return getValuesForTags("property"), getValuesForTags("item")
}
