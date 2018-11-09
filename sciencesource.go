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

type ItemType string

type ScienceSourceAnnotation struct {
    // Exists purely to let us look up the item ID on sci source
	Item              ItemType `item:"annotation"`

	// These fields we know beforehand
	TermFound         string   `json:"term" property:"term found"`
	LengthOfTermFound int      `json:"length" property:"length of term found"`
	WikiDataItemCode  string   `json:"wikidata" property:"Wikidata item code"`
	DictionaryName    string   `json:"dictionary" property:"dictionary name"`
	TimeCode          string   `json:"time" property:"time code1"`

    // These fields we only know from the science source instance
	InstanceOf        string   `json:"instance_of" property:"instance of"`

	// Used to let us look the item up later
	ScienceSourceItemID   string `json:"id"`
}

type ScienceSourceAnchorPoint struct {
    // Exists purely to let us look up the item ID on sci source
	Item                      ItemType `item:"anchor point"`

	// These fields we know beforehand
	PrecedingPhrase           string   `json:"preceding_phrase" property:"preceding phrase"`
	FollowingPhrase           string   `json:"following_phrase" property:"following phrase"`
	DistanceToPreceding       int      `json:"preceding_distance" property:"distance to preceding"`
	DistanceToFollowing       int      `json:"follow_distance" property:"distance to following"`
	CharacterNumber           int      `json:"character" property:"character number"`
	TimeCode                  string   `json:"time" property:"time code1"`

    // These fields we only know from the science source instance
	InstanceOf                string   `json:"instance_of" property:"instance of"`

    // These we only know after we've uploaded the article document
	ScienceSourceArticleTitle string   `json:"science_source_title" property:"ScienceSource article title"`

    // These fields we know after we've created the article item
	AnchorPoint               string   `json:"point" property:"anchor point in"` // Ref to article

	// These we only know once we've uploaded all the annotations
	PrecedingAnchorPoint      string  `json:"preceding_anchor" property:"preceding anchor point"` // Ref to anchor point/article
	FollowingAnchorPoint      string  `json:"follow_anchor" property:"following anchor point"` // Ref to anchor point/terminus
	Anchors                   string   `json:"anchors" property:"anchors"`

	// Used to let us look the item up later
	ScienceSourceItemID   string `json:"id"`

    // Internal program management
	Annotation ScienceSourceAnnotation `json:"annotation"`
}

type ScienceSourceArticle struct {
    // Exists purely to let us look up the item ID on sci source
	Item                      ItemType `item:"article"`

    // These fields we know beforehand
	WikiDataItemCode          string   `json:"wikidata" property:"Wikidata item code"`
	ArticleTextTitle          string   `json:"title" property:"article text title"`
	PublicationDate           string   `json:"publication_date" property:"publication date"`
	TimeCode                  string   `json:"time" property:"time code1"`
	CharacterNumber           int      `json:"character" property:"character number"` // always 0?
	PrecedingPhrase           string   `json:"preceding_phrase" property:"preceding phrase"`
	FollowingPhrase           string   `json:"follow_phrase" property:"following phrase"`

    // These fields we only know from the science source instance
	InstanceOf                string   `json:"instance_of" property:"instance of"`

    // These we only know after we've uploaded the article
	ScienceSourceArticleTitle string   `json:"science_source_title" property:"ScienceSource article title"`

	// These we only know once we've uploaded all the annotations
	FollowingAnchorPoint      string  `json:"following_anchor" property:"following anchor point"`

	// Used to let us look the item up later
	ScienceSourceItemID   string `json:"id"`

    // Internal program management
	Annotations []ScienceSourceAnchorPoint `json:"annotations"`
}

// terminus needs looking up too

type ScienceSourceClient struct {
	wikiDataClient *WikiDataClient

	propertyMap map[string]string
	itemMap     map[string]string
}

func NewScienceSourceClient(consumerKey string, consumerSecret string, urlbase string) *ScienceSourceClient {

	res := &ScienceSourceClient{
		wikiDataClient: NewWikiDataClient(consumerKey, consumerSecret, urlbase),
		propertyMap: make(map[string]string, 0),
		itemMap: make(map[string]string, 0),
	}

	return res
}

func (c *ScienceSourceClient) GetPropertyAndItemConfigurationFromServer() error {

	list := getValuesForTags("property")
	for _, i := range list {
		label, err := c.wikiDataClient.GetPropertyForLabel(i)
		if err != nil {
			return err
		}
		c.propertyMap[i] = label
	}

	list = getValuesForTags("item")
	for _, i := range list {
		label, err := c.wikiDataClient.GetItemForLabel(i)
		if err != nil {
			return err
		}
		c.itemMap[i] = label
	}

	return nil
}

func getValuesForTags(tagname string) []string {
	tagset := make(map[string]bool, 0)

	types := [3]reflect.Type{
		reflect.TypeOf(ScienceSourceAnnotation{}),
		reflect.TypeOf(ScienceSourceAnchorPoint{}),
		reflect.TypeOf(ScienceSourceArticle{}),
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
