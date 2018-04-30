package log

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseFileNameIntoBlock(t *testing.T) {
	Convey("ParseFileNameIntoBlock", t, func() {
		filename := "1234-5678"
		b, err := ParseFileNameIntoBlock(filename)
		So(b, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(b.StartTime, ShouldEqual, 1234)
		So(b.EndTime, ShouldEqual, 5678)
	})
}

func TestReduceToTimeRange(t *testing.T) {
	Convey("ReduceToTimeRange", t, func() {
		block := &Block{
			StartTime: 13000,
			EndTime:   17000,
			Service:   "test",
			Level:     "reader",
			Messages: []*Message{
				&Message{Text: "Foo2", Timestamp: 13000},
				&Message{Text: "Bar2", Timestamp: 13500},
				&Message{Text: "Baz2", Timestamp: 17000},
			},
		}
		block.ReduceToTimeRange(13001, 17000)
		So(len(block.Messages), ShouldEqual, 2)
		So(block.StartTime, ShouldEqual, 13500)
		So(block.EndTime, ShouldEqual, 17000)
	})
}
