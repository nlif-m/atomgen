package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/nlif-m/atomgen/config"
	"github.com/nlif-m/atomgen/utils"
	"github.com/nlif-m/atomgen/ytdlp"
)

var (
	programConfig string
	genConfig     string
)

func main() {
	// TODO: Add a ability to make this configs for each url individually
	flag.StringVar(&genConfig, "genConfig", "", "generate default config file")
	flag.StringVar(&programConfig, "config", "", "config file")
	flag.Parse()

	if genConfig != "" {
		err := config.WriteDefaultTo(genConfig)
		utils.CheckErr(err)
		return

	}

	if programConfig == "" {
		flag.Usage()
		return
	}

	cfg, err := config.NewFromFile(programConfig)
	utils.CheckErr(err)
	yt := ytdlp.New(cfg.YtdlpProgram)

	atomgen := newAtomgen(yt, cfg)

	startWorkChan := make(chan bool)
	var wg sync.WaitGroup
	go func(ch chan bool) {
		for {
			<-ch
			if atomgen.cfg.VideosToDowload != 0 {
				err := atomgen.DownloadVideos()
				utils.CheckErr(err)
			}

			if atomgen.cfg.WeeksToDelete != 0 {
				err := atomgen.deleteOldFiles()
				utils.CheckErr(err)
			}

			err = atomgen.generateAtomFeed()
			utils.CheckErr(err)
			wg.Done()
		}
	}(startWorkChan)
	timeToSleep := time.Duration(cfg.ProgramRestartIntervalMinutes * uint(time.Minute))
	tick := time.Tick(timeToSleep)
	for {
		<-tick
		startWorkChan <- true
		wg.Add(1)
		wg.Wait()
		log.Printf("Start sleeping regular download for  %f minutes\n", timeToSleep.Minutes())

	}
}
