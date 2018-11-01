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

// Encoding of structures in json comes from data schema found here:
//      https://sciencesource.wmflabs.org/wiki/Data_schema
type Annotation struct {
	// P0-P1 ??
	WikiDataItemCode string  `json:"P2,omitempty"`
	InstanceOf       string `json:"P3,omitempty"`
	SubclassOf       string `json:"P4,omitempty"`
	// P5 ??
	PrecedingAnchorPoint      *string `json:"P6,omitempty"`
	FollowingAnchorPoint      *string `json:"P7,omitempty"`
	DistanceToPreceding       int     `json:"P8,omitempty"`
	DistanceToFollowing       int     `json:"P9,omitempty"`
	CharacterNumber           int     `json:"P10,omitempty"`
	ArticleTextTitle          string `json:"P11,omitempty"`
	AnchorPoint               string `json:"P12,omitempty"`
	PrecedingPhrase           string `json:"P13,omitempty"`
	FollowingPhrase           string `json:"P14,omitempty"`
	TermFound                 string `json:"P15,omitempty"`
	DictionaryName            string `json:"P16,omitempty"`
	PublicationName           string `json:"P17,omitempty"`
	LengthOfTermFound         int     `json:"P18,omitempty"`
	BasedOn                   string `json:"P19,omitempty"`
	ScienceSourceArticleTitle string `json:"P20,omitempty"`
	// P21 ??
	TimeCode string `json:"P22,omitempty"`
	// P23 ??
	Anchors      string `json:"P24,omitempty"`
	FormatterURL string `json:"P25,omitempty"`
}
