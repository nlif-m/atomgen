package main

import (
	"log"
	"os/exec"
)

type Ytdlp struct {
	programName string
}

func newYtdlp(programName string) Ytdlp {
	return Ytdlp{programName: programName}
}

func (yt *Ytdlp) newCmdWithArgs(Args ...string) *exec.Cmd {
	return exec.Command(yt.programName, Args...)

}

func (yt *Ytdlp) GetChannelNameFromURL(URL string) (channelName string, err error) {
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

func (yt *Ytdlp) GetVersion() (version string, err error) {
	cmd := yt.newCmdWithArgs("--version")
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: failed to get version of '%s'\n cmd: [%s]\t%s\n", yt.programName, cmd, err)
		return "", err
	}
	return string(cmdOutput), nil
}
