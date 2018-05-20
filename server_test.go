package log

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	. "github.com/smartystreets/goconvey/convey"

	"net/http/httptest"
)

func TestPostEndpoint(t *testing.T) {
	Convey("PostEndpoint", t, func() {

		pathPrefix = "test"
		b := &Block{
			StartTime: 5002,
			EndTime:   10001,
			Service:   "test",
			Level:     "endpoint",
			Messages: []*Message{
				&Message{Text: "Foo", Timestamp: 5002},
				&Message{Text: "Bar", Timestamp: 7005},
				&Message{Text: "Baz", Timestamp: 10001},
			},
		}
		b2 := &Block{
			StartTime: 4999,
			EndTime:   9999,
			Service:   "test",
			Level:     "endpoint2",
			Messages: []*Message{
				&Message{Text: "Foob", Timestamp: 4999},
				&Message{Text: "Barb", Timestamp: 7005},
				&Message{Text: "Bazb", Timestamp: 9999},
			},
		}
		postRequest := &PostRequest{Blocks: []*Block{b, b2}}
		byteArray, _ := proto.Marshal(postRequest)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(byteArray))
		resp := httptest.NewRecorder()

		s := NewDefaultServer()
		s.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, 200)
		time.Sleep(10 * time.Millisecond)

		byteArray, _ = ioutil.ReadFile(pathPrefix + "/test/endpoint/5002-10001")
		outputBlock := &Block{}
		proto.Unmarshal(byteArray, outputBlock)
		So(outputBlock.Messages[0].Text, ShouldEqual, "Foo")

		byteArray, _ = ioutil.ReadFile(pathPrefix + "/test/endpoint2/4999-9999")
		outputBlock = &Block{}
		proto.Unmarshal(byteArray, outputBlock)
		So(outputBlock.Messages[0].Text, ShouldEqual, "Foob")
	})

	os.RemoveAll(pathPrefix)
}
func TestGetEndpoint(t *testing.T) {
	Convey("Get Endpoint", t, func() {

		pathPrefix = "test"
		b := &Block{
			StartTime: 5002,
			EndTime:   10001,
			Service:   "test",
			Level:     "endpoint",
			Messages: []*Message{
				&Message{Text: "Foo", Timestamp: 5002},
				&Message{Text: "Bar", Timestamp: 7005},
				&Message{Text: "Baz", Timestamp: 10001},
			},
		}
		b2 := &Block{
			StartTime: 5003,
			EndTime:   10002,
			Service:   "test",
			Level:     "endpoint2",
			Messages: []*Message{
				&Message{Text: "Foo2", Timestamp: 5003},
				&Message{Text: "Bar2", Timestamp: 7006},
				&Message{Text: "Baz2", Timestamp: 10002},
			},
		}
		b3 := &Block{
			StartTime: 5004,
			EndTime:   10004,
			Service:   "test2",
			Level:     "endpoint2",
			Messages: []*Message{
				&Message{Text: "Foo3", Timestamp: 5004},
				&Message{Text: "Bar3", Timestamp: 7007},
				&Message{Text: "Baz3", Timestamp: 10003},
			},
		}
		b.WriteToFile()
		b2.WriteToFile()
		b3.WriteToFile()
		Convey("given a service and level gets the correct logs", func() {

			url := "/?from_time=5001&to_time=8008&service=test&level=endpoint"
			req := httptest.NewRequest("GET", url, bytes.NewReader([]byte{}))
			resp := httptest.NewRecorder()

			s := NewDefaultServer()
			s.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, 200)

			byteArray, _ := ioutil.ReadAll(resp.Body)
			response := &GetServiceLevelResponse{}
			err := proto.Unmarshal(byteArray, response)
			So(err, ShouldBeNil)
			So(len(response.Messages), ShouldEqual, 2)
			So(response.Messages[0].GetLogMessage().Text, ShouldEqual, "Foo")
			So(response.Messages[1].GetLogMessage().Text, ShouldEqual, "Bar")
		})

		Convey("given only a service gets the correct logs", func() {
			url := "/?from_time=5001&to_time=8008&service=test"
			req := httptest.NewRequest("GET", url, bytes.NewReader([]byte{}))
			resp := httptest.NewRecorder()

			s := NewDefaultServer()
			s.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, 200)

			byteArray, _ := ioutil.ReadAll(resp.Body)
			response := &GetServiceResponse{}
			err := proto.Unmarshal(byteArray, response)
			So(err, ShouldBeNil)
			messages := response.Messages
			So(len(response.Messages), ShouldEqual, 4)
			So(messages[0].GetLogMessage().Text, ShouldEqual, "Foo")
			So(messages[1].GetLogMessage().Text, ShouldEqual, "Foo2")
			So(messages[2].GetLogMessage().Text, ShouldEqual, "Bar")
			So(messages[3].GetLogMessage().Text, ShouldEqual, "Bar2")
		})

		Convey("given no service or level gets the correct logs", func() {
			url := "/?from_time=5001&to_time=8008"
			req := httptest.NewRequest("GET", url, bytes.NewReader([]byte{}))
			resp := httptest.NewRecorder()

			s := NewDefaultServer()
			s.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, 200)

			byteArray, _ := ioutil.ReadAll(resp.Body)
			response := &GetResponse{}
			err := proto.Unmarshal(byteArray, response)
			So(err, ShouldBeNil)
			messages := response.Messages
			So(len(response.Messages), ShouldEqual, 6)
			So(messages[0].GetLogMessage().Text, ShouldEqual, "Foo")
			So(messages[1].GetLogMessage().Text, ShouldEqual, "Foo2")
			So(messages[2].GetLogMessage().Text, ShouldEqual, "Foo3")
			So(messages[3].GetLogMessage().Text, ShouldEqual, "Bar")
			So(messages[4].GetLogMessage().Text, ShouldEqual, "Bar2")
			So(messages[5].GetLogMessage().Text, ShouldEqual, "Bar3")

		})
	})

	os.RemoveAll(pathPrefix)
}

func TestGetEndpointCache(t *testing.T) {
	Convey("Get Endpoint Caching", t, func() {
		pathPrefix = "test"
		b := &Block{
			StartTime: 5002,
			EndTime:   10001,
			Service:   "test",
			Level:     "endpoint",
			Messages: []*Message{
				&Message{Text: "Foo", Timestamp: 5002},
				&Message{Text: "Bar", Timestamp: 7005},
				&Message{Text: "Baz", Timestamp: 10001},
			},
		}
		postRequest := &PostRequest{Blocks: []*Block{b}}
		byteArray, _ := proto.Marshal(postRequest)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(byteArray))
		resp := httptest.NewRecorder()

		s := NewDefaultServer()
		s.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, 200)
		time.Sleep(10 * time.Millisecond)
		// Remove from disk
		os.RemoveAll(pathPrefix)

		url := "/?from_time=5001&to_time=12000&service=test&level=endpoint"
		getReq := httptest.NewRequest("GET", url, bytes.NewReader([]byte{}))
		getResp := httptest.NewRecorder()

		s.ServeHTTP(getResp, getReq)

		So(resp.Code, ShouldEqual, 200)

		time.Sleep(10 * time.Microsecond)

		getByteArray, _ := ioutil.ReadAll(getResp.Body)
		response := &GetServiceLevelResponse{}
		err := proto.Unmarshal(getByteArray, response)

		So(err, ShouldBeNil)

		expectedMessages := []*PlainMessage{
			&PlainMessage{Message: &Message{Text: "Foo", Timestamp: 5002}},
			&PlainMessage{Message: &Message{Text: "Bar", Timestamp: 7005}},
			&PlainMessage{Message: &Message{Text: "Baz", Timestamp: 10001}},
		}
		So(response.Messages, ShouldResemble, expectedMessages)
	})
	os.RemoveAll(pathPrefix)
}
