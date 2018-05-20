package log

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileReader(t *testing.T) {
	pathPrefix = "test"
	Convey("FileReader", t, func() {
		b1 := &Block{
			StartTime: 5002,
			EndTime:   10001,
			Service:   "test",
			Level:     "file_reader",
			Messages: []*Message{
				&Message{Text: "Foo", Timestamp: 5002},
				&Message{Text: "Bar", Timestamp: 7005},
				&Message{Text: "Baz", Timestamp: 10001},
			},
		}
		b2 := &Block{
			StartTime: 13000,
			EndTime:   17000,
			Service:   "test",
			Level:     "file_reader2",
			Messages: []*Message{
				&Message{Text: "Foo2", Timestamp: 13000},
				&Message{Text: "Bar2", Timestamp: 13500},
				&Message{Text: "Baz2", Timestamp: 17000},
			},
		}
		b3 := &Block{
			StartTime: 30000,
			EndTime:   40003,
			Service:   "test2",
			Level:     "file_reader3",
			Messages: []*Message{
				&Message{Text: "Foo3", Timestamp: 30000},
				&Message{Text: "Bar3", Timestamp: 34000},
				&Message{Text: "Baz3", Timestamp: 40003},
			},
		}
		b1.WriteToFile()
		b2.WriteToFile()
		b3.WriteToFile()

		r := FileReader{}
		Convey("get services", func() {
			services := r.GetServices()
			So(len(services), ShouldEqual, 2)
			So(services[0], ShouldEqual, "test")
			So(services[1], ShouldEqual, "test2")
		})
		Convey("get levels", func() {
			levels := r.GetLevels("test")
			So(len(levels), ShouldEqual, 2)
			So(levels[0], ShouldEqual, "file_reader")
			So(levels[1], ShouldEqual, "file_reader2")
		})
	})
	os.RemoveAll(pathPrefix)
}
