package main

import (
	"flag"
	"fmt"
)

const (
	detectContentTypeMost = 512
)

var (
	ytdlpProgram                string
	ytdlpProgramDefault         string = "yt-dlp"
	srcFolder                   string
	srcFolderDefault            string = "src"
	urlsFile                    string
	urlsFileDefault             string = "urls.csv"
	atomFile                    string
	atomFileDefault             string = "atom.xml"
	channelTitle                string
	channelTitleDefault         string = "test page"
	channelLink                 string
	channelLinkDefault          string = "https://rss.yasal.xyz"
	ytdlpDownloadArchive        string
	ytdlpDownloadArchiveDefault string = "downloaded.txt"
	weeksToDownload             uint
	weeksToDownloadDefault      uint = 4
	weeksToDelete               uint
	weeksToDeleteDefault        uint = 0
	videosToDownload            int
	videosToDownloadDefault     int = 10
	generateAtomFile            bool
	generateAtomFileDefault     bool = true
)

func main() {
	// TODO: Add a ability to make this configs for each url individually

	flag.StringVar(&ytdlpProgram, "ytdlpProgram", ytdlpProgramDefault, "What the name of yt-dlp binary?")
	flag.StringVar(&srcFolder, "srcFolder", srcFolderDefault, "What folder to use for downloads?")
	flag.StringVar(&urlsFile, "urlsFile", urlsFileDefault, "What file that contain urls to download ?")
	flag.StringVar(&atomFile, "atomFile", atomFileDefault, "What file to write atom feed?")
	flag.StringVar(&channelTitle, "channelTitle", channelTitleDefault, "What title of atom feed?")
	flag.StringVar(&channelLink, "channelLink", channelLinkDefault, "What url of atom feed? examaple: https://example.com/atom.xml , there https://example.com is channelLink")
	flag.StringVar(&ytdlpDownloadArchive, "ytdlpDownloadArchive", ytdlpDownloadArchiveDefault, "")
	flag.UintVar(&weeksToDownload,
		"weeksToDownload",
		weeksToDownloadDefault,
		fmt.Sprintf("How many weeks of video to download? if equal 0, try to download all videos on url only limited by amount of videosToDownload(default %v)\t", videosToDownloadDefault))

	flag.UintVar(&weeksToDelete,
		"weeksToDelete",
		weeksToDeleteDefault,
		"How many weeks must pass to delete video? if equal 0, don't delete")

	flag.IntVar(&videosToDownload,
		"videosToDownload",
		videosToDownloadDefault,
		fmt.Sprintf("How many videos to maximux download for from specific amount of weeks(default %v weeks)? if equal 0, don't download. if equal -1 don't limit amount of downloading.", weeksToDownloadDefault))

	flag.BoolVar(&generateAtomFile, "generateAtom", generateAtomFileDefault, "Generate atom file ?")
	flag.Parse()

	yt := newYtdlp(ytdlpProgram)
	if videosToDownload != 0 {
		yt.DownloadVideosFromFile(urlsFile)
	}

	if weeksToDelete != 0 {
		deleteOldFilesFromFolder(srcFolder, weeksToDelete)
	}

	if generateAtomFile {
		generateAtomRssFile(atomFile, srcFolder)
	}

}
