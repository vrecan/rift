package main

import (
	log "github.com/cihub/seelog"
	"github.com/vrecan/rift/start"
        "runtime"
)

func main() {

	defer log.Flush()
	var err error
        var logger log.LoggerInterface
	if(runtime.GOOS == "windows") {
	logger, err = log.LoggerFromConfigAsFile("seelog_windows.xml")
	} else {
	logger, err = log.LoggerFromConfigAsFile("seelog.xml")
 	}

	if err != nil {
		log.Warn("Failed to load config", err)
	} else {
		log.ReplaceLogger(logger)
	}

	start.Run()
}
