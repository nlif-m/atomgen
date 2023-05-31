package main

import (
	"flag"
	"github.com/nlif-m/atomgen/utils"
	"github.com/nlif-m/atomgen/ytdlp"
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
		utils.CheckErr(err)
		return

	}

	if config == "" {
		flag.Usage()
		return
	}

	cfg, err := newCfgFromFile(config)
	utils.CheckErr(err)
	yt := ytdlp.New(cfg.YtdlpProgram)

	atomgen := newAtomgen(yt, cfg)
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
}
