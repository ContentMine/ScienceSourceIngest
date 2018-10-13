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
	"encoding/xml"
	"fmt"
	"os"
)

type ContributorName struct {
	Surname    string `xml:"surname"`
	GivenNames string `xml:"given-names"`
}

type Contributor struct {
	Name ContributorName `xml:"name"`
}

type ContribGroup struct {
	XMLName      xml.Name      `xml:"contrib-group"`
	Contributors []Contributor `xml:"contrib"`
}

type TitleGroup struct {
	XMLName          xml.Name `xml:"title-group"`
	ArticleTitle     string   `xml:"article-title"`
	AlternativeTitle string   `xml:"alt-title"`
}

type JournalMeta struct {
	XMLName xml.Name `xml:"journal-meta"`
}

type ArticleMeta struct {
	XMLName           xml.Name       `xml:"article-meta"`
	TitleGroup        TitleGroup     `xml:"title-group"`
	ContributorGroups []ContribGroup `xml:"contrib-group"`
}

type Front struct {
	XMLName     xml.Name    `xml:"front"`
	JournalMeta JournalMeta `xml:"journal-meta"`
	ArticleMeta ArticleMeta `xml:"article-meta"`
}

type OpenXMLPaper struct {
	XMLName xml.Name `xml:"article"`
	Front   Front    `xml:"front"`
}

// Parsing

func LoadPaperXMLFromFile(path string) (OpenXMLPaper, error) {
	var paper OpenXMLPaper

	f, err := os.Open(path)
	if err != nil {
		return OpenXMLPaper{}, err
	}

	err = xml.NewDecoder(f).Decode(&paper)
	return paper, err
}

// Convenience functions

func (author *ContributorName) String() string {
	if author == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s %s", author.GivenNames, author.Surname)
}

func (paper OpenXMLPaper) Title() string {
	return paper.Front.ArticleMeta.TitleGroup.ArticleTitle
}

func (paper OpenXMLPaper) FirstAuthor() *ContributorName {
	contrib_groups := paper.Front.ArticleMeta.ContributorGroups
	if len(contrib_groups) > 0 {
		author_list := contrib_groups[0].Contributors
		if len(author_list) > 0 {
			return &(author_list[0].Name)
		}
	}
	return nil
}
