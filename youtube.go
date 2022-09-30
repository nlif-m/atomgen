package main

import (
	// "fmt"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// TODO: pass YTDLP arguments as some type of struct for example
// to avoid using global variables
type ytdlpChannels struct {
	in  chan string
	out chan string
}

func downloadChannelAsAudio(chs ytdlpChannels) {
	downloadURL := <-chs.in
	cmd := exec.Command(ytdlp, "--playlist-end", "10", "--dateafter", fmt.Sprint("today-", howManyWeeksDownload, "weeks"), "-x", "--download-archive", ytdlpDownloadArchive, "-f",
		"bestaudio", "-o", ytdlpOutputTemplate, "--no-simulate", "-O", "Downloading %(title)s", "--no-progress", downloadURL)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	cOut, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Warning:  failed to run ", "[", cmd, "]", err, string(cOut))
	}

	chs.out <- downloadURL
}

// TODO: there exists some type of problem if same channel url exist >1 times, but not critical
// Make a set from records, to prevent mistaken double urls
func downloadVideosFromFile(file string) {
	log.Println("Start downloading videos from urls in", file)
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
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

	for range records {
		go downloadChannelAsAudio(*chs)
	}

	for _, record := range records {
		source := record[0]
		chs.in <- source
	}

	for index := range records {
		log.Printf("[%d/%d] Download %s\n", index+1, length, <-chs.out)
	}
	log.Println("Finished downloading videos from urls in", file)
}
