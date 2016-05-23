package sensupluginsslack

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHandlerSlack(t *testing.T) {

	Convey("When creating a new slack message", t, func() {
		var out string

		Convey("The command should return some output", func() {
			So(out, ShouldBeBlank)
		})
	})
}
