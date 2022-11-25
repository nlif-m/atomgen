package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/tools/blog/atom"
)

func getMimeType(filePath string) (mimeType string, err error) {
	// TODO: I am sure there something wrong, but what ?
	fd, err := os.Open(filePath)
	if err != nil {
		defer fd.Close()
		return "", err
	}
	defer fd.Close()

	r := io.Reader(fd)
	r1 := io.LimitReader(r, detectContentTypeMost)
	head, err := io.ReadAll(r1)

	if err != nil {
		return "", err
	}

	mimeType = http.DetectContentType(head)
	return mimeType, nil
}
func newAtomEntry(name string, fileLocation string, mimeType string, length uint, fileModificationTime time.Time) *atom.Entry {
	return &atom.Entry{
		Title: name,
		ID:    fileLocation,
		Link: []atom.Link{
			atom.Link{
				Rel:    string("enclosure"),
				Href:   fileLocation,
				Type:   mimeType,
				Length: uint(length),
			},
		},
		Updated: atom.Time(fileModificationTime),
		Summary: &atom.Text{
			Type: "text",
			Body: name}}
}

func getEntriesFromSrcFolder(srcFolder string) (entries []*atom.Entry, err error) {
	files, err := os.ReadDir(srcFolder)
	if err != nil {
		log.Printf("ERROR: when getting entries from %s", srcFolder)
		return nil, err
	}
	entries = make([]*atom.Entry, 0, 100)
filesLoop:
	for _, file := range files {
		Name := file.Name()

		Ext := filepath.Ext(Name)
		switch Ext {
		case ".part", ".ytdl":
			continue filesLoop
		}

		switch Name {
		case filepath.Base(ytdlpDownloadArchive):
			continue filesLoop
		}

		mimeType, err := getMimeType(srcFolder + string(os.PathSeparator) + file.Name())
		if err != nil {
			log.Printf("ERROR: while getting Mimetype of %s%c%s\n%s", srcFolder, os.PathSeparator, file.Name(), err)
			return nil, err
		}

		urlEncodedName := url.PathEscape(Name)
		fileLocation := channelLink + string(os.PathSeparator) + srcFolder + string(os.PathSeparator) + urlEncodedName
		fileInfo, err := file.Info()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		fileModificationTime := fileInfo.ModTime()
		length, err := strconv.ParseUint((strconv.Itoa(int(fileInfo.Size()))), 10, 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// log.Printf("Generated entry for '%s': '%s' %s\n", ATOM_FILE, Name, fileTime)

		entries = append(entries, newAtomEntry(Name, fileLocation, mimeType, uint(length), fileModificationTime))
	}
	return entries, nil
}

func newAtomFeed(channelTitle string, channelLink string, authorName string, entries []*atom.Entry) *atom.Feed {
	return &atom.Feed{
		XMLName: xml.Name{},
		Title:   channelTitle,
		Link: []atom.Link{
			atom.Link{Href: channelLink}},
		Updated: atom.Time(time.Now()),
		Author:  &atom.Person{Name: channelLink},
		ID:      channelLink,
		Entry:   entries}
}
func generateAtomRssFile(rssFile string, srcFolder string) error {
	log.Println("Start generating ", rssFile)

	entries, err := getEntriesFromSrcFolder(srcFolder)
	if err != nil {
		return err
	}
	log.Printf("Generated %d entries for '%s'\n", len(entries), rssFile)
	v := newAtomFeed(channelTitle, channelLink, channelLink, entries)
	data, err := xml.MarshalIndent(v, " ", "  ")

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(rssFile, data, 0600)

	if err != nil {
		log.Fatalf("ERROR: failed to write rss file to %s %s", rssFile, err)
	}
	log.Println("Finish generating ", rssFile)
	return nil
}
