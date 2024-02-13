package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: Separate extraction and counting into separate functions.
func extractAndCountWords(odtFilePath string) int {
	zipReader, err := zip.OpenReader(odtFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == "content.xml" {
			fileReader, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer fileReader.Close()

			decoder := xml.NewDecoder(fileReader)
			totalWords := 0
			inParagraph := false

			for {
				token, _ := decoder.Token()
				if token == nil {
					break
				}

				switch startElement := token.(type) {
				case xml.StartElement:
					if startElement.Name.Local == "p" && startElement.Name.Space == "urn:oasis:names:tc:opendocument:xmlns:text:1.0" {
						inParagraph = true
					}
				case xml.EndElement:
					if startElement.Name.Local == "p" {
						inParagraph = false
					}
				case xml.CharData:
					if inParagraph {
						words := strings.Fields(string(startElement))
						totalWords += len(words)
					}
				}
			}

			return totalWords
		}
	}

	return 0
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get home directory: ", err)
	}

	// TODO: Enable counting of all files of type in directory.
	odtFilePath := filepath.Join(homeDir, "mywriting", "UntitledDoc2-2.odt")
	wordCount := extractAndCountWords(odtFilePath)
	fmt.Printf("Total word count: %d\n", wordCount)
}
