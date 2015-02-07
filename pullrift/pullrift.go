package pullrift

import (
	"errors"
	log "github.com/cihub/seelog"
	zmq "github.com/pebbe/zmq4"
	"github.com/vrecan/rift/c"
	"math/rand"
)

var InvalidSampleRateError = errors.New("Invalid sample rate, it must be between 1 and 100")

type PullRift struct {
	pull  *zmq.Socket
	rifts []PushSocket
}

type PushSocket struct {
	push       *zmq.Socket
	URL        string
	SampleRate int
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
	rift.pull, err = zmq.NewSocket(zmq.PULL)
	if nil != err {
		return rift, err
	}
	err = rift.pull.Bind(rcvURL)
	if nil != err {
		return rift, err
	}
	for _, conf := range confs {
		push, err_ := NewPushSocket(conf.URL, conf.SampleRate)
		if nil != err_ {
			log.Error("Failed to rift to url: ", conf.URL, " error: ", err_)
			continue
		}
		rift.rifts = append(rift.rifts, push)
	}
	return rift, err
}

//Main loop for pullRift.
func (r *PullRift) Run() {
	log.Info("Starting Rift :", r)
	var bytes []byte
	var err error
	for {
		bytes, err = r.pull.RecvBytes(0)
		if nil == err {
			random := rand.Int()%100 + 1 // random number between 0 - 100
			for _, push := range r.rifts {
				if random <= push.SampleRate {
					push.push.SendBytes(bytes, zmq.DONTWAIT)
				}
			}
		} else {
			log.Error("Failled to pull from socket: ", err)
		}
	}
}

//Close our sockets.
func (r *PullRift) Close() {
	if nil != r.pull {
		r.pull.Close()
	}
	for _, push := range r.rifts {
		push.push.Close()
	}
}
