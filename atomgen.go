package main

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
	// fd, err := os.OpenFile("log.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	// if err != nil {
	// 	panic(err)
	// }
	// defer fd.Close()
	// log.SetOutput(fd)

	downloadVideosFromFile(URLS_CSV_FILE)
	generateAtomRssFile(ATOM_FILE, SRC_FOLDER)
}
