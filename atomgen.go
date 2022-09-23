package main

// TODO: make a daemon that run program eveny N time
// TODO: Make a separate commands, for example only 'atomgen download' to download and nothing else and etc.

import (
	"os"
	// "flag"
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
	downloadVideosFromFile(URLS_CSV_FILE)
	generateAtomRssFile(ATOM_FILE, src_folder)
}
