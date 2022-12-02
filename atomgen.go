package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	detectContentTypeMost = 512

	defautlYtdlp = "yt-dlp"
	src_folder   = "src"
	urlsFile     = "urls.csv"

	ytdlpDownloadArchive     = "downloaded.txt"
	ytdlpOutputTemplate      = src_folder + string(os.PathSeparator) + "%(uploader)s %(title)s.%(ext)s"
	howManyWeeksIsOld    int = 4 // media older than that amount will be ignored

	atomFile = "atom.xml"

	channelTitle = "test page1"
	channelLink  = "https://rss.yasal.xyz"
)

func main() {
	// TODO: Add a flag to change download limit
	downloadVideosFromFileFlag := flag.Bool("dwl", true, fmt.Sprint("Download videos from ", urlsFile, " ?"))
	generateAtomRssFileFlag := flag.Bool("atom", true, "Generate atom file ?")
	deleteOldVideosFromFolderFlag := flag.Bool("delete-old", false, fmt.Sprint("Delete files older than ", howManyWeeksIsOld, " weeks ?"))
	flag.Parse()

	yt := newYtdlp(defautlYtdlp)
	if *downloadVideosFromFileFlag {
		yt.DownloadVideosFromFile(urlsFile)
	}

	if *deleteOldVideosFromFolderFlag {
		deleteOldFilesFromFolder(src_folder, howManyWeeksIsOld)
	}

	if *generateAtomRssFileFlag {
		generateAtomRssFile(atomFile, src_folder)
	}

}
