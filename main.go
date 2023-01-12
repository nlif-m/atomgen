package main

import (
	"flag"
)

var (
	config    string
	genConfig string
)

func main() {
	// TODO: Add a ability to make this configs for each url individually
	flag.StringVar(&genConfig, "genConfig", "", "generate default config file")
	flag.StringVar(&config, "config", "", "config file")
	flag.Parse()

	if genConfig != "" {
		err := writeDefaultCfgTo(genConfig)
		checkErr(err)
		return

	}

	if config == "" {
		flag.Usage()
		return
	}

	cfg, err := newCfgFromFile(config)
	checkErr(err)
	yt := newYtdlp(cfg.YtdlpProgram)

	atomgen := newAtomgen(yt, cfg)
	if atomgen.cfg.VideosToDowload != 0 {
		err := atomgen.DownloadVideos()
		checkErr(err)
	}

	if atomgen.cfg.WeeksToDelete != 0 {
		err := atomgen.deleteOldFiles()
		checkErr(err)
	}

	err = atomgen.generateAtomFeed()
	checkErr(err)
}
