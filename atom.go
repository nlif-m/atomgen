package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/tools/blog/atom"
)

func getMimeType(filePath string) (mimeType string, err error) {
	// TODO: I am sure there something wrong, but what ?
	fd, err := os.Open(filePath)
	if err != nil {
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
			{
				Rel:    string("enclosure"),
				Href:   fileLocation,
				Type:   mimeType,
				Length: uint(length),
			},
		},
		Published: atom.Time(fileModificationTime),
		Updated:   atom.Time(fileModificationTime),
		Summary: &atom.Text{
			Type: "text",
			Body: name}}
}

func newAtomFeed(channelTitle string, channelLink string, authorLink string, entries []*atom.Entry) *atom.Feed {
	return &atom.Feed{
		XMLName: xml.Name{},
		Title:   channelTitle,
		Link: []atom.Link{
			{Href: channelLink}},
		Updated: atom.Time(time.Now()),
		Author:  &atom.Person{Name: channelLink},
		ID:      channelLink,
		Entry:   entries}
}
