package pullrift

import (
	log "github.com/cihub/seelog"
	zmq "github.com/pebbe/zmq4"
)

type pullRift struct {
	pull  *zmq.Socket
	rifts []pushSocket
}

type pushSocket struct {
	push       *zmq.Socket
	URL        string
	SampleRate int32
}

func NewPushSocket(url string, sampleRate int32) (push pushSocket, err error) {
	push = pushSocket{}
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
func NewPullRift(rcvURL string, pushURLs []string) (rift pullRift, err error) {
	rift = pullRift{}
	rift.pull, err = zmq.NewSocket(zmq.PULL)
	if nil != err {
		return rift, err
	}
	err = rift.pull.Bind(rcvURL)
	if nil != err {
		return rift, err
	}
	for _, url := range pushURLs {
		push, err_ := NewPushSocket(url, 100)
		if nil != err_ {
			log.Error("Failed to rift to url: ", url, " error: ", err_)
			continue
		}
		rift.rifts = append(rift.rifts, push)
	}
	return rift, err
}

func (r *pullRift) Run() {
	log.Info("Starting Rift :", r)
	for {
		bytes, err := r.pull.RecvBytes(0)
		if nil == err {
			for _, push := range r.rifts {
				push.push.SendBytes(bytes, zmq.DONTWAIT)
			}
		} else {
			log.Error("Failled to pull from socket: ", err)
		}
	}
}

func (r *pullRift) Close() {
	if nil != r.pull {
		r.pull.Close()
	}
	for _, push := range r.rifts {
		push.push.Close()
	}
}
