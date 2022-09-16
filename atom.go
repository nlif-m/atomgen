package main

import (
	"encoding/xml"
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
