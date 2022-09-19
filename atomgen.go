package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

const (
	DetectContentTypeMost = 512

	YTDLP                  = "yt-dlp"
	SRC_FOLDER             = "src"
	URLS_CSV_FILE          = "urls.csv"
	YTDLP_DOWNLOAD_ARCHIVE = SRC_FOLDER + string(os.PathSeparator) + "downloaded.txt"
	YTDLP_OUTPUT_TEMPLATE  = SRC_FOLDER + string(os.PathSeparator) + "%(uploader)s_%(title)s.%(ext)s"

	ATOM_FILE = "rss.xml"

	CHANNEL_TITLE = "test page"
	CHANNEL_LINK  = "https://rss.yasal.xyz"
)

// TODO: pass YTDLP arguments as some type of struct for example
// to avoid using global variables
type ytdlp_channels struct {
	in  chan string
	out chan string
}

func downloadChannelAsAudio(chs ytdlp_channels, length int) {
	download_url := <-chs.in
	cmd := exec.Command(YTDLP, "--playlist-end", "10", "--dateafter", "today-4weeks", "-x", "--download-archive", YTDLP_DOWNLOAD_ARCHIVE, "-f",
		"bestaudio", "-o", YTDLP_OUTPUT_TEMPLATE, "--no-simulate", "-O", "Downloading %(title)s", "--no-progress", download_url)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	chs.out <- download_url

}

func downloadVideosFromFile(file string) {
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(fd)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	length := len(records)
	chs := &ytdlp_channels{
		make(chan string),
		make(chan string),
	}

	for _, record := range records {
		source := record[0]
		go downloadChannelAsAudio(*chs, length)
		chs.in <- source
	}

	for index := range records {
		log.Printf("[%d/%d] Download %s\n", index+1, length, <-chs.out)
	}
}

func generateAtomRssFile(rssFile string, src_folder string) {

	files, err := ioutil.ReadDir(src_folder)
	if err != nil {
		panic(err)
	}

	var entries = make([]Entry, 0, 10)
	files, err = ioutil.ReadDir(src_folder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		Name := file.Name()

		Ext := filepath.Ext(Name)
		if Ext == ".part" || Ext == ".ytdl" {
			continue
		}

		if Name == filepath.Base(YTDLP_DOWNLOAD_ARCHIVE) {
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
		fileLoc := CHANNEL_LINK + "/" + src_folder + "/" + urlEncodedName
		fileTime := file.ModTime().Format(time.RFC3339)
		log.Printf("Generated entry for '%s': '%s' %s\n", ATOM_FILE, Name, fileTime)
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
}

func main() {
	log.Println("Start downloading videos from urls in", URLS_CSV_FILE)
	downloadVideosFromFile(URLS_CSV_FILE)
	log.Println("Finished downloading videos from urls in", URLS_CSV_FILE)

	log.Println("Start generating ", ATOM_FILE)
	generateAtomRssFile(ATOM_FILE, SRC_FOLDER)
	log.Println("Finish generating ", ATOM_FILE)
}
