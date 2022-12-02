package main

import (
	// "fmt"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type ytdlp struct {
	programName string
}

func newYtdlp(programName string) ytdlp {
	return ytdlp{programName: programName}
}

func (yt *ytdlp) newCmdWithArgs(Args ...string) *exec.Cmd {
	return exec.Command(yt.programName, Args...)

}

func (yt *ytdlp) GetChannelNameFromURL(URL string) (channelName string, err error) {
	cmd := yt.newCmdWithArgs(
		"-O", "%(channel)s",
		"--playlist-end", "1",
		URL)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Warning:  failed to run ", "[", cmd, "]", err)
		return "", err
	}
	return string(cmdOutput), nil
}

func (yt *ytdlp) GetVersion() (version string, err error) {
	cmd := yt.newCmdWithArgs("--version")
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Warning:  failed to run ", "[", cmd, "]", err)
		return "", err
	}
	return string(cmdOutput), nil
}

func (yt *ytdlp) DownloadURLAsAudio(URL string) error {
	log.Println("Start downloading:", URL)
	cmd := yt.newCmdWithArgs(
		"--playlist-end", "10",
		"--dateafter", fmt.Sprint("today-", howManyWeeksIsOld, "weeks"),
		"-x",
		"--download-archive", ytdlpDownloadArchive,
		"-f", "bestaudio",
		"-o", ytdlpOutputTemplate,
		"--no-simulate", "-O", "Downloading %(title)s",
		URL)

	err := cmd.Run()
	if err != nil {
		log.Println("Warning:  failed to run ", "[", cmd, "]", err)
		return err
	}
	log.Println("Finish downloading:", URL)
	return nil
}

func (yt *ytdlp) DownloadVideosFromFile(file string) {
	log.Println("Start downloading videos from urls in", file)
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	// TODO: Think about chnaging csv to another format
	// maybe json
	r := csv.NewReader(fd)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	recordsSet := map[string]struct{}{}
	for _, record := range records {
		recordsSet[record[0]] = struct{}{}
	}

	var wg sync.WaitGroup
	for record := range recordsSet {
		wg.Add(1)
		go func(URL string) {
			yt.DownloadURLAsAudio(URL)
			wg.Done()
		}(record)
	}

	wg.Wait()
	log.Println("Finished downloading videos from urls in", file)
}
