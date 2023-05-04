package main

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

const (
	// Cfg
	ytdlpProgramDefault         string = "yt-dlp"
	outputFolderDefault         string = "/tmp/test"
	srcFolderDefault            string = "src"
	urlsFileDefault             string = "urls.csv"
	atomFileDefault             string = "atom.xml"
	channelTitleDefault         string = "test page"
	authorLinkDefault           string = "https://rss.yasal.xyz"
	locationLinkDefault         string = "https://rss.yasal.xyz" // .../src
	ytdlpDownloadArchiveDefault string = "downloaded.txt"
	weeksToDownloadDefault      uint   = 4
	weeksToDeleteDefault        uint   = 0
	videosToDownloadDefault     int    = 10
	generateAtomFileDefault     bool   = true
	limitDownloadDefault        uint   = 10
	downloadAudioFormatDefault  string = "mp3"

	// MimeDetect
	detectContentTypeMost = 512

	// LocationType
	HttpLocation LocationType = "http"
	S3Location   LocationType = "s3" // TODO: implement s3 support
)

type LocationType string

type Cfg struct {
	AtomFile             string
	AuthorLink           string
	ChannelTitle         string
	LocationLink         string
	LocationType         LocationType
	Urls                 []string
	VideosToDowload      int
	WeeksToDelete        uint
	WeeksToDownload      uint
	YtdlpDownloadArchive string
	YtdlpProgram         string
	SrcFolder            string
	OutputFolder         string
	LimitDownload        uint
	DownloadAudioFormat  string
}

func newCfgFromFile(filePath string) (Cfg, error) {
	newCfg := Cfg{}
	body, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read '%s' \n", filePath)
		return Cfg{}, err
	}

	err = json.Unmarshal(body, &newCfg)
	if err != nil {
		log.Printf("Failed to unmarshal '%s'\n", body)
		return Cfg{}, err
	}

	newCfg.validate()

	return newCfg, nil
}

func (cfg *Cfg) writeCfg(w io.Writer) error {
	body, err := json.MarshalIndent(*cfg, " ", " ")
	if err != nil {
		log.Println("Failed to marshal cfg:", cfg)
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func writeDefaultCfg(w io.Writer) error {
	defaultCfg := newCfgDefault()
	return defaultCfg.writeCfg(w)
}

func writeDefaultCfgTo(filepath string) error {
	fd, err := os.Create(filepath)
	if err != nil {
		log.Println("Failed to open ", filepath)
		return err
	}
	defer fd.Close()
	return writeDefaultCfg(fd)
}

func (cfg *Cfg) validate() {
	newPath := func(path string) string {
		return filepath.Join(cfg.OutputFolder, path)
	}
	checkIsPathAbs(cfg.OutputFolder)
	cfg.AtomFile = newPath(cfg.AtomFile)
	checkIsPathAbs(cfg.AtomFile)
	locationLink, err := url.JoinPath(cfg.LocationLink, cfg.SrcFolder)
	checkErr(err)
	cfg.LocationLink = locationLink
	cfg.YtdlpDownloadArchive = newPath(cfg.YtdlpDownloadArchive)
	checkIsPathAbs(cfg.YtdlpDownloadArchive)
	cfg.SrcFolder = newPath(cfg.SrcFolder)
	checkIsPathAbs(cfg.SrcFolder)
	if cfg.LimitDownload < 1 {
		log.Fatalln("Warning: LimitDowload must be at least 1")

	}
	existedAudioFormat := false
	for _, audioFormat := range YtdlpAudioFormats {
		if audioFormat == cfg.DownloadAudioFormat {
			existedAudioFormat = true
			break
		}
	}

	if !existedAudioFormat {
		log.Fatalf("Warning: DownloadAudioFormat must be choosen from %v, your provided format is '%s'\n", YtdlpAudioFormats, cfg.DownloadAudioFormat)
	}
}

func newCfgDefault() Cfg {
	cfg := Cfg{
		AtomFile:             atomFileDefault,
		AuthorLink:           authorLinkDefault,
		ChannelTitle:         channelTitleDefault,
		LocationLink:         locationLinkDefault,
		LocationType:         HttpLocation,
		Urls:                 []string{"", ""},
		VideosToDowload:      videosToDownloadDefault,
		WeeksToDelete:        weeksToDeleteDefault,
		WeeksToDownload:      weeksToDownloadDefault,
		YtdlpDownloadArchive: ytdlpDownloadArchiveDefault,
		YtdlpProgram:         ytdlpProgramDefault,
		SrcFolder:            srcFolderDefault,
		OutputFolder:         outputFolderDefault,
		LimitDownload:        limitDownloadDefault,
		DownloadAudioFormat:  downloadAudioFormatDefault,
	}

	return cfg
}
