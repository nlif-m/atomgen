package main

import (
	"log"
	"path"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkIsPathAbs(filepath string) {
	if !path.IsAbs(filepath) {
		log.Fatalf("'%s' is %s\n", filepath, path.ErrBadPattern)
	}
}
