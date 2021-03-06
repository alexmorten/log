package log

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	Convey("Cache", t, func() {
		cache := NewCache()
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

		cache.inChannel <- b1
		time.Sleep(10 * time.Millisecond)
		So(cache.GetBlock(5000, 12000, b1.Service, b1.Level), ShouldResemble, b1)

		cache.inChannel <- b2
		time.Sleep(10 * time.Millisecond)
		expectedBlock := &Block{
			StartTime: 5002,
			EndTime:   13500,
			Service:   "test",
			Level:     "reader",
			Messages:  append(b1.Messages, b2.Messages[0], b2.Messages[1]),
		}
		So(cache.GetBlock(5000, 14000, b1.Service, b1.Level), ShouldResemble, expectedBlock)

		cache.inChannel <- b3
		time.Sleep(10 * time.Millisecond)
		expectedBlock = &Block{
			StartTime: 10001,
			EndTime:   40003,
			Service:   "test",
			Level:     "reader",
			Messages:  append([]*Message{b1.Messages[2]}, append(b2.Messages, b3.Messages...)...),
		}
		So(cache.GetBlock(10000, 50000, b1.Service, b1.Level), ShouldResemble, expectedBlock)

		cache.Shutdown()
	})
}

func TestCacheCleaning(t *testing.T) {
	limitBefore := cacheMessageCountLimit

	Convey("cleanCache", t, func() {
		cacheMessageCountLimit = 5
		cache := NewCache()
		b1 := &Block{
			StartTime: 5002,
			EndTime:   10001,
			Service:   "test",
			Level:     "cache",
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
			Level:     "cache",
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
			Level:     "cache2",
			Messages: []*Message{
				&Message{Text: "Foo3", Timestamp: 30000},
			},
		}
		b4 := &Block{
			StartTime: 50000,
			EndTime:   50000,
			Service:   "test",
			Level:     "cache2",
			Messages: []*Message{
				&Message{Text: "FooBar2", Timestamp: 50000},
			},
		}
		cache.AddBlock(b1)
		cache.AddBlock(b2)
		cache.AddBlock(b3)
		cache.AddBlock(b4)
		time.Sleep(10 * time.Millisecond)
		So(cache.GetBlock(5000, 50000, "test", "cache"), ShouldResemble, b2)
		expectedBlock := b3.Copy()
		expectedBlock.Merge(b4)
		So(cache.GetBlock(5000, 50000, "test", "cache2"), ShouldResemble, expectedBlock)
	})

	cacheMessageCountLimit = limitBefore
}
