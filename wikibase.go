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
	"strings"
	"sync"

	"github.com/mrjones/oauth"
)

type WikiBaseType string
const (
    WikiBaseProperty WikiBaseType = "property"
    WikiBaseItem WikiBaseType = "item"
)

type GeneralAPIResponse struct {
	BatchComplete string  `json:"batchcomplete"`
	RequestID     *string `json:"requestid"`
}

type Token struct {
	CSRFToken *string `json:"csrftoken"`
}

type TokensQuery struct {
	Tokens Token `json:"tokens"`
}

type TokenRequestResponse struct {
	GeneralAPIResponse
	Query TokensQuery `json:"query"`
}

type WikiBaseSearchItem struct {
	Duration    int    `json:"ns"`
	Title       string `json:"title"`
	PageID      int    `json:"pageid"`
	DisplayText string    `json:"displaytext"`
}

type WikiBaseSearchQuery struct {
	Items []WikiBaseSearchItem `json:"wbsearch"`
}

type WikiBaseSearchResponse struct {
	GeneralAPIResponse
	Query WikiBaseSearchQuery `json:"query"`
}

// TODO - clearly needs to die
const WikiBaseConsumerKey string = "44fc577d47dd15516fcbe4dfe78777cd"
const WikiBaseConsumerSecret string = "183d5f332679148ce5090558b5d58afb1689a9fb"

const accessToken string = "cd7c6e9a2954c52a29ccb942cf356d46"
const accessSecret string = "0c46c3ffba5d8aff28665786ed342b361d14ba10"

type WikiDataClient struct {
	ConsumerKey    string
	ConsumerSecret string

	URLBase string

	AccessToken *oauth.AccessToken
	consumer    *oauth.Consumer

	// Don't read directly - use GetEditingToken()
	editToken     *string
	editTokenLock sync.RWMutex
}

func NewWikiDataClient(consumerKey string, consumerSecret string, urlbase string) *WikiDataClient {

	res := WikiDataClient{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		URLBase:        urlbase,
	}

	res.consumer = oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("%s/wiki/Special:OAuth/initiate", urlbase),
			AuthorizeTokenUrl: fmt.Sprintf("%s/wiki/Special:OAuth/authorize", urlbase),
			AccessTokenUrl:    fmt.Sprintf("%s/wiki/Special:OAuth/token", urlbase),
		})

	// Kill this
	aToken := oauth.AccessToken{Token: accessToken, Secret: accessSecret}
	res.AccessToken = &aToken

	return &res
}

func (c *WikiDataClient) GetEditingToken() (string, error) {

	c.editTokenLock.RLock()
	initVal := c.editToken
	c.editTokenLock.RUnlock()

	if initVal != nil {
		return *initVal, nil
	}

	c.editTokenLock.Lock()
	defer c.editTokenLock.Unlock()

	// at start of day there's a big risk all go-routines race on getting
	// the edit token, so bail early if someone else has won
	if c.editToken != nil {
		return *c.editToken, nil
	}

	response, err := c.consumer.Get(
		fmt.Sprintf("%s/w/api.php", c.URLBase),
		map[string]string{
			"action": "query",
			"meta":   "tokens",
			"format": "json",
		},
		c.AccessToken)

	if err != nil {
		return "", err
	}

	var token TokenRequestResponse
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		return "", err
	}

	if token.Query.Tokens.CSRFToken == nil {
		return "", fmt.Errorf("Failed to get token in response from server: %v", token)
	}

	c.editToken = token.Query.Tokens.CSRFToken

	return *c.editToken, nil
}

func (c *WikiDataClient) getWikibaseThingForLabel(thing WikiBaseType, label string) (string, error) {

	response, err := c.consumer.Get(
		fmt.Sprintf("%s/w/api.php", c.URLBase),
		map[string]string{
			"action":      "query",
			"list":        "wbsearch",
			"wbssearch":   label,
			"wbstype":     string(thing),
			"wbslanguage": "en",
			"format":      "json",
		},
		c.AccessToken)

	if err != nil {
		return "", err
	}

	var search WikiBaseSearchResponse
	err = json.NewDecoder(response.Body).Decode(&search)
	if err != nil {
		return "", err
	}

	switch len(search.Query.Items) {
	case 0:
		return "", fmt.Errorf("Failed to find any matching properties for %s", label)
	case 1:
        parts := strings.Split(search.Query.Items[0].Title, ":")
        if len(parts) != 2 {
            return "", fmt.Errorf("We expected type:value in reply, but got: %v", search.Query.Items[0].Title)
        }
		return parts[1], nil // TODO fix
	default:
		return "", fmt.Errorf("Too many matches returned: %d", len(search.Query.Items))
	}
}

func (c *WikiDataClient) GetPropertyForLabel(label string) (string, error) {
    return c.getWikibaseThingForLabel(WikiBaseProperty, label)
}

func (c *WikiDataClient) GetItemForLabel(label string) (string, error) {
    return c.getWikibaseThingForLabel(WikiBaseItem, label)
}
