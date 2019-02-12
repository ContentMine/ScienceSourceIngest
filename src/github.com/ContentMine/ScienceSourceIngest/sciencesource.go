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
	"time"

	"github.com/ContentMine/wikibase"
)

// Encoding of structures in json comes from data schema found here:
//      https://sciencesource.wmflabs.org/wiki/Data_schema

type ScienceSourceAnnotation struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	wikibase.ItemHeader `json:"item" item:"annotation"`

	// These fields we know beforehand
	TermFound         string    `json:"term" property:"term found"`
	LengthOfTermFound int       `json:"length" property:"length of term found"`
	WikiDataItemCode  string    `json:"wikidata" property:"Wikidata item code"`
	DictionaryName    string    `json:"dictionary" property:"dictionary name"`
	TimeCode          time.Time `json:"time" property:"time code1"`

	// These fields we only know from the science source instance
	InstanceOf wikibase.ItemPropertyType `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article document
	ScienceSourceArticleTitle string `json:"science_source_title" property:"ScienceSource article title"`

	// These fields we know after we've created the anchor point item
	BasedOn wikibase.ItemPropertyType `json:"based_on" property:"based on,omitoncreate"` // Ref to article
}

type ScienceSourceAnchorPoint struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	wikibase.ItemHeader `json:"item" item:"anchor point"`

	// These fields we know beforehand
	PrecedingPhrase     string    `json:"preceding_phrase" property:"preceding phrase"`
	FollowingPhrase     string    `json:"following_phrase" property:"following phrase"`
	DistanceToPreceding *int      `json:"preceding_distance,omitempty" property:"distance to preceding"`
	DistanceToFollowing *int      `json:"following_distance,omitempty" property:"distance to following"`
	CharacterNumber     int       `json:"character" property:"character number"`
	TimeCode            time.Time `json:"time" property:"time code1"`

	// These fields we only know from the science source instance
	InstanceOf wikibase.ItemPropertyType `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article document
	ScienceSourceArticleTitle string `json:"science_source_title" property:"ScienceSource article title"`

	// These fields we know after we've created the article item
	AnchorPoint wikibase.ItemPropertyType `json:"point" property:"anchor point in,omitoncreate"` // Ref to article

	// These we only know once we've uploaded all the annotations
	PrecedingAnchorPoint *wikibase.ItemPropertyType `json:"preceding_anchor,omitempty" property:"preceding anchor point,omitoncreate"` // Ref to anchor point/article
	FollowingAnchorPoint wikibase.ItemPropertyType  `json:"following_anchor" property:"following anchor point,omitoncreate"`           // Ref to anchor point/terminus
	Anchors              wikibase.ItemPropertyType  `json:"anchors" property:"anchors,omitoncreate"`

	// Internal program management
	Annotation ScienceSourceAnnotation `json:"annotation"`
}

type ScienceSourceArticle struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	wikibase.ItemHeader `json:"item" item:"article"`

	// These fields we know beforehand
	ScienceSourceArticleTitle string                     `json:"science_source_title" property:"ScienceSource article title"`
	WikiDataItemCode          string                     `json:"wikidata" property:"Wikidata item code"`
	ArticleTextTitle          string                     `json:"title" property:"article text title"`
	PublicationDate           time.Time                  `json:"publication_date" property:"publication date"`
	TimeCode                  time.Time                  `json:"time" property:"time code1"`
	CharacterNumber           int                        `json:"character" property:"character number"` // always 0?
	PrecedingPhrase           *string                    `json:"preceding_phrase,omitempty" property:"preceding phrase"`
	FollowingPhrase           *string                    `json:"following_phrase,omitempty" property:"following phrase"`
	PrecedingAnchorPoint      *wikibase.ItemPropertyType `json:"preceding_anchor,omitempty" property:"preceding anchor point"` // Always nil on article

	// These fields we only know from the science source instance
	InstanceOf wikibase.ItemPropertyType `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article
	PageID int `json:"page_id" property:"page ID"`

	// These we only know once we've uploaded all the annotations
	FollowingAnchorPoint wikibase.ItemPropertyType `json:"following_anchor" property:"following anchor point,omitoncreate"`

	// Internal program management
	Annotations []ScienceSourceAnchorPoint `json:"annotations"`
}

// terminus needs looking up too

type ScienceSourceClient struct {
	wikiBaseClient *wikibase.Client
}

func NewScienceSourceClient(oauthInfo wikibase.OAuthInformation, urlbase string) *ScienceSourceClient {

	oauth_client := wikibase.NewOAuthNetworkClient(oauthInfo, urlbase)

	res := &ScienceSourceClient{
		wikiBaseClient: wikibase.NewClient(oauth_client),
	}

	return res
}

func (c *ScienceSourceClient) GetConfigurationFromServer() error {

	err := c.wikiBaseClient.MapPropertyAndItemConfiguration(ScienceSourceArticle{}, true)
	if err != nil {
		return err
	}
	err = c.wikiBaseClient.MapPropertyAndItemConfiguration(ScienceSourceAnchorPoint{}, true)
	if err != nil {
		return err
	}
	err = c.wikiBaseClient.MapPropertyAndItemConfiguration(ScienceSourceAnnotation{}, true)
	if err != nil {
		return err
	}

	err = c.wikiBaseClient.MapItemConfigurationByLabel("terminus", true)
	if err != nil {
		return err
	}

	return nil
}

func (c *ScienceSourceClient) UploadPaper(article *ScienceSourceArticle, htmlFileName string) error {

	data, err := ioutil.ReadFile(htmlFileName)
	if err != nil {
		return err
	}

	page_id, upload_error := c.wikiBaseClient.CreateOrUpdateArticle(article.ScienceSourceArticleTitle, string(data))
	if upload_error != nil {

		// if we get a page exists error then ignore for now and move on, as we assume the title is unique
		ignore_error := false
		if err, ok := upload_error.(*wikibase.APIError); ok {
			ignore_error = err.Code == "articleexists"
		}

		if ignore_error != true {
			return upload_error
		}
	}

	article.PageID = page_id

	return c.wikiBaseClient.ProtectPageByID(article.PageID)
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
	article.InstanceOf = c.wikiBaseClient.ItemMap["article"]
	if len(article.ID) == 0 {
		err := c.wikiBaseClient.CreateItemInstance("article instance", article)
		if err != nil {
			return err
		}
	}

	// Create an item for all the anchors and their articles
	for i := 0; i < len(article.Annotations); i++ {
		article.Annotations[i].InstanceOf = c.wikiBaseClient.ItemMap["anchor point"]

		if len(article.Annotations[i].ID) == 0 {
			err := c.wikiBaseClient.CreateItemInstance("anchor instance", &(article.Annotations[i]))
			if err != nil {
				return err
			}
		}

		article.Annotations[i].Annotation.InstanceOf = c.wikiBaseClient.ItemMap["annotation"]
		if len(article.Annotations[i].Annotation.ID) == 0 {
			err := c.wikiBaseClient.CreateItemInstance("annotation instance", &(article.Annotations[i].Annotation))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ScienceSourceClient) ReconsileArticleItemTree(article *ScienceSourceArticle) error {

	// Patch the article first
	if len(article.Annotations) == 0 {
		article.FollowingAnchorPoint = c.wikiBaseClient.ItemMap["terminus"]
	} else {
		article.FollowingAnchorPoint = article.Annotations[0].ID
	}

	for i := 0; i < len(article.Annotations); i++ {

		// Patch anchor point first
		if i != 0 {
			article.Annotations[i].PrecedingAnchorPoint = &article.Annotations[i-1].ID
		} else {
			article.Annotations[i].PrecedingAnchorPoint = nil
		}
		if i != len(article.Annotations)-1 {
			article.Annotations[i].FollowingAnchorPoint = article.Annotations[i+1].ID
		} else {
			article.Annotations[i].FollowingAnchorPoint = c.wikiBaseClient.ItemMap["terminus"]
		}
		article.Annotations[i].AnchorPoint = article.ID
		article.Annotations[i].Anchors = article.Annotations[i].Annotation.ID

		// Patch annotation second
		article.Annotations[i].Annotation.BasedOn = article.Annotations[i].ID
	}

	return nil
}

func (c *ScienceSourceClient) PopulateAritcleItemTree(article *ScienceSourceArticle) error {

	err := c.wikiBaseClient.UploadClaimsForItem(article, false)
	if err != nil {
		return err
	}

	for i := 0; i < len(article.Annotations); i++ {
		err := c.wikiBaseClient.UploadClaimsForItem(&article.Annotations[i], false)
		if err != nil {
			return err
		}
		err = c.wikiBaseClient.UploadClaimsForItem(&(article.Annotations[i].Annotation), false)
		if err != nil {
			return err
		}
	}

	return nil
}
