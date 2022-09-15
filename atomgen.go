package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	DetectContentTypeMost = 512
)

func main() {
	channelTitle := "test page"
	channelLink := "https://rss.yasal.xyz"

	ytDlAlreadyDownloadFile := "downloaded.txt"

	rssFile := "rss.xml"
	srcFolder := "src"

	files, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		panic(err)
	}

	var entries = make([]Entry, 0, 10)
	// for _, file := range files {
	// 	oldName := srcFolder + "/" + file.Name()

	// 	oldExt := filepath.Ext(oldName)
	// 	if oldExt == ".part" {
	// 		continue
	// 	}
	// 	rightName := GetRidOfWrongCharacters(file.Name())
	// 	newName := srcFolder + "/" + rightName
	// 	if oldName != newName {
	// 		err = os.Rename(oldName, newName)
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "ERROR: when renaming %s to %s", oldName, newName)
	// 			panic(err)
	// 		}
	// 		fmt.Printf("Rename '%s' to '%s'\n", oldName, newName)
	// 	}
	// }

	files, err = ioutil.ReadDir(srcFolder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		Name := file.Name()

		Ext := filepath.Ext(Name)
		if Ext == ".part" {
			continue
		}

		if Name == ytDlAlreadyDownloadFile {
			fmt.Printf("%s \n", Name)
			continue
		}

		// TODO: I am sure there something wrong, but what ?
		fd, err := os.Open("src/" + file.Name())
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		r := io.Reader(fd)
		r1 := io.LimitReader(r, DetectContentTypeMost)
		head, err := io.ReadAll(r1)

		if err != nil {
			panic(err)
		}

		mimeType := http.DetectContentType(head)

		urlEncodedName := url.PathEscape(Name)
		fileLoc := channelLink + "/" + srcFolder + "/" + urlEncodedName
		fileTime := file.ModTime().Format(time.RFC3339)
		fmt.Printf("Generated entry for '%s' %s\n", Name, fileTime)
		entries = append(entries,
			Entry{
				Title: Name,
				Link: Link{
					Rel: "enclosure", Length: strconv.Itoa(int(file.Size())),
					Type: mimeType,
					Href: fileLoc},
				Id:      fileLoc,
				Updated: fileTime,
				Summary: Name})
	}

	v := &Atom{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   channelTitle,
		Link:    Link{Href: channelLink},
		Updated: time.Now().Format(time.RFC3339),
		Author:  Author{Name: "rss.yasal.xyz"},
		Id:      channelLink,
		Entries: entries}

	data, err := xml.MarshalIndent(v, " ", "  ")

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(rssFile, data, 0600)

	if err != nil {
		panic(err)

	}
}
