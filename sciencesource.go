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

	"github.com/mdales/wikibase"
)

// Encoding of structures in json comes from data schema found here:
//      https://sciencesource.wmflabs.org/wiki/Data_schema

type ScienceSourceAnnotation struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	Item wikibase.ItemPropertyType `json:"item" item:"annotation"`

	// These fields we know beforehand
	TermFound         string    `json:"term" property:"term found"`
	LengthOfTermFound int       `json:"length" property:"length of term found"`
	WikiDataItemCode  string    `json:"wikidata" property:"Wikidata item code"`
	DictionaryName    string    `json:"dictionary" property:"dictionary name"`
	TimeCode          time.Time `json:"time" property:"time code1"`

	// These fields we know after we've created the anchor point item
	BasedOn wikibase.ItemPropertyType `json:"based_on" property:"based on"` // Ref to article

	// These fields we only know from the science source instance
	InstanceOf wikibase.ItemPropertyType `json:"instance_of" property:"instance of"`
}

type ScienceSourceAnchorPoint struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	Item wikibase.ItemPropertyType `json:"item" item:"anchor point"`

	// These fields we know beforehand
	PrecedingPhrase     string    `json:"preceding_phrase" property:"preceding phrase"`
	FollowingPhrase     string    `json:"following_phrase" property:"following phrase"`
	DistanceToPreceding int       `json:"preceding_distance" property:"distance to preceding"`
	DistanceToFollowing int       `json:"following_distance" property:"distance to following"`
	CharacterNumber     int       `json:"character" property:"character number"`
	TimeCode            time.Time `json:"time" property:"time code1"`

	// These fields we only know from the science source instance
	InstanceOf wikibase.ItemPropertyType `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article document
	ScienceSourceArticleTitle string `json:"science_source_title" property:"ScienceSource article title"`

	// These fields we know after we've created the article item
	AnchorPoint wikibase.ItemPropertyType `json:"point" property:"anchor point in"` // Ref to article

	// These we only know once we've uploaded all the annotations
	PrecedingAnchorPoint wikibase.ItemPropertyType `json:"preceding_anchor" property:"preceding anchor point"` // Ref to anchor point/article
	FollowingAnchorPoint wikibase.ItemPropertyType `json:"following_anchor" property:"following anchor point"` // Ref to anchor point/terminus
	Anchors              wikibase.ItemPropertyType `json:"anchors" property:"anchors"`

	// Internal program management
	Annotation ScienceSourceAnnotation `json:"annotation"`
}

type ScienceSourceArticle struct {
	// Exists partly to let us look up the item ID on sci source, and as a place to store the uploaded
	// wikibase item ID when we cache state to disk
	Item wikibase.ItemPropertyType `json:"item" item:"article"`

	// These fields we know beforehand
	ScienceSourceArticleTitle string    `json:"science_source_title" property:"ScienceSource article title"`
	WikiDataItemCode          string    `json:"wikidata" property:"Wikidata item code"`
	ArticleTextTitle          string    `json:"title" property:"article text title"`
	PublicationDate           time.Time `json:"publication_date" property:"publication date"`
	TimeCode                  time.Time `json:"time" property:"time code1"`
	CharacterNumber           int       `json:"character" property:"character number"` // always 0?
	//PrecedingPhrase           string    `json:"preceding_phrase" property:"preceding phrase"`
	//FollowingPhrase           string    `json:"following_phrase" property:"following phrase"`

	// These fields we only know from the science source instance
	InstanceOf wikibase.ItemPropertyType `json:"instance_of" property:"instance of"`

	// These we only know after we've uploaded the article
	PageID int `json:"page_id" property:"page ID"`

	// These we only know once we've uploaded all the annotations
	FollowingAnchorPoint wikibase.ItemPropertyType `json:"following_anchor" property:"following anchor point"`

	// Internal program management
	Annotations []ScienceSourceAnchorPoint `json:"annotations"`
}

// terminus needs looking up too

type ScienceSourceClient struct {
	wikiBaseClient *wikibase.WikiBaseClient
}

func NewScienceSourceClient(oauthInfo wikibase.WikiBaseOAuthInformation, urlbase string) *ScienceSourceClient {

	oauth_client := wikibase.NewOAuthClient(oauthInfo, urlbase)

	res := &ScienceSourceClient{
		wikiBaseClient: wikibase.NewWikiBaseClient(oauth_client),
	}

	return res
}

func (c *ScienceSourceClient) GetConfigurationFromServer() error {

    err := c.wikiBaseClient.MapPropertyAndItemConfiguration(ScienceSourceArticle{})
    if err != nil {
        return err
    }
    err = c.wikiBaseClient.MapPropertyAndItemConfiguration(ScienceSourceAnchorPoint{})
    if err != nil {
        return err
    }
    err = c.wikiBaseClient.MapPropertyAndItemConfiguration(ScienceSourceAnnotation{})
    if err != nil {
        return err
    }

    err = c.wikiBaseClient.MapItemConfigurationByLabel("terminus")
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

	page_id, upload_error := c.wikiBaseClient.CreateArticle(article.ScienceSourceArticleTitle, string(data))
	if upload_error != nil {

		// if we get a page exists error then ignore for now and move on, as we assume the title is unique
		ignore_error := false
		if err, ok := upload_error.(*wikibase.WikiBaseError); ok {
			ignore_error = err.Code == "articleexists"
		}

		if ignore_error != true {
			return upload_error
		}
	}

	article.PageID = page_id

	return nil
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
		item_id, err := c.wikiBaseClient.CreateItemInstance("article instance")
		if err != nil {
			return err
		}
		article.Item = item_id
	}

	// Create an item for all the anchors and their articles
	for i := 0; i < len(article.Annotations); i++ {

		if len(article.Annotations[i].Item) == 0 {
			item_id, err := c.wikiBaseClient.CreateItemInstance("anchor instance")
			if err != nil {
				return err
			}
			article.Annotations[i].Item = item_id
		}

		if len(article.Annotations[i].Annotation.Item) == 0 {
			item_id, err := c.wikiBaseClient.CreateItemInstance("annotation instance")
			if err != nil {
				return err
			}
			article.Annotations[i].Annotation.Item = item_id
		}
	}

	return nil
}

func (c *ScienceSourceClient) ReconsileArticleItemTree(article *ScienceSourceArticle) error {

	// Patch the article first
	article.InstanceOf = c.wikiBaseClient.ItemMap["article"]
	if len(article.Annotations) == 0 {
		article.FollowingAnchorPoint = c.wikiBaseClient.ItemMap["terminus"]
	} else {
		article.FollowingAnchorPoint = article.Annotations[0].Item
	}

	for i := 0; i < len(article.Annotations); i++ {

		// Patch anchor point first
		article.Annotations[i].InstanceOf = c.wikiBaseClient.ItemMap["anchor point"]
		article.Annotations[i].ScienceSourceArticleTitle = article.ScienceSourceArticleTitle
		if i != 0 {
			article.Annotations[i].PrecedingAnchorPoint = article.Annotations[i-1].Item
		} else {
			article.Annotations[i].PrecedingAnchorPoint = c.wikiBaseClient.ItemMap["terminus"]
		}
		if i != len(article.Annotations)-1 {
			article.Annotations[i].FollowingAnchorPoint = article.Annotations[i+1].Item
		} else {
			article.Annotations[i].FollowingAnchorPoint = c.wikiBaseClient.ItemMap["terminus"]
		}
		article.Annotations[i].AnchorPoint = article.Item
		article.Annotations[i].Anchors = article.Annotations[i].Annotation.Item

		// Patch annotation second
		article.Annotations[i].Annotation.InstanceOf = c.wikiBaseClient.ItemMap["annotation"]
		article.Annotations[i].Annotation.BasedOn = article.Annotations[i].Item
	}

	return nil
}

func (c *ScienceSourceClient) PopulateAritcleItemTree(article *ScienceSourceArticle) error {

    err := c.wikiBaseClient.UploadClaimsForItem(article.Item, *article)
    if err != nil {
        return err
    }

	for i := 0; i < len(article.Annotations); i++ {
        err := c.wikiBaseClient.UploadClaimsForItem(article.Annotations[i].Item, article.Annotations[i])
        if err != nil {
            return err
        }
        err = c.wikiBaseClient.UploadClaimsForItem(article.Annotations[i].Annotation.Item, article.Annotations[i].Annotation)
        if err != nil {
            return err
        }
    }

    return nil
}
