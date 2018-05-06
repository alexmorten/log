package client

import (
	"testing"

	"github.com/alexmorten/log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	Convey("Client Cache", t, func() {
		cache := NewCache()

		m1 := &log.Message{
			Timestamp: 10001,
			Text:      "foo",
		}
		m2 := &log.Message{
			Timestamp: 10001,
			Text:      "bar",
		}
		m3 := &log.Message{
			Timestamp: 10001,
			Text:      "baz",
		}
		m4 := &log.Message{
			Timestamp: 10001,
			Text:      "bab",
		}

		cache.AddMessage("level1", m1)
		cache.AddMessage("level1", m2)
		cache.AddMessage("level2", m3)
		cache.AddMessage("level3", m4)

		messagesMap := cache.GetCachedMessagesAndReset()

		So(messagesMap, ShouldNotBeEmpty)
		So(cache.messages, ShouldBeEmpty)
		So(messagesMap["level1"], ShouldResemble, []*log.Message{m1, m2})
		So(messagesMap["level2"], ShouldResemble, []*log.Message{m3})
		So(messagesMap["level3"], ShouldResemble, []*log.Message{m4})

	})
}
