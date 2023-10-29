package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	aatom "github.com/nlif-m/atomgen/atom"
	"github.com/nlif-m/atomgen/config"
	"github.com/nlif-m/atomgen/utils"
	"github.com/nlif-m/atomgen/ytdlp"

	"golang.org/x/tools/blog/atom"
)

type Atomgen struct {
	ytdlp ytdlp.Ytdlp
	cfg   config.Cfg
}

func newAtomgen(ytdlp ytdlp.Ytdlp, cfg config.Cfg) Atomgen {
	return Atomgen{ytdlp, cfg}
}

func (atomgen *Atomgen) fullUpdate() error {
	if atomgen.cfg.VideosToDowload != 0 {
		err := atomgen.DownloadVideos()
		return err
	}

	if atomgen.cfg.WeeksToDelete != 0 {
		err := atomgen.deleteOldFiles()
		return err
	}

	return atomgen.generateAtomFeed()
}
func (atomgen *Atomgen) generateAtomFeed() error {
	entries, err := atomgen.getEntries()
	if err != nil {
		return err
	}
	atomFeed := aatom.NewFeed(atomgen.cfg.ChannelTitle, atomgen.cfg.AuthorLink, atomgen.cfg.AuthorLink, entries)

	log.Printf("Generated %d entries for '%s'\n", len(entries), atomgen.cfg.AtomFile)
	data, err := xml.MarshalIndent(atomFeed, " ", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(atomgen.cfg.AtomFile, data, 0644)

	if err != nil {
		return err
	}
	log.Printf("Finish generating: '%s'\n", atomgen.cfg.AtomFile)
	return nil
}

func (atomgen *Atomgen) getEntries() (entries []*atom.Entry, err error) {
	files, err := os.ReadDir(atomgen.cfg.SrcFolder)
	if err != nil {
		log.Printf("ERROR: when getting entries from %s", atomgen.cfg.SrcFolder)
		return nil, err
	}
	entries = make([]*atom.Entry, 0, len(files))
filesLoop:
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		Name := file.Name()

		Ext := filepath.Ext(Name)
		switch Ext {
		case ".part", ".ytdl", ".xml", ".json":
			continue filesLoop
		}

		switch Name {
		case filepath.Base(atomgen.cfg.YtdlpDownloadArchive):
			continue filesLoop
		}

		mimeType, err := aatom.GetMimeType(filepath.Join(atomgen.cfg.SrcFolder, file.Name()))
		if err != nil {
			log.Printf("ERROR: while getting Mimetype of %s%c%s\n%s", atomgen.cfg.SrcFolder, os.PathSeparator, file.Name(), err)
			return nil, err
		}

		urlEncodedName := url.PathEscape(Name)
		fileLocation, err := url.JoinPath(atomgen.cfg.LocationLink, urlEncodedName)
		utils.CheckErr(err)
		fileInfo, err := file.Info()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		fileModificationTime := fileInfo.ModTime()
		length, err := strconv.ParseUint((strconv.Itoa(int(fileInfo.Size()))), 10, 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		entries = append(entries, aatom.NewEntry(Name, fileLocation, mimeType, uint(length), fileModificationTime, Name))
	}
	return entries, nil
}

func (atomgen *Atomgen) deleteOldFiles() error {
	log.Println("Start deleting old videos")
	files, err := os.ReadDir(atomgen.cfg.SrcFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(atomgen.cfg.SrcFolder, file.Name())
		fileInfo, err := file.Info()

		if err != nil {
			log.Println("Warning: failed to delete file at ", filePath, err)
			continue
		}
		if !fileInfo.ModTime().Before(time.Now().AddDate(0, 0, -int(atomgen.cfg.WeeksToDelete)*7)) {
			continue
		}

		log.Println(fmt.Sprint("Deleting file older than ", atomgen.cfg.WeeksToDelete, " weeks "), file.Name())
		err = os.Remove(filePath)
		if err != nil {
			log.Println("Warning: failed to delete file at ", filePath, err)
			continue
		}
	}

	log.Println("Finish deleting old videos")
	return nil
}

func (atomgen *Atomgen) DownloadURL(URL string, withoutTimeLimit bool, usingDownloadArchive bool) error {
	channelName, err := atomgen.ytdlp.GetChannelNameFromURL(URL)
	if err != nil {
		return err
	}
	channelName = strings.TrimSpace(channelName)
	log.Printf("Start downloading: %v\t%v\n", channelName, URL)
	var cmd *exec.Cmd

	ytdlpOutputTemplate := filepath.Join(atomgen.cfg.SrcFolder, "%(uploader)s %(title)s.%(ext)s")
	cmd = atomgen.ytdlp.NewCmdWithArgs(
		"--playlist-items", fmt.Sprintf("0:%v", atomgen.cfg.VideosToDowload),
		"-x",
		"--match-filters", fmt.Sprintf("!is_live & duration>%d", atomgen.cfg.YtdlpDurationLowerLimit),
		"-f", "ba/ba*",
		"--audio-format", fmt.Sprintf("%s/best", atomgen.cfg.DownloadAudioFormat),
		"-o", ytdlpOutputTemplate,
		"--no-simulate")

	if usingDownloadArchive {
		cmd.Args = append(cmd.Args, "--download-archive", atomgen.cfg.YtdlpDownloadArchive)
	}
	if !withoutTimeLimit && !(atomgen.cfg.WeeksToDownload == 0) {
		cmd.Args = append(cmd.Args, "--dateafter", fmt.Sprint("today-", atomgen.cfg.WeeksToDownload, "weeks"))
	}
	cmd.Args = append(cmd.Args, URL)

	body, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: failed to download '%s' as audio\n cmd: [%s]\t%s\n%s\n", URL, cmd, err, string(body))
		return err
	}
	log.Printf("Finish downloading: %v\t%v\n%s\n", channelName, URL, string(body))
	return nil
}

func (atomgen *Atomgen) DownloadVideos() error {
	log.Printf("Start downloading videos to '%s'\n", atomgen.cfg.SrcFolder)
	records := atomgen.cfg.Urls
	records = utils.Unique(records)
	isNotEmpty := func(record string) bool {
		return record != ""
	}
	records = utils.Filter(records, isNotEmpty)

	var wg sync.WaitGroup

	limitDownloadBuffer := make(chan int, atomgen.cfg.LimitDownload)

	// TODO: fix that for each record goroutine is createad but downloading not starting because of buffering
	// suggest to don't create goroutine until it needed

	for _, record := range records {
		wg.Add(1)
		go func(URL string) {
			defer wg.Done()
			limitDownloadBuffer <- 1
			atomgen.DownloadURL(URL, false, true)
			<-limitDownloadBuffer
		}(record)
	}

	wg.Wait()
	log.Printf("Finished downloading videos to '%s'\n", atomgen.cfg.SrcFolder)
	return nil
}
