package pullrift

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vrecan/rift/c"
	"testing"
)

func TestPullRift(t *testing.T) {

	Convey("Test adding empty zmq pullers", t, func() {
		pullConf := []c.PullConf{}
		rift, err := NewPullRift("tcp://127.0.0.1:9900", pullConf)
		defer rift.Close()
		So(err, ShouldEqual, NoValidRifts)

	})

	Convey("Test adding multiple valid pullers", t, func() {
		pullConf := []c.PullConf{}
		pullConf = append(pullConf, c.PullConf{URL: "tcp://127.0.0.1:9880", SampleRate: 100})
		pullConf = append(pullConf, c.PullConf{URL: "tcp://127.0.0.1:9881", SampleRate: 100})
		rift, err := NewPullRift("tcp://127.0.0.1:9901", pullConf)
		defer rift.Close()
		So(err, ShouldBeNil)
		So(rift.RiftSize(), ShouldEqual, 2)
	})

	Convey("Test adding one invalid one valid, should log error for invalid and continue.", t, func() {
		pullConf := []c.PullConf{}
		pullConf = append(pullConf, c.PullConf{URL: "invalid", SampleRate: 100})
		pullConf = append(pullConf, c.PullConf{URL: "tcp://127.0.0.1:9881", SampleRate: 100})
		rift, err := NewPullRift("tcp://127.0.0.1:9902", pullConf)
		defer rift.Close()
		So(err, ShouldBeNil)
		So(rift.RiftSize(), ShouldEqual, 1)
	})

}
