// Copyright ContentMine Ltd 2018

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
