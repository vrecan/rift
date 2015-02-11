package pullrift

import (
	"errors"
	log "github.com/cihub/seelog"
	zmq "github.com/pebbe/zmq4"
	"github.com/vrecan/rift/c"
	"math/rand"
	"time"
)

var InvalidSampleRateError = errors.New("Invalid sample rate, it must be between 1 and 100")
var NoValidRifts = errors.New("No valid rifts created")

type stats struct {
	riftQueue string
	puller    string
	success   int64
	failed    int64
}

type PullRift struct {
	URL        string
	pull       *zmq.Socket
	pushers    []PushSocket
	logChannel chan *stats
	poller     *zmq.Poller
}

type PushSocket struct {
	push       *zmq.Socket
	URL        string
	SampleRate int
	success    int64
	failed     int64
}

func NewPushSocket(url string, sampleRate int) (push PushSocket, err error) {
	push = PushSocket{}
	if sampleRate > 100 || sampleRate <= 0 {
		return push, InvalidSampleRateError
	}
	push.push, err = zmq.NewSocket(zmq.PUSH)
	if nil != err {
		return push, err
	}
	push.URL = url
	push.SampleRate = sampleRate
	err = push.push.Connect(url)
	return push, err
}

//Create new pull rift that sends to any number of push sockets
func NewPullRift(rcvURL string, confs []c.PullConf) (rift PullRift, err error) {
	rift = PullRift{}
	rand.Seed(100)
	rift.poller = zmq.NewPoller()
	rift.logChannel = make(chan *stats, 100)
	rift.pull, err = zmq.NewSocket(zmq.PULL)
	if nil != err {
		return rift, err
	}
	err = rift.pull.Bind(rcvURL)
	rift.URL = rcvURL
	if nil != err {
		return rift, err
	}
	rift.poller.Add(rift.pull, zmq.POLLIN)
	for _, conf := range confs {
		push, err_ := NewPushSocket(conf.URL, conf.SampleRate)
		if nil != err_ {
			log.Error("Failed to rift to url: ", conf.URL, " error: ", err_)
			continue
		}
		rift.pushers = append(rift.pushers, push)
	}

	if rift.RiftSize() <= 0 {
		err = NoValidRifts
	}

	return rift, err
}

//Main loop for pullRift.
func (r *PullRift) Run() {
	log.Info("Starting Rift :", r)
	go r.logRift()
	var bytes []byte

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, push := range r.pushers {
				r.logChannel <- &stats{riftQueue: r.URL, puller: push.URL, success: push.success, failed: push.failed}
				push.failed = 0
				push.success = 0
			}

		default:
			res, err := r.poller.Poll(5 * time.Second)
			if nil != err {
				log.Error("Failed on poll: ", err)
			}
			if len(res) > 0 {
				if res[0].Events&zmq.POLLIN != 0 {
					bytes, err = r.pull.RecvBytes(0)
					if nil == err {
						random := rand.Int()%100 + 1 // random number between 0 - 100
						for _, push := range r.pushers {
							if random <= push.SampleRate {
								_, err = push.push.SendBytes(bytes, zmq.DONTWAIT)
								if nil != err {
									push.failed++
								} else {
									push.success++
								}
							}
						}
					} else {
						log.Error("Failed to recv from rift: ", err)
					}
				}
			}

		}
	}

}

func (r *PullRift) logRift() {
	for stats := range r.logChannel {
		log.Info("RiftQueue: ", stats.riftQueue, " PullQueue: ", stats.puller, " Sucess: ", stats.success, " Failed: ", stats.failed)
		// log.Info(stats)
	}
}

//Returns the number of queues we are sending to
func (r *PullRift) RiftSize() (size int) {
	return len(r.pushers)
}

//Close our sockets.
func (r *PullRift) Close() {
	if nil != r.pull {
		r.pull.Close()
	}
	for _, push := range r.pushers {
		push.push.Close()
	}
}
