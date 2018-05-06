package client

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("Client", t, func() {
		cache := NewCache()
		client := &Client{
			Cache:           cache,
			Config:          NewConfig(),
			shutdownChannel: make(chan struct{}),
		}

		client.Log("Foo", "Bar", "Baz", []int{1, 2, 3, 4})
		client.Log("Bla", "Blap")

		messagesMap := client.Cache.GetCachedMessagesAndReset()
		So(len(messagesMap["standard"]), ShouldEqual, 2)
		So(messagesMap["standard"][0].Text, ShouldResemble, "Foo Bar Baz [1 2 3 4]")
		So(messagesMap["standard"][1].Text, ShouldResemble, "Bla Blap")
	})
}
