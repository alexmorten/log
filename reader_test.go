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
	Convey("mergeOrderedMessageStacks", t, func() {
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
			StartTime: 8000,
			EndTime:   12000,
			Service:   "test",
			Level:     "reader3",
			Messages: []*Message{
				&Message{Text: "Foo2", Timestamp: 8000},
				&Message{Text: "Bar2", Timestamp: 9000},
				&Message{Text: "Baz2", Timestamp: 12000},
			},
		}
		b3 := &Block{
			StartTime: 30000,
			EndTime:   40003,
			Service:   "test2",
			Level:     "reader",
			Messages: []*Message{
				&Message{Text: "Foo3", Timestamp: 30000},
				&Message{Text: "Bar3", Timestamp: 34000},
				&Message{Text: "Baz3", Timestamp: 40003},
			},
		}
		stacks := []*MessageContainerStack{
			b1.toCompleteMessageStack(),
			b2.toCompleteMessageStack(),
			b3.toCompleteMessageStack(),
		}

		mergedStack := mergeOrderedMessageStacks(stacks)
		completeMessages := []*CompleteMessage{}
		for !mergedStack.Empty() {
			completeMessages = append(completeMessages, mergedStack.PopMessageContainer().(*CompleteMessage))
		}
		expectedMessages := []*CompleteMessage{
			&CompleteMessage{Service: "test", Level: "reader", Message: &Message{Text: "Foo", Timestamp: 5002}},
			&CompleteMessage{Service: "test", Level: "reader", Message: &Message{Text: "Bar", Timestamp: 7005}},
			&CompleteMessage{Service: "test", Level: "reader3", Message: &Message{Text: "Foo2", Timestamp: 8000}},
			&CompleteMessage{Service: "test", Level: "reader3", Message: &Message{Text: "Bar2", Timestamp: 9000}},
			&CompleteMessage{Service: "test", Level: "reader", Message: &Message{Text: "Baz", Timestamp: 10001}},
			&CompleteMessage{Service: "test", Level: "reader3", Message: &Message{Text: "Baz2", Timestamp: 12000}},
			&CompleteMessage{Service: "test2", Level: "reader", Message: &Message{Text: "Foo3", Timestamp: 30000}},
			&CompleteMessage{Service: "test2", Level: "reader", Message: &Message{Text: "Bar3", Timestamp: 34000}},
			&CompleteMessage{Service: "test2", Level: "reader", Message: &Message{Text: "Baz3", Timestamp: 40003}},
		}
		So(completeMessages, ShouldResemble, expectedMessages)

	})
	os.RemoveAll(pathPrefix)
}
