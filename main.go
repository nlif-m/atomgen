package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"time"

	"log"

	"github.com/nlif-m/atomgen/config"
	"github.com/nlif-m/atomgen/utils"
	"github.com/nlif-m/atomgen/ytdlp"
)

var (
	programConfig string
	genConfig     string
	portConfig    string
)

func main() {
	// TODO: Add a ability to make this configs for each url individually
	flag.StringVar(&genConfig, "genConfig", "", "generate default config file")
	flag.StringVar(&programConfig, "config", "", "config file")
	flag.StringVar(&portConfig, "address", ":3000", "address to listen")
	flag.Parse()

	if genConfig != "" {
		err := config.WriteDefaultTo(genConfig)
		if err != nil {
			log.Fatalf("ERROR: failed to generate config to: %q:%q\n", genConfig, err)
		}
		os.Exit(0)
	}

	if programConfig == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.NewFromFile(programConfig)
	if err != nil {
		log.Fatalf("ERROR: failed to read config from %q:%q\n", programConfig, err)
	}
	yt := ytdlp.New(cfg.YtdlpProgram)
	atomgen := newAtomgen(yt, cfg)

	fullUpdateChan := make(chan bool)
	atomFileUpdateChan := make(chan bool)

	// Run initial update
	go func(fullUpdateChan chan bool, atomFileUpdateChan chan bool) {
		atomFileUpdateChan <- true
		fullUpdateChan <- true
	}(fullUpdateChan, atomFileUpdateChan)

	// Run Telegram bot
	go func() {
		TgBot(atomgen, atomFileUpdateChan)
	}()

	var wg sync.WaitGroup // Why i use it since it non needed

	go func(fullUpdateChan chan bool, atomFileUpdateChan chan bool) {
		for {
			select {
			case <-fullUpdateChan:
				wg.Add(1)
				log.Println("Receive fullUpdateChan request")
				err = atomgen.fullUpdate()
				utils.CheckErr(err)
				err = atomgen.generateAtomFeed()
				utils.CheckErr(err)
				wg.Done()

			case <-atomFileUpdateChan:
				wg.Add(1)
				log.Println("Receive atomFileUpdateChan request")
				err = atomgen.generateAtomFeed()
				utils.CheckErr(err)
				wg.Done()
			}
		}
	}(fullUpdateChan, atomFileUpdateChan)

	// Run timer to call update every cfg.ProgramRestartIntervalMinutes minutes
	go func() {
		for {
			time.Sleep(time.Duration(cfg.ProgramRestartIntervalMinutes) * time.Minute)
			log.Println("Time to update, timer tick")
			fullUpdateChan <- true
		}

	}()

	http.Handle("/", http.FileServer(http.Dir(cfg.OutputFolder)))
	log.Fatalln(http.ListenAndServe(portConfig, nil))
}
