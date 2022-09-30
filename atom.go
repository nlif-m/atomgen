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
	"path"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/tools/blog/atom"
)

func generateAtomRssFile(rssFile string, srcFolder string) {
	log.Println("Start generating ", atomFile)
	files, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		log.Fatal(err)
	}

	var entries = make([]*atom.Entry, 0, 10)
	files, err = ioutil.ReadDir(srcFolder)
	if err != nil {
		log.Fatal(err)
	}
	entriesCount := 0
	for _, file := range files {
		Name := file.Name()

		Ext := filepath.Ext(Name)
		log.Println(Name, Ext)
		switch Ext {
		case ".part", ".ytdl":
			continue
		}

		if Name == filepath.Base(ytdlpDownloadArchive) {
			continue
		}

		if file.ModTime().Before(time.Now().AddDate(0, 0, -howManyWeeksDownload*7)) {
			filePath := srcFolder + string(os.PathSeparator) + file.Name()
			log.Println(fmt.Sprint("Deleting file older than ", howManyWeeksDownload, " weeks "), filePath)
			err = os.Remove(filePath)
			if err != nil {
				log.Println("Warning: failed to delete file at", filePath, err)
			}
			continue
		}

		// TODO: I am sure there something wrong, but what ?
		fd, err := os.Open(path.Join(srcFolder, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		r := io.Reader(fd)
		r1 := io.LimitReader(r, detectContentTypeMost)
		head, err := io.ReadAll(r1)

		if err != nil {
			log.Fatal(err)
		}

		mimeType := http.DetectContentType(head)

		urlEncodedName := url.PathEscape(Name)
		fileLoc := channelLink + string(os.PathSeparator) + srcFolder + string(os.PathSeparator) + urlEncodedName
		fileTime := file.ModTime()
		length, err := strconv.ParseUint((strconv.Itoa(int(file.Size()))), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("Generated entry for '%s': '%s' %s\n", ATOM_FILE, Name, fileTime)

		entries = append(entries,
			&atom.Entry{
				Title: Name,
				ID:    fileLoc,
				Link: []atom.Link{
					atom.Link{
						Rel:    string("enclosure"),
						Href:   fileLoc,
						Type:   mimeType,
						Length: uint(length),
					},
				},
				Updated: atom.Time(fileTime),
				Summary: &atom.Text{
					Type: "text",
					Body: Name}},
		)
		entriesCount++
	}

	log.Printf("Generated %d entries for '%s'\n", entriesCount, atomFile)
	v := &atom.Feed{
		XMLName: xml.Name{},
		Title:   channelTitle,
		Link: []atom.Link{
			atom.Link{Href: channelLink}},
		Updated: atom.Time(time.Now()),
		Author:  &atom.Person{Name: "rss.yasal.xyz"},
		ID:      channelLink,
		Entry:   entries}
	data, err := xml.MarshalIndent(v, " ", "  ")

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(rssFile, data, 0600)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Finish generating ", atomFile)
}
