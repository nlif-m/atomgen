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
	ytdlpProgramDefault         string        = "yt-dlp"
	outputFolderDefault         string        = "/tmp/test"
	srcFolderDefault            string        = "src"
	urlsFileDefault             string        = "urls.csv"
	atomFileDefault             string        = "atom.xml"
	channelTitleDefault         string        = "test page"
	authorLinkDefault           string        = "https://example.com"
	locationLinkDefault         string        = "https://example.com" // .../src
	ytdlpDownloadArchiveDefault string        = "downloaded.txt"
	weeksToDownloadDefault      uint          = 4
	weeksToDeleteDefault        uint          = 0
	videosToDownloadDefault     int           = 10
	generateAtomFileDefault     bool          = true
	limitDownloadDefault        uint          = 10
	downloadAudioFormatDefault  string        = "mp3"
	programRestartIntervalMinutesDefault      uint = 60
	
	// MimeDetect
	DetectContentTypeMost = 512

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
	ProgramRestartIntervalMinutes uint
	
}

func NewFromFile(filePath string) (Cfg, error) {
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

	newCfg.Validate()

	return newCfg, nil
}

func (cfg *Cfg) Write(w io.Writer) error {
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

func WriteDefault(w io.Writer) error {
	defaultCfg := NewDefault()
	return defaultCfg.Write(w)
}

func WriteDefaultTo(filepath string) error {
	fd, err := os.Create(filepath)
	if err != nil {
		log.Println("Failed to open ", filepath)
		return err
	}
	defer fd.Close()
	return WriteDefault(fd)
}

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
		ProgramRestartIntervalMinutes: programRestartIntervalMinutesDefault,
	}

	return cfg
}
