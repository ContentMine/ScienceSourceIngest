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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"time"
)

type PaperProcessor struct {
	Paper           Paper
	TargetDirectory string
}

const HTMLHeader string = `{{headertemplate
| title = %s
| publication_date = %s
| initial_author_first = %s
| initial_author_last = %s
| license = %s
| wikidata = %s
| main_subject = %s
| DOI =
| PubMed_ID =
| PMC_ID =
| Generator = %s/%s
}}
`

type SearchDirection int

const (
	SearchDirectionBackward SearchDirection = -1
	SearchDirectionForward                  = 1
)

const PhraseTargetSize int = 100

// Generic helpers

func fetchResource(url string, filename string) error {

	// if it already exists, don't fetch it again
	if _, err := os.Stat(filename); err == nil {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, resp_err := http.Get(url)
	if resp_err != nil {
		return resp_err
	}
	defer resp.Body.Close()

	_, copy_err := io.Copy(f, resp.Body)
	return copy_err
}

func findPhrase(prose []byte, startOffset int, direction SearchDirection) string {

	targetOffset := startOffset + (PhraseTargetSize * int(direction))

	// Need better terminating condition here
	for true {
		if direction == SearchDirectionBackward {
			if targetOffset < 0 {
				targetOffset = 0
				break
			}
		} else {
			if targetOffset > (len(prose) - 1) {
				targetOffset = len(prose) - 1
				break
			}
		}

		if prose[targetOffset] == byte(' ') {
			break
		}

		targetOffset = targetOffset + (1 * int(direction))
	}

	if startOffset > targetOffset {
		startOffset, targetOffset = targetOffset, startOffset
	}

	return string(prose[startOffset:targetOffset])
}

// Computed properties

func (processor PaperProcessor) folderName() string {
	return path.Join(processor.TargetDirectory, processor.Paper.ID())
}

func (processor PaperProcessor) targetXMLFileName() string {
	return path.Join(processor.folderName(), "paper.xml")
}

func (processor PaperProcessor) targetHTMLFileName() string {
	return path.Join(processor.folderName(), "paper.html")
}

func (processor PaperProcessor) targetTextFileName() string {
	return path.Join(processor.folderName(), "paper.txt")
}

func (processor PaperProcessor) targetAnnotationsFileName() string {
	return path.Join(processor.folderName(), "annotations.json")
}

func (processor PaperProcessor) targetSupplementaryArchiveFileName() string {
	return path.Join(processor.folderName(), "supplementary.zip")
}

// Side effect heavy functions

func (processor PaperProcessor) createFolderIfRequired() error {
	return os.MkdirAll(processor.folderName(), 0755)
}

func (processor PaperProcessor) fetchPaperTextToDisk() error {
	return fetchResource(processor.Paper.FullTextURL(), processor.targetXMLFileName())
}

func (processor PaperProcessor) fetchPaperSupplementaryFilesToDisk() error {
	return fetchResource(processor.Paper.SupplementaryFilesURL(), processor.targetSupplementaryArchiveFileName())
}

func (processor PaperProcessor) processXMLToHTML(FirstAuthor *ContributorName) error {

	f, err := os.Create(processor.targetHTMLFileName())
	if err != nil {
		return err
	}
	defer f.Close()

	firstName := ""
	surname := ""
	if FirstAuthor != nil {
		surname = FirstAuthor.Surname
		firstName = FirstAuthor.GivenNames // TODO!
	}

	header := fmt.Sprintf(HTMLHeader,
		processor.Paper.Title.Value,
		processor.Paper.Date.Value,
		firstName, surname,
		processor.Paper.LicenseLabel.Value,
		processor.Paper.ItemLabel.Value,
		processor.Paper.MainSubjectLabel.Value,
		Remote, Version,
	)

	f.Write([]byte(header))

	cmd := exec.Cmd{
		Path: "/usr/bin/xsltproc",
		Args: []string{"xsltproc", "jats-parsoid.xsl", processor.targetXMLFileName()},
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	_, copy_err := io.Copy(f, stdout)
	if copy_err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	// write the footer
	f.Write([]byte("{{footertemplate}}\n"))

	return nil
}

func (processor PaperProcessor) processXMLToText() error {

	f, err := os.Create(processor.targetTextFileName())
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Cmd{
		Path: "/usr/bin/xsltproc",
		Args: []string{"xsltproc", "jats-text.xsl", processor.targetXMLFileName()},
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	_, copy_err := io.Copy(f, stdout)
	if copy_err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (processor PaperProcessor) findAnnotations(dictionaries []Dictionary,
	articleTitle string, journalTitle string) ([]Annotation, error) {

	data, err := ioutil.ReadFile(processor.targetTextFileName())
	if err != nil {
		return nil, err
	}

	total_matches := make([]DictionaryMatch, 0)

	for _, dictionary := range dictionaries {
		total_matches = append(total_matches, dictionary.FindMatches(data)...)
	}

	sort.Sort(DictionaryMatchesByOffset(total_matches))

	res := make([]Annotation, len(total_matches))

	for i := 0; i < len(total_matches); i++ {
		match := total_matches[i]

		distanceToPreceding := 0
		if i > 0 {
			distanceToPreceding = match.Offset - total_matches[i-1].Offset
		}
		distanceToFollowing := 0
		if i < (len(total_matches) - 1) {
			distanceToFollowing = total_matches[i+1].Offset - match.Offset
		}

		annotation := Annotation{
			WikiDataItemCode: match.Entry.Identifiers.WikiData,
			// InstanceOf
			// SubclassOf
			// PrecedingAnchorPoint
			// FollowingAnchorPoint
			DistanceToPreceding: distanceToPreceding,
			DistanceToFollowing: distanceToFollowing,
			CharacterNumber:     match.Offset,
			ArticleTextTitle:    articleTitle,
			// AnchorPoint,
			PrecedingPhrase:   findPhrase(data, match.Offset, SearchDirectionBackward),
			FollowingPhrase:   findPhrase(data, match.Offset+len(match.Entry.Term), SearchDirectionForward),
			TermFound:         match.Entry.Term,
			DictionaryName:    match.Dictionary.Identifier,
			PublicationName:   journalTitle,
			LengthOfTermFound: len(match.Entry.Term),
			// BasedOn
			// ScienceSourceArticleTitle
			TimeCode: time.Now().String(), // TODO: Format properly
			// Anchors
			FormatterURL: "https://www.wikidata.org/wiki/$",
		}

		res[i] = annotation
	}

	return res, nil
}

func (processor PaperProcessor) saveAnnotations(annotations []Annotation) error {

	f, err := os.Create(processor.targetAnnotationsFileName())
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(annotations)
}

// main entry point

func (processor PaperProcessor) ProcessPaper(dictionaries []Dictionary) error {

	err := processor.createFolderIfRequired()
	if err != nil {
		return err
	}

	err = processor.fetchPaperTextToDisk()
	if err != nil {
		return err
	}

	err = processor.fetchPaperSupplementaryFilesToDisk()
	if err != nil {
		return err
	}

	openXMLdoc, xml_err := LoadPaperXMLFromFile(processor.targetXMLFileName())
	if xml_err != nil {
		return xml_err
	}

	err = processor.processXMLToHTML(openXMLdoc.FirstAuthor())
	if err != nil {
		return err
	}

	err = processor.processXMLToText()
	if err != nil {
		return err
	}

	annotations, aerr := processor.findAnnotations(dictionaries, openXMLdoc.Title(),
		openXMLdoc.JournalTitle())
	if aerr != nil {
		return aerr
	}

	log.Printf("Count %d", len(annotations))

	err = processor.saveAnnotations(annotations)
	if err != nil {
		return err
	}

	return nil
}
