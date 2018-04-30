package log

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsInTimeRange(t *testing.T) {
	Convey("IsInTimeRange", t, func() {
		m1 := Message{Timestamp: 5000}
		m2 := Message{Timestamp: 7000}

		So(m1.IsInTimeRange(4000, 6000), ShouldBeTrue)
		So(m2.IsInTimeRange(4000, 6000), ShouldBeFalse)
	})
}
