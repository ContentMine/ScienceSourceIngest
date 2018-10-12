// Copyright ContentMine Ltd 2018

package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

type PaperProcessor struct {
	Paper           Paper
	TargetDirectory string
}

func (processor PaperProcessor) folderName() string {
	return path.Join(processor.TargetDirectory, processor.Paper.ID())
}

func (processor PaperProcessor) targetXMLFileName() string {
	return path.Join(processor.folderName(), "paper.xml")
}

func (processor PaperProcessor) createFolderIfRequired() error {
	return os.MkdirAll(processor.folderName(), 0755)
}

func (processor PaperProcessor) fetchPaperTextToDisk() error {
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

func (processor PaperProcessor) ProcessPaper() error {

	err := processor.createFolderIfRequired()
	if err != nil {
		return err
	}

	err = processor.fetchPaperTextToDisk()
	if err != nil {
		return err
	}

	return nil
}
