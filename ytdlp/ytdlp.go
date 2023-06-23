package ytdlp

import (
	"log"
	"os/exec"
	"regexp"
)

const (
	InfoJsonExtension = ".info.json"
)

type YtdlpURLType uint

const (
	YoutubeVideoType YtdlpURLType = iota
	YoutubePlaylistType
	VkVideoType
	UndefinedType
)

func (y YtdlpURLType) String() string {
	switch y {
	case YoutubeVideoType:
		return "youtube_video"
	case YoutubePlaylistType:
		return "youtube_playlist"
	case VkVideoType:
		return "vk_video"
	case UndefinedType:
		return "undefined"
	}
	log.Printf("YtdlpURLType %d failed to convert to string\n", y)
	return "unknown"
}

var AudioFormats = [...]string{"aac", "alac", "flac", "m4a", "mp3", "popus", "vorbis", "wav"}

type Ytdlp struct {
	programName string
}

type YtdlpInfoJson struct {
	Description string `json:"description"`
}

func New(programName string) Ytdlp {
	return Ytdlp{programName: programName}
}
func NewDefault() Ytdlp {
	return Ytdlp{programName: "yt-dlp"}
}

func (yt *Ytdlp) NewCmdWithArgs(Args ...string) *exec.Cmd {
	return exec.Command(yt.programName, Args...)

}

func (yt *Ytdlp) GetChannelNameFromURL(URL string) (channelName string, err error) {
	cmd := yt.NewCmdWithArgs(
		"-O", "%(channel)s",
		"--playlist-end", "1",
		URL)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: failed to get channel name for '%s'\n cmd: [%s]\t%s\n%s\n", URL, cmd, err, string(cmdOutput))
		return "", err
	}
	return string(cmdOutput), nil
}

func (yt *Ytdlp) GetVersion() (version string, err error) {
	cmd := yt.NewCmdWithArgs("--version")
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: failed to get version of '%s'\n cmd: [%s]\t%s\n", yt.programName, cmd, err)
		return "", err
	}
	return string(cmdOutput), nil
}

var youtubeVideoRegexp = regexp.MustCompile(`(https:\/\/|)(www\.|)youtube\.com\/watch\?v=.+`)
var youtubeVideoRegexp2 = regexp.MustCompile(`(https:\/\/|)youtu\.be\/.+`)
var youtubePlaylistRegexp = regexp.MustCompile(`(https:\/\/|)(www\.|)youtube\.com\/playlist\?list=.+`)
var vkVideoRegexp = regexp.MustCompile(`(https:\/\/|)vk\.com\/video(s|).+`)

// TODO: Write a tests
func (yt *Ytdlp) IsDownloadable(rawURL string) (ytType YtdlpURLType, URL string, downloadable bool) {
	URL = youtubeVideoRegexp.FindString(rawURL)
	if URL != "" {
		return YoutubeVideoType, URL, true
	}
	URL = youtubeVideoRegexp2.FindString(rawURL)
	if URL != "" {
		return YoutubeVideoType, URL, true
	}

	URL = youtubePlaylistRegexp.FindString(rawURL)
	if URL != "" {
		return YoutubePlaylistType, URL, true
	}

	URL = vkVideoRegexp.FindString(rawURL)
	if URL != "" {
		return VkVideoType, URL, true
	}

	return UndefinedType, "", false
}
