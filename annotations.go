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
	WikiDataItemCode string `json:"P2"`
	InstanceOf       string `json:"P3"`
	SubclassOf       string `json:"P4"`
	// P5 ??
	PrecedingAnchorPoint      string `json:"P6"`
	FollowingAnchorPoint      string `json:"P7"`
	DistanceToPreceding       int    `json:"P8"`
	DistanceToFollowing       int    `json:"P9"`
	CharacterNumber           int    `json:"P10"`
	ArticleTextTitle          string `json:"P11"`
	AnchorPoint               string `json:"P12"`
	PreceedingPhrase          string `json:"P13"`
	FollowingPhrase           string `json:"P14"`
	TermFound                 string `json:"P15"`
	DictionaryName            string `json:"P16"`
	PublicationName           string `json:"P17"`
	LengthOfTermFound         int    `json:"P18"`
	BasedOn                   string `json:"P19"`
	ScienceSourceArticleTitle string `json:"P20"`
	// P21 ??
	TimeCode string `json:"P22"`
	// P23 ??
	Anchors      string `json:"P24"`
	FormatterURL string `json:"P25"`
}
