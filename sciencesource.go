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
	"reflect"

	"github.com/mdales/wikibase"
)

// Encoding of structures in json comes from data schema found here:
//      https://sciencesource.wmflabs.org/wiki/Data_schema

type ItemType string

type ScienceSourceAnnotation struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	Item ItemType `json:"item" item:"annotation"`

	// These fields we know beforehand
	TermFound         string `json:"term" property:"term found"`
	LengthOfTermFound int    `json:"length" property:"length of term found"`
	WikiDataItemCode  string `json:"wikidata" property:"Wikidata item code"`
	DictionaryName    string `json:"dictionary" property:"dictionary name"`
	TimeCode          string `json:"time" property:"time code1"`

	// These fields we only know from the science source instance
	InstanceOf string `json:"instance_of" property:"instance of"`

	// Used to let us look the item up later
	ScienceSourceItemID string `json:"id"`
}

type ScienceSourceAnchorPoint struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	Item ItemType `json:"item" item:"anchor point"`

	// These fields we know beforehand
	PrecedingPhrase     string `json:"preceding_phrase" property:"preceding phrase"`
	FollowingPhrase     string `json:"following_phrase" property:"following phrase"`
	DistanceToPreceding int    `json:"preceding_distance" property:"distance to preceding"`
	DistanceToFollowing int    `json:"following_distance" property:"distance to following"`
	CharacterNumber     int    `json:"character" property:"character number"`
	TimeCode            string `json:"time" property:"time code1"`

	// These fields we only know from the science source instance
	InstanceOf string `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article document
	ScienceSourceArticleTitle string `json:"science_source_title" property:"ScienceSource article title"`

	// These fields we know after we've created the article item
	AnchorPoint string `json:"point" property:"anchor point in"` // Ref to article

	// These we only know once we've uploaded all the annotations
	PrecedingAnchorPoint string `json:"preceding_anchor" property:"preceding anchor point"` // Ref to anchor point/article
	FollowingAnchorPoint string `json:"following_anchor" property:"following anchor point"` // Ref to anchor point/terminus
	Anchors              string `json:"anchors" property:"anchors"`

	// Used to let us look the item up later
	ScienceSourceItemID string `json:"id"`

	// Internal program management
	Annotation ScienceSourceAnnotation `json:"annotation"`
}

type ScienceSourceArticle struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	Item ItemType `json:"item" item:"article"`

	// These fields we know beforehand
	WikiDataItemCode string `json:"wikidata" property:"Wikidata item code"`
	ArticleTextTitle string `json:"title" property:"article text title"`
	PublicationDate  string `json:"publication_date" property:"publication date"`
	TimeCode         string `json:"time" property:"time code1"`
	CharacterNumber  int    `json:"character" property:"character number"` // always 0?
	PrecedingPhrase  string `json:"preceding_phrase" property:"preceding phrase"`
	FollowingPhrase  string `json:"following_phrase" property:"following phrase"`

	// These fields we only know from the science source instance
	InstanceOf string `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article
	ScienceSourceArticleTitle string `json:"science_source_title" property:"ScienceSource article title"`
	PageID                    int    `json:"page_id" property:"page ID"`

	// These we only know once we've uploaded all the annotations
	FollowingAnchorPoint string `json:"following_anchor" property:"following anchor point"`

	// Internal program management
	Annotations []ScienceSourceAnchorPoint `json:"annotations"`
}

// terminus needs looking up too

type ScienceSourceClient struct {
	wikiBaseClient *wikibase.WikiBaseClient

	PropertyMap map[string]string
	ItemMap     map[string]string
}

func NewScienceSourceClient(oauthInfo wikibase.WikiBaseOAuthInformation, urlbase string) *ScienceSourceClient {

	oauth_client := wikibase.NewOAuthClient(oauthInfo, urlbase)

	res := &ScienceSourceClient{
		wikiBaseClient: wikibase.NewWikiBaseClient(oauth_client),
		PropertyMap:    make(map[string]string, 0),
		ItemMap:        make(map[string]string, 0),
	}

	return res
}

func (c *ScienceSourceClient) GetPropertyAndItemConfigurationFromServer() error {

	list := getValuesForTags("property")
	for _, i := range list {
		label, err := c.wikiBaseClient.GetPropertyForLabel(i)
		if err != nil {
			return err
		}
		c.PropertyMap[i] = label
	}

	list = getValuesForTags("item")
	for _, i := range list {
		label, err := c.wikiBaseClient.GetItemForLabel(i)
		if err != nil {
			return err
		}
		c.ItemMap[i] = label
	}

	return nil
}

func (c *ScienceSourceClient) UploadPaper(article *ScienceSourceArticle, htmlFileName string) error {

	data, err := ioutil.ReadFile(htmlFileName)
	if err != nil {
		return err
	}

	page_id, upload_error := c.wikiBaseClient.CreateArticle(article.ScienceSourceArticleTitle, string(data))
	if upload_error != nil {
		return upload_error
	}

	article.PageID = page_id

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

// Article helper functions

func (article *ScienceSourceArticle) Save(filename string) error {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(article)
}

func LoadScienceSourceArticle(filename string) (*ScienceSourceArticle, error) {

	var article ScienceSourceArticle

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&article)
	return &article, err
}

// Wiki base item related code

func (c *ScienceSourceClient) CreateArticleItemTree(article *ScienceSourceArticle) error {

	// Create the node for the article in the wiki base if necessary
	if len(article.Item) == 0 {
		item_id, err := c.wikiBaseClient.CreateItemInstance("article")
		if err != nil {
			return err
		}
		article.Item = ItemType(item_id)
	}

	// Create an item for all the anchors and their articles
	for _, anchor := range article.Annotations {

		if len(anchor.Item) == 0 {
			item_id, err := c.wikiBaseClient.CreateItemInstance("anchor")
			if err != nil {
				return err
			}
			anchor.Item = ItemType(item_id)
		}

		if len(anchor.Annotation.Item) == 0 {
			item_id, err := c.wikiBaseClient.CreateItemInstance("anchor")
			if err != nil {
				return err
			}
			anchor.Annotation.Item = ItemType(item_id)
		}
	}

	return nil
}

func (c *ScienceSourceClient) PopulateAritcleItemTree(article *ScienceSourceArticle) error {
	return nil
}
