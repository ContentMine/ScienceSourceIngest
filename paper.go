// Copyright ContentMine Ltd 2018

package main

import (
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


// Side effect heavy functions

func (processor PaperProcessor) createFolderIfRequired() error {
	return os.MkdirAll(processor.folderName(), 0755)
}

func (processor PaperProcessor) fetchPaperTextToDisk() error {

	// if it already exists, don't fetch it again
	if _, err := os.Stat(processor.targetXMLFileName()); err == nil {
		return nil
	}

	f, err := os.Create(processor.targetXMLFileName())
	if err != nil {
		return err
	}
	defer f.Close()

	resp, resp_err := http.Get(processor.Paper.FullTextURL())
	if resp_err != nil {
		return resp_err
	}
	defer resp.Body.Close()

	_, copy_err := io.Copy(f, resp.Body)
	return copy_err
}

func (processor PaperProcessor) processXMLToHTML() error {

	f, err := os.Create(processor.targetHTMLFileName())
	if err != nil {
		return err
	}
	defer f.Close()

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

	err = processor.processXMLToHTML()
	if err != nil {
		return err
	}

	return nil
}
