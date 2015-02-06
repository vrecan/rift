package start

import (
	log "github.com/cihub/seelog"
	rift "github.com/vrecan/rift/pullrift"
	"os"
	"os/signal"
	"sync"
)

func Run() {
	log.Info("Started")

	shutdown := make(chan int, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//start signal handling routine
	go Death(signals, shutdown, wg)

	pushURLs := []string{"tcp://10.1.10.42:13109", "tcp://10.1.10.172:13100"}
	pullRift, err := rift.NewPullRift("tcp://10.1.10.42:13100", pushURLs)
	if nil != err {
		log.Critical("Failed to create rift: ", err)
	}
	go pullRift.Run()

	//hanndle shutdown
	wg.Wait()
	log.Info("Shutting down go routines")
	close(shutdown)
	close(signals)
	pullRift.Close()
	log.Info("Shutdown complete")
}

//Manage death of application by signal
func Death(c <-chan os.Signal, death chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for sig := range c {
			switch sig.String() {
			case "terminated":
				{
					death <- 1
					return
				}
			case "interrupt":
				{
					death <- 2
					return
				}
			default:
				{
					log.Debug("Signal: ", sig.String())
					death <- 3
					return
				}
			}

		}
	}
}
