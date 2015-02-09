package c

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPullRift(t *testing.T) {
	Convey("Read file that doesn't exist", t, func() {
		_, err := GetConf("/path/that/should/never/exist/woo.yaml")
		ShouldNotBeNil(err)
	})
}
