package main

// TODO: make a daemon that run program eveny N time
// TODO: Make a separate commands, for example only 'atomgen download' to download and nothing else and etc.

import (
	"os"
	// "flag"
)

const (
	DetectContentTypeMost = 512

	YTDLP                  = "yt-dlp"
	SRC_FOLDER             = "src"
	URLS_CSV_FILE          = "urls.csv"
	YTDLP_DOWNLOAD_ARCHIVE = SRC_FOLDER + string(os.PathSeparator) + "downloaded.txt"
	YTDLP_OUTPUT_TEMPLATE  = SRC_FOLDER + string(os.PathSeparator) + "%(uploader)s %(title)s.%(ext)s"

	ATOM_FILE = "rss.xml"

	CHANNEL_TITLE = "test page"
	CHANNEL_LINK  = "https://rss.yasal.xyz"
)

func main() {
	downloadVideosFromFile(URLS_CSV_FILE)
	generateAtomRssFile(ATOM_FILE, SRC_FOLDER)
}
