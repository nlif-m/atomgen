package main

// TODO: make a daemon that run program eveny N time
// TODO: Make a separate commands, for example only 'atomgen download' to download and nothing else and etc.

import (
	"flag"
	"fmt"
	"os"
)

const (
	detectContentTypeMost = 512

	ytdlp                = "yt-dlp"
	src_folder           = "src"
	urlsFile             = "urls.csv"
	ytdlpDownloadArchive = src_folder + string(os.PathSeparator) + "downloaded.txt"
	ytdlpOutputTemplate  = src_folder + string(os.PathSeparator) + "%(uploader)s %(title)s.%(ext)s"

	atomFile = "rss.xml"

	channelTitle = "test page1"
	channelLink  = "https://rss.yasal.xyz"
)

func main() {
	downloadVideosFromFlagStr := flag.String("dwl", "yes", fmt.Sprint("Download videos from ", urlsFile, " ?"))
	generateAtomRssFileFlagStr := flag.String("atom", "yes", "Generate atom file ?")

	flag.Parse()

	switch *downloadVideosFromFlagStr {
	case "yes", "y":
		downloadVideosFromFile(urlsFile)
	}

	switch *generateAtomRssFileFlagStr {
	case "yes", "y":
		generateAtomRssFile(atomFile, src_folder)
	}
}
