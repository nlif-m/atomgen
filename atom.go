package main

import (
	"encoding/xml"
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
)

type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Length  string   `xml:"length,attr"`
	Src     string   `xml:"src,attr"`
	Type    string   `xml:"type,attr"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	Link    Link     `xml:"link"`
	Id      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Summary string   `xml:"summary"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
}

type Atom struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	Link    Link     `xml:"link"`
	Updated string   `xml:"updated"`
	Author  Author   `xml:"author"`
	Id      string   `xml:"id"`
	Entries []Entry  `xml:"entry"`
}

func generateAtomRssFile(rssFile string, src_folder string) {
	log.Println("Start generating ", ATOM_FILE)
	files, err := ioutil.ReadDir(src_folder)
	if err != nil {
		panic(err)
	}

	var entries = make([]Entry, 0, 10)
	files, err = ioutil.ReadDir(src_folder)
	if err != nil {
		panic(err)
	}
	entries_count := 0
	for _, file := range files {
		Name := file.Name()

		Ext := filepath.Ext(Name)
		if Ext == ".part" || Ext == ".ytdl" {
			continue
		}

		if Name == filepath.Base(YTDLP_DOWNLOAD_ARCHIVE) {
			continue
		}

		// TODO: I am sure there something wrong, but what ?
		fd, err := os.Open(path.Join(src_folder, file.Name()))
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
		fileLoc := path.Join(CHANNEL_LINK, src_folder, urlEncodedName)
		fileTime := file.ModTime().Format(time.RFC3339)
		// log.Printf("Generated entry for '%s': '%s' %s\n", ATOM_FILE, Name, fileTime)

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
		entries_count += 1
	}

	log.Printf("Generated %d entries for '%s'\n", entries_count, ATOM_FILE)
	v := &Atom{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   CHANNEL_TITLE,
		Link:    Link{Href: CHANNEL_LINK},
		Updated: time.Now().Format(time.RFC3339),
		Author:  Author{Name: "rss.yasal.xyz"},
		Id:      CHANNEL_LINK,
		Entries: entries}

	data, err := xml.MarshalIndent(v, " ", "  ")

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(rssFile, data, 0600)

	if err != nil {
		panic(err)
	}
	log.Println("Finish generating ", ATOM_FILE)
}
