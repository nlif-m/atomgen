package ytdlp

import (
	"log"
	"os/exec"
)

const (
	InfoJsonExtension = ".info.json"
)

var AudioFormats = [...]string{"aac", "alac", "flac", "m4a", "mp3", "opus", "vorbis", "wav"}

type Ytdlp struct {
	programName string
}

type YtdlpInfoJson struct {
	Description string `json:"description"`
}

func New(programName string) Ytdlp {
	return Ytdlp{programName: programName}
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
