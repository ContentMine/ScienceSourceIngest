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
	Paper               Paper
	TargetDirectory     string
	ScienceSourceRecord *ScienceSourceArticle
}

const HTMLHeader string = `{{articleheader
| title = %s
| publication_date = %s
| initial_author_first = %s
| initial_author_last = %s
| wikidata_code = %s
| Generator = %s/%s
}}
`

const HTMLFooter string = `{{articlefooter
| pmcid = %s
| license = %s
| main_subject = %s
| batch_date = %d-%d-%d
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

func (processor PaperProcessor) targetScienceSourceStateFileName() string {
	return path.Join(processor.folderName(), "scisource.json")
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

// Main processing functions

func (processor PaperProcessor) populateScienceSourceArticle() *ScienceSourceArticle {
	article := &ScienceSourceArticle{
		WikiDataItemCode:          processor.Paper.WikiDataID(),
		ArticleTextTitle:          processor.Paper.Title.Value,
		ScienceSourceArticleTitle: fmt.Sprintf("%s (%s)", processor.Paper.Title.Value, processor.Paper.ID()),
		PublicationDate:           processor.Paper.Date.Value,
		TimeCode:                  time.Now().String(),
	}

	return article
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
		processor.Paper.ItemLabel.Value,
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

	now := time.Now()
	footer := fmt.Sprintf(HTMLFooter,
		processor.Paper.PMCID.Value,
		processor.Paper.LicenseLabel.Value,
		processor.Paper.MainSubjectLabel.Value,
		now.Year(), now.Month(), now.Day(),
	)

	// write the footer
	f.Write([]byte(footer))

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
	articleTitle string, journalTitle string) ([]ScienceSourceAnchorPoint, error) {

	data, err := ioutil.ReadFile(processor.targetTextFileName())
	if err != nil {
		return nil, err
	}

	total_matches := make([]DictionaryMatch, 0)

	for _, dictionary := range dictionaries {
		total_matches = append(total_matches, dictionary.FindMatches(data)...)
	}

	sort.Sort(DictionaryMatchesByOffset(total_matches))

	res := make([]ScienceSourceAnchorPoint, len(total_matches))

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

		now := time.Now().String() // TODO: Format properly

		annotation := ScienceSourceAnnotation{
			TermFound:         match.Entry.Term,
			DictionaryName:    match.Dictionary.Identifier,
			WikiDataItemCode:  match.Entry.Identifiers.WikiData,
			LengthOfTermFound: len(match.Entry.Term),
			TimeCode:          now,
		}

		anchorPoint := ScienceSourceAnchorPoint{
			PrecedingPhrase:     findPhrase(data, match.Offset, SearchDirectionBackward),
			FollowingPhrase:     findPhrase(data, match.Offset+len(match.Entry.Term), SearchDirectionForward),
			DistanceToPreceding: distanceToPreceding,
			DistanceToFollowing: distanceToFollowing,
			CharacterNumber:     match.Offset,
			TimeCode:            now,

			Annotation: annotation,
		}

		res[i] = anchorPoint
	}

	return res, nil
}

// main entry point

func (processor PaperProcessor) ProcessPaper(dictionaries []Dictionary, sciSourceClient *ScienceSourceClient) error {

	err := processor.createFolderIfRequired()
	if err != nil {
		return err
	}

	// Have we already processed this paper?
	processor.ScienceSourceRecord, err = LoadScienceSourceArticle(processor.targetScienceSourceStateFileName())
	if err != nil {
		processor.ScienceSourceRecord = processor.populateScienceSourceArticle()

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

		processor.ScienceSourceRecord.Annotations = annotations

		// Save the record with annotations
		err = processor.ScienceSourceRecord.Save(processor.targetScienceSourceStateFileName())
		if err != nil {
			return err
		}
	}
	log.Printf("Count %d", len(processor.ScienceSourceRecord.Annotations))

	if processor.ScienceSourceRecord.PageID == 0 {
		log.Printf("Uploading paper %s", processor.Paper.ID())
		upload_err := sciSourceClient.UploadPaper(processor.ScienceSourceRecord, processor.targetHTMLFileName())
		if upload_err != nil {
			return upload_err
		}

		log.Printf("Page ID is %d", processor.ScienceSourceRecord.PageID)

		// Save the record again as it'll have an updated Page ID
		err = processor.ScienceSourceRecord.Save(processor.targetScienceSourceStateFileName())
		if err != nil {
			return err
		}
	}

	// Creating all the wikibase items related to the paper is a two pass process, due to the fact that
	// the virtual data structure that is described in [0] and related examples has two way links between
	// items (e.g., an Anchor Node item references an Annotation item, and that Annotation item needs to
	// refer to the Anchor Node).
	//
	// So to simplify the logic we only add properties to items once we have created all the items, as that's
	// the only time when we have all the information about all properties for each item.
	//
	// [0] https://sciencesource.wmflabs.org/wiki/Data_schema
	upload_err := sciSourceClient.CreateArticleItemTree(processor.ScienceSourceRecord)
	// regardless of whether we error, do another save to record any partial changes to the tree
	err = processor.ScienceSourceRecord.Save(processor.targetScienceSourceStateFileName())
	if err != nil || upload_err != nil {

		// if we had two errors combine them into one
		if err != nil && upload_err != nil {
			err = fmt.Errorf("Failed to both create wikibase items (%v) and save state (%v)", upload_err, err)
		} else if upload_err != nil {
			err = upload_err
		}

		return err
	}

	// If we got here then now we have an item for every part of the data structure, so upload all the properties.
	err = sciSourceClient.ReconsileArticleItemTree(processor.ScienceSourceRecord)
	if err != nil {
	    return err
	}
    err = processor.ScienceSourceRecord.Save(processor.targetScienceSourceStateFileName())
    if err != nil {
        return err
    }

	return nil
}
