package main

import (
	log "github.com/cihub/seelog"
	"github.com/vrecan/rift/start"
)

func main() {

	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		log.Warn("Failed to load config", err)
	} else {
		log.ReplaceLogger(logger)
	}

	start.Run()
}
