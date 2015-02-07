package start

import (
	log "github.com/cihub/seelog"
	"github.com/vrecan/rift/c"
	rift "github.com/vrecan/rift/pullrift"
	"os"
	"os/signal"
	"sync"
)

func Run() {
	log.Info("Started")
	conf, err := c.GetConf("c/c.yml")
	if nil != err {
		log.Error(err)
	}
	log.Info("conf: ", conf)

	shutdown := make(chan int, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//start signal handling routine
	go Death(signals, shutdown, wg)
	rifts := make([]rift.PullRift, len(conf.Rifts))
	for _, r := range conf.Rifts {
		nRift, err := rift.NewPullRift(r.Pull, r.Push)
		if nil == err {
			go nRift.Run()
			rifts = append(rifts, nRift)
		} else {
			log.Error("Failed to construct rift... : ", err)
		}
	}

	//hanndle shutdown
	wg.Wait()
	log.Info("Shutting down go routines")
	close(shutdown)
	close(signals)
	for _, r := range rifts {
		r.Close()
	}
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
