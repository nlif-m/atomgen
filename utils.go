package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func deleteOldFilesFromFolder(srcFolder string, howManyWeeksIsOld int) {

	log.Println("Start deleting old videos")
	files, err := os.ReadDir(srcFolder)

	if err != nil {
		log.Printf("WARNING: failed to read %s folder to delete old videos\n", srcFolder)
		return
	}

filesLoop:
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		if !fileInfo.ModTime().Before(time.Now().AddDate(0, 0, -howManyWeeksIsOld*7)) {
			continue filesLoop
		}

		filePath := srcFolder + string(os.PathSeparator) + file.Name()
		log.Println(fmt.Sprint("Deleting file older than ", howManyWeeksIsOld, " weeks "), file.Name())
		err = os.Remove(filePath)
		if err != nil {
			log.Println("Warning: failed to delete file at ", filePath, err)
		}

		continue filesLoop

	}

	log.Println("Finish deleting old videos")

}
