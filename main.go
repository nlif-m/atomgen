package main

import (
	"flag"
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

	fullUpdateChan := make(chan bool)
	atomFileUpdateChan := make(chan bool)
	atomgen.generateAtomFeed()
	atomgen.fullUpdate()
	var wg sync.WaitGroup
	go func(fullUpdateChan chan bool, atomFileUpdateChan chan bool) {
		go func() {
			TgBot(atomgen, atomFileUpdateChan)
		}()
		for {
			select {
			case <-fullUpdateChan:
				wg.Add(1)
				atomgen.fullUpdate()
				wg.Done()

			case <-atomFileUpdateChan:
				wg.Add(1)
				err = atomgen.generateAtomFeed()
				utils.CheckErr(err)
				wg.Done()
			}
		}
	}(fullUpdateChan, atomFileUpdateChan)
	timeToSleep := time.Duration(cfg.ProgramRestartIntervalMinutes * uint(time.Minute))
	tick := time.Tick(timeToSleep)
	for {
		<-tick
		fullUpdateChan <- true
	}
}
