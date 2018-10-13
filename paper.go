// Copyright ContentMine Ltd 2018

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
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
| wikiedata = %s
| main_subject = %s
| DOI =
| PubMed_ID =
| PMC_ID =
}}
`

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
	)

	f.Write([]byte(header))

	cmd := exec.Cmd{
		Path: "/usr/bin/xsltproc",
		Args: []string{"xsltproc", "jats-html.xsl", processor.targetXMLFileName()},
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

// main entry point

func (processor PaperProcessor) ProcessPaper() error {

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

	return nil
}
