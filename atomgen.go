package main

// TODO: make a daemon that run program eveny N time
// TODO: Make a separate commands, for example only 'atomgen download' to download and nothing else and etc.

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	// "flag"
)

const (
	detectContentTypeMost = 512

	ytdlp                = "yt-dlp"
	src_folder           = "src"
	urlsFile             = "urls.csv"
	ytdlpDownloadArchive = src_folder + string(os.PathSeparator) + "downloaded.txt"
	ytdlpOutputTemplate  = src_folder + string(os.PathSeparator) + "%(uploader)s %(title)s.%(ext)s"

	atomFile = "rss.xml"

	channelTitle = "test page1"
	channelLink  = "https://rss.yasal.xyz"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hi, check <a href=/rss.xml >rss.xml</a>"))
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	generateAtomRssFile(atomFile, src_folder)
	fd, err := os.Open(atomFile)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	defer fd.Close()

	body, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
}

var srcPath = regexp.MustCompile("^/(src)/(.*)$")

func srcHandler(w http.ResponseWriter, r *http.Request) {

	t := srcPath.FindStringSubmatch(r.URL.Path)
	name := t[2]
	path := src_folder + string(os.PathSeparator) + name
	fd, err := os.Open(path)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	defer fd.Close()

	body, err := ioutil.ReadAll(fd)
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
}

func main() {

	http.HandleFunc("/rss.xml/", rssHandler)
	http.HandleFunc("/src/", srcHandler)
	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe("192.168.43.39:8080", nil))
	// 	downloadVideosFromFile(URLS_CSV_FILE)
}
