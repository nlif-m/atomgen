package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func GetRidOfWrongCharacters(filename string) (rightName string) {
OuterLoop:
	for _, char := range filename {
		switch char {
		case '&':
			rightName += "_and_"
			continue OuterLoop
		case ' ':
			rightName += "_"
			continue OuterLoop
		default:
			rightName += string(char)

		}
	}
	return rightName
}

func getLinkType(filename string) (linkType string) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp3", ".opus":
		return "audio/mpeg"
	case ".mp4":
		return "video/wmv"
	default:
		return "UNKNOWN"
	}
}

func main() {
	channelTitle := "test page"
	channelLink := "https://rss.yasal.xyz"
	//	channelDesc := "Hello"
	pubDate := time.Now().UTC().String()

	rssFile := "rss.xml"
	srcFolder := "src"

	files, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		panic(err)
	}
	fmt.Println(GetRidOfWrongCharacters("Hello world"))
	var entries = make([]Entry, 0, 10)
	for _, file := range files {
		oldName := srcFolder + "/" + file.Name()
		rightName := GetRidOfWrongCharacters(file.Name())
		newName := srcFolder + "/" + rightName
		if oldName != newName {
			os.Rename(oldName, newName)
			fmt.Printf("Rename '%s' to '%s'\n", oldName, newName)			
		}
	}

	files, err = ioutil.ReadDir(srcFolder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fileLoc := channelLink + "/" + srcFolder + "/" + file.Name()
		fmt.Printf("Generated entry for '%s'\n", file.Name())
		entries = append(entries, Entry{Title: file.Name(),
			Link: Link{Rel: "enclosure", Length: strconv.Itoa(int(file.Size())),
				Type: getLinkType(file.Name()),
				Href: fileLoc},
			Id:      fileLoc,
			Updated: file.ModTime().String(),
			Summary: file.Name()})
	}

	v := &Atom{Xmlns: "http://www.w3.org/2005/Atom",
		Title:   channelTitle,
		Link:    Link{Href: channelLink},
		Updated: pubDate,
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

	fmt.Printf("Write entries based on files in '%s' to '%s'\n", srcFolder, rssFile)
}
