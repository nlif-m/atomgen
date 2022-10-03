package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func deleteOldVideosFromDir(srcFolder string, howManyWeeksDownload int) {

	log.Println("Start deleting old videos")
	files, err := ioutil.ReadDir(srcFolder)

	if err != nil {
		log.Println("WARNING: failed to delete old videos")
		return
	}

filesLoop:
	for _, file := range files {
		if !file.ModTime().Before(time.Now().AddDate(0, 0, -howManyWeeksDownload*7)) {
			continue filesLoop
		}

		filePath := srcFolder + string(os.PathSeparator) + file.Name()
		log.Println(fmt.Sprint("Deleting file older than ", howManyWeeksDownload, " weeks "), file.Name())
		err = os.Remove(filePath)
		if err != nil {
			log.Println("Warning: failed to delete file at ", filePath, err)
		}

		continue filesLoop

	}

}
