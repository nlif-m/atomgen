package config

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/nlif-m/atomgen/utils"
	"github.com/nlif-m/atomgen/ytdlp"
)

const (
	// Cfg
	ytdlpProgramDefault                  string = "yt-dlp"
	outputFolderDefault                  string = "/tmp/test"
	srcFolderDefault                     string = "src"
	atomFileDefault                      string = "atom.xml"
	channelTitleDefault                  string = "test page"
	authorLinkDefault                    string = "https://example.com"
	locationLinkDefault                  string = "https://example.com" // .../src
	ytdlpDownloadArchiveDefault          string = "downloaded.txt"
	weeksToDownloadDefault               uint   = 4
	weeksToDeleteDefault                 uint   = 0
	videosToDownloadDefault              int    = 10
	generateAtomFileDefault              bool   = true
	limitDownloadDefault                 uint   = 10
	downloadAudioFormatDefault           string = "mp3"
	programRestartIntervalMinutesDefault uint   = 60
	ytdlpDurationLowerLimitDefault       uint   = 120

	// MimeDetect
	DetectContentTypeMost = 512

	// LocationType
	HttpLocation LocationType = "http"
	S3Location   LocationType = "s3" // TODO: implement s3 support
)

type LocationType string

type Cfg struct {
	AtomFile                      string
	AuthorLink                    string
	ChannelTitle                  string
	LocationLink                  string
	LocationType                  LocationType
	Urls                          []string
	VideosToDowload               int
	WeeksToDelete                 uint
	WeeksToDownload               uint
	YtdlpDownloadArchive          string
	YtdlpProgram                  string
	YtdlpDurationLowerLimit       uint
	SrcFolder                     string
	OutputFolder                  string
	LimitDownload                 uint
	DownloadAudioFormat           string
	ProgramRestartIntervalMinutes uint
	TelegramBotToken              string
	TelegramAdminId               string
}

func NewFromFile(filePath string) (cfg Cfg, err error) {
	fd, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open '%s' \n", filePath)
		return cfg, err
	}
	defer fd.Close()

	cfg, err = Read(fd)
	if err != nil {
		log.Printf("failed to read as json %s\n", filePath)
		return cfg, err
	}
	cfg.Validate()
	return cfg, nil
}

func Read(r io.Reader) (cfg Cfg, err error) {
	err = json.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil

}

func (cfg *Cfg) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent(" ", " ")
	err := encoder.Encode(cfg)
	if err != nil {
		log.Println("Failed to encode cfg:", cfg)
		return err
	}
	return nil
}

func WriteDefault(w io.Writer) error {
	defaultCfg := NewDefault()
	return defaultCfg.Write(w)
}

// TODO: Write a test that generate default config and then validate it
func WriteDefaultTo(filepath string) error {
	fd, err := os.Create(filepath)
	if err != nil {
		log.Println("Failed to open ", filepath)
		return err
	}
	defer fd.Close()
	return WriteDefault(fd)
}

// TODO: Make it return error instead of just panic
func (cfg *Cfg) Validate() {
	newPath := func(path string) string {
		return filepath.Join(cfg.OutputFolder, path)
	}
	utils.CheckIsPathAbs(cfg.OutputFolder)
	cfg.AtomFile = newPath(cfg.AtomFile)
	utils.CheckIsPathAbs(cfg.AtomFile)
	locationLink, err := url.JoinPath(cfg.LocationLink, cfg.SrcFolder)
	utils.CheckErr(err)
	cfg.LocationLink = locationLink
	cfg.YtdlpDownloadArchive = newPath(cfg.YtdlpDownloadArchive)
	utils.CheckIsPathAbs(cfg.YtdlpDownloadArchive)
	cfg.SrcFolder = newPath(cfg.SrcFolder)
	os.MkdirAll(cfg.SrcFolder, 0)
	utils.CheckIsPathAbs(cfg.SrcFolder)
	if cfg.LimitDownload < 1 {
		log.Fatalln("Warning: LimitDowload must be at least 1")

	}
	existedAudioFormat := false
	for _, audioFormat := range ytdlp.AudioFormats {
		if audioFormat == cfg.DownloadAudioFormat {
			existedAudioFormat = true
			break
		}
	}

	if !existedAudioFormat {
		log.Fatalf("Warning: DownloadAudioFormat must be choosen from %v, your provided format is '%s'\n", ytdlp.AudioFormats, cfg.DownloadAudioFormat)
	}
}

func NewDefault() Cfg {
	cfg := Cfg{
		AtomFile:                      atomFileDefault,
		AuthorLink:                    authorLinkDefault,
		ChannelTitle:                  channelTitleDefault,
		LocationLink:                  locationLinkDefault,
		LocationType:                  HttpLocation,
		Urls:                          []string{"", ""},
		VideosToDowload:               videosToDownloadDefault,
		WeeksToDelete:                 weeksToDeleteDefault,
		WeeksToDownload:               weeksToDownloadDefault,
		YtdlpDownloadArchive:          ytdlpDownloadArchiveDefault,
		YtdlpProgram:                  ytdlpProgramDefault,
		YtdlpDurationLowerLimit:       ytdlpDurationLowerLimitDefault,
		SrcFolder:                     srcFolderDefault,
		OutputFolder:                  outputFolderDefault,
		LimitDownload:                 limitDownloadDefault,
		DownloadAudioFormat:           downloadAudioFormatDefault,
		ProgramRestartIntervalMinutes: programRestartIntervalMinutesDefault,
	}

	return cfg
}
