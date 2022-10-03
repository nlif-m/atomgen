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

	ytdlp      = "yt-dlp"
	src_folder = "src"
	urlsFile   = "urls.csv"

	ytdlpDownloadArchive     = "downloaded.txt"
	ytdlpOutputTemplate      = src_folder + string(os.PathSeparator) + "%(uploader)s %(title)s.%(ext)s"
	howManyWeeksDownload int = 4

	atomFile = "rss.xml"

	channelTitle = "test page1"
	channelLink  = "https://rss.yasal.xyz"
)

func main() {
	// TODO: Add a flat to change download limit
	downloadVideosFromFlagStr := flag.String("dwl", "yes", fmt.Sprint("Download videos from ", urlsFile, " ?"))
	generateAtomRssFileFlagStr := flag.String("atom", "yes", "Generate atom file ?")
	deleteOldVideosFromDirStr := flag.String("delete-old", "no", fmt.Sprint("Delete files older than ", howManyWeeksDownload, " weeks ?"))

	flag.Parse()

	switch *downloadVideosFromFlagStr {
	case "yes", "y":
		downloadVideosFromFile(urlsFile)
	}

	switch *deleteOldVideosFromDirStr {
	case "yes", "y":
		deleteOldVideosFromDir(src_folder, howManyWeeksDownload)
	}

	switch *generateAtomRssFileFlagStr {
	case "yes", "y":
		generateAtomRssFile(atomFile, src_folder)
	}

}
