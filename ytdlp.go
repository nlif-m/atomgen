package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
		log.Printf("Warning: failed to get channel name for '%s'\n cmd: [%s]\t%s\n", URL, cmd, err)
		return "", err
	}
	return string(cmdOutput), nil
}

func (yt *ytdlp) GetVersion() (version string, err error) {
	cmd := yt.newCmdWithArgs("--version")
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: failed to get version of '%s'\n cmd: [%s]\t%s\n", yt.programName, cmd, err)
		return "", err
	}
	return string(cmdOutput), nil
}

func (yt *ytdlp) DownloadURLAsAudio(URL string) error {
	channelName, _ := yt.GetChannelNameFromURL(URL)
	channelName = strings.TrimSpace(channelName)
	log.Printf("Start downloading: %v\t%v\n", channelName, URL)
	var cmd *exec.Cmd

	ytdlpOutputTemplate := srcFolder + string(os.PathSeparator) + "%(uploader)s %(title)s.%(ext)s"
	// TODO: looks very bad, fix it
	if weeksToDownload == 0 {
		cmd = yt.newCmdWithArgs(
			"--playlist-items", fmt.Sprintf("0:%v", videosToDownload),
			"-x",
			"--download-archive", ytdlpDownloadArchive,
			"-f", "bestaudio",
			"-o", ytdlpOutputTemplate,
			"--no-simulate", "-O", "Downloading %(title)s",
			URL)
	} else {
		cmd = yt.newCmdWithArgs(
			"--playlist-items", fmt.Sprintf("0:%v", videosToDownload),
			"--dateafter", fmt.Sprint("today-", weeksToDownload, "weeks"),
			"-x",
			"--download-archive", ytdlpDownloadArchive,
			"-f", "bestaudio",
			"-o", ytdlpOutputTemplate,
			"--no-simulate", "-O", "Downloading %(title)s",
			URL)
	}

	body, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: failed to download '%s' as audio\n cmd: [%s]\t%s\n", URL, cmd, err)
		return err
	}
	log.Printf("Finish downloading: %v\t%v\n%s\n", channelName, URL, string(body))
	return nil
}

func (yt *ytdlp) DownloadVideosFromFile(file string) {
	log.Printf("Start downloading videos from urls in '%s'\n", file)
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
	log.Printf("Finished downloading videos from urls in '%s'\n", file)
}
