package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
)

func getLinkType(filename string) (linkType string) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp3", ".m4a":
		return "audio/mpeg"
	case ".opus":
		return "audio/ogg"
	case ".mp4":
		return "video/wmv"
	default:
		return "UNKNOWN"
	}
}

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

		urlEncodedName := url.PathEscape(Name)
		fileLoc := channelLink + "/" + srcFolder + "/" + urlEncodedName
		fmt.Printf("Generated entry for '%s' %s\n", Name, file.ModTime().String())
		entries = append(entries, Entry{Title: Name,
			Link: Link{Rel: "enclosure", Length: strconv.Itoa(int(file.Size())),
				Type: getLinkType(Name),
				Href: fileLoc},
			Id:      fileLoc,
			Updated: file.ModTime().String(),
			Summary: Name})
	}

	v := &Atom{Xmlns: "http://www.w3.org/2005/Atom",
		Title: channelTitle,
		Link:  Link{Href: channelLink},
		//		Updated: pubDate,
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
