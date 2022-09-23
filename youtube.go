package main

import (
	// "fmt"
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// TODO: pass YTDLP arguments as some type of struct for example
// to avoid using global variables
type ytdlpChannels struct {
	in  chan string
	out chan string
}

func downloadChannelAsAudio(chs ytdlpChannels, length int) {
	downloadURL := <-chs.in
	cmd := exec.Command(ytdlp, "--playlist-end", "10", "--dateafter", "today-4weeks", "-x", "--download-archive", ytdlpDownloadArchive, "-f",
		"bestaudio", "-o", ytdlpOutputTemplate, "--no-simulate", "-O", "Downloading %(title)s", "--no-progress", downloadURL)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	chs.out <- downloadURL
}

func downloadVideosFromFile(file string) {
	log.Println("Start downloading videos from urls in", urlsFile)
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
	chs := &ytdlpChannels{
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
	log.Println("Finished downloading videos from urls in", urlsFile)
}

func getYoutubeXml(url string) (body []byte, err error) {
	resp, err := http.Get("https://www.youtube.com/feeds/videos.xml?channel_id=UCoZiaJ2Yks-XzhyXoaNe3Dw")

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}
	return
}
