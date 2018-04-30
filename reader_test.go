package log

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type CacheAccessorMock struct{}

func (m *CacheAccessorMock) GetCache() *Cache {
	return NewCache()
}

func TestGetBlocksInTimeRange(t *testing.T) {
	pathPrefix = "test"
	Convey("GetBlocksInTimeRange", t, func() {
		b1 := &Block{
			StartTime: 5002,
			EndTime:   10001,
			Service:   "test",
			Level:     "reader",
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
			Level:     "reader",
			Messages: []*Message{
				&Message{Text: "Foo2", Timestamp: 13000},
				&Message{Text: "Bar2", Timestamp: 13500},
				&Message{Text: "Baz2", Timestamp: 17000},
			},
		}
		b3 := &Block{
			StartTime: 30000,
			EndTime:   40003,
			Service:   "test",
			Level:     "reader",
			Messages: []*Message{
				&Message{Text: "Foo3", Timestamp: 30000},
				&Message{Text: "Bar3", Timestamp: 34000},
				&Message{Text: "Baz3", Timestamp: 40003},
			},
		}
		b4 := &Block{
			StartTime: 15000,
			EndTime:   25000,
			Service:   "test2",
			Level:     "should_not_be_included",
			Messages: []*Message{
				&Message{Text: "Foo3", Timestamp: 30000},
				&Message{Text: "Bar3", Timestamp: 34000},
				&Message{Text: "Baz3", Timestamp: 40003},
			},
		}
		b1.WriteToFile()
		b2.WriteToFile()
		b3.WriteToFile()
		b4.WriteToFile()
		blocks, err := GetBlocksInTimeRange(5001, 29999, "test", "reader", &CacheAccessorMock{})

		So(err, ShouldBeNil)

		So(len(blocks), ShouldEqual, 2)
		So(blocks[0], ShouldResemble, b1)
		So(blocks[1], ShouldResemble, b2)
	})
	os.RemoveAll(pathPrefix)
}
