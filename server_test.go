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
		byteArray, _ := proto.Marshal(b)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(byteArray))
		resp := httptest.NewRecorder()

		s := NewServer()
		s.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, 200)
		time.Sleep(10 * time.Millisecond)

		byteArray, _ = ioutil.ReadFile(pathPrefix + "/test/endpoint/5002-10001")
		outputBlock := &Block{}
		proto.Unmarshal(byteArray, outputBlock)
		So(outputBlock.Messages[0].Text, ShouldEqual, "Foo")
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
		b.WriteToFile()
		url := "/?from_time=5001&to_time=8008&service=test&level=endpoint"
		req := httptest.NewRequest("GET", url, bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()

		s := NewServer()
		s.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, 200)

		time.Sleep(10 * time.Millisecond)

		byteArray, _ := ioutil.ReadAll(resp.Body)
		response := &GetResponse{}
		err := proto.Unmarshal(byteArray, response)
		So(err, ShouldBeNil)
		So(response.Messages[0].Text, ShouldEqual, "Foo")
	})

	os.RemoveAll(pathPrefix)
}

func TestGetEndpointCache(t *testing.T) {
	Convey("Get Endpoint Cache", t, func() {
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

		byteArray, _ := proto.Marshal(b)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(byteArray))
		resp := httptest.NewRecorder()

		s := NewServer()
		s.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, 200)
		time.Sleep(10 * time.Millisecond)
		// Remove from disk
		os.RemoveAll(pathPrefix)

		url := "/?from_time=5001&to_time=8008&service=test&level=endpoint"
		getReq := httptest.NewRequest("GET", url, bytes.NewReader([]byte{}))
		getResp := httptest.NewRecorder()

		s.ServeHTTP(getResp, getReq)

		So(resp.Code, ShouldEqual, 200)

		time.Sleep(10 * time.Microsecond)

		getByteArray, _ := ioutil.ReadAll(getResp.Body)
		response := &GetResponse{}
		err := proto.Unmarshal(getByteArray, response)

		So(err, ShouldBeNil)

		So(response.Messages, ShouldResemble, []*Message{b.Messages[0], b.Messages[1]})
	})
	os.RemoveAll(pathPrefix)
}
