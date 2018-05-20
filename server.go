package log

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
)

// Server handles incoming blocks
type Server struct {
	WriterCollection *WriterCollection
	Reader           *Reader
}

//NewDefaultServer creates a new Server and initializes its members
func NewDefaultServer() *Server {
	cache := NewCache()
	fileReader := &FileReader{}
	reader := NewReader(cache, fileReader)
	return &Server{
		Reader:           reader,
		WriterCollection: NewWriterCollection(cache),
	}
}

// StartServer starts a new Server
func StartServer() {
	s := NewDefaultServer()
	http.Handle("/", s)
	err := http.ListenAndServe(":7654", nil)
	if err != nil {
		fmt.Println(err)
	}
	s.Shutdown()
}

//Shutdown the server and all its components
func (s *Server) Shutdown() {
	s.WriterCollection.Shutdown()
	s.Reader.Shutdown()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	switch r.Method {
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodGet:
		s.handleGet(w, r)
	}
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	postRequest := &PostRequest{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = proto.Unmarshal(bytes, postRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, block := range postRequest.Blocks {
		if !block.Valid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	for _, block := range postRequest.Blocks {
		storageWriter := s.WriterCollection.GetWriter(block.Service, block.Level)
		storageWriter.InChannel <- block
	}

	w.WriteHeader(http.StatusOK)
}

type getParams struct {
	startTime int64
	endTime   int64
	service   string
	level     string
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	parsedParams, err := parseParams(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if parsedParams.service != "" && parsedParams.level != "" {
		s.handleServiceLevelGet(w, parsedParams)
	} else if parsedParams.service != "" && parsedParams.level == "" {
		s.handleServiceGet(w, parsedParams)
	} else {
		s.handlePlainGet(w, parsedParams)
	}

}

func (s *Server) handleServiceLevelGet(w http.ResponseWriter, p *getParams) {
	messages := s.Reader.GetServiceLevelMessagesInTimeRange(
		p.startTime,
		p.endTime,
		p.service,
		p.level,
	)

	response := pools.GetServiceLevelResponses.Get().(*GetServiceLevelResponse)
	response.Reset()
	response.Messages = messages

	bytes, err := proto.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)

	// put objects back into their pools
	for _, message := range messages {
		pools.PlainMessages.Put(message)
	}
	pools.GetServiceLevelResponses.Put(response)
}

func (s *Server) handleServiceGet(w http.ResponseWriter, p *getParams) {
	messages := s.Reader.GetServiceMessagesInTimeRange(
		p.startTime,
		p.endTime,
		p.service,
	)

	response := pools.GetServiceResponses.Get().(*GetServiceResponse)
	response.Reset()
	response.Messages = messages

	bytes, err := proto.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)

	// put objects back into their pools
	for _, message := range messages {
		pools.ServiceMessages.Put(message)
	}
	pools.GetServiceResponses.Put(response)

}

func (s *Server) handlePlainGet(w http.ResponseWriter, p *getParams) {
	messages := s.Reader.GetCompleteMessagesInTimeRange(
		p.startTime,
		p.endTime,
	)
	response := pools.GetResponses.Get().(*GetResponse)
	response.Reset()
	response.Messages = messages

	bytes, err := proto.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)

	// put objects back into their pools
	for _, message := range messages {
		pools.CompleteMessages.Put(message)
	}
	pools.GetResponses.Put(response)
}

func parseParams(params url.Values) (p *getParams, err error) {
	p = &getParams{}
	startTimeParam := params.Get("from_time")
	endTimeParam := params.Get("to_time")
	if endTimeParam == "" {
		p.endTime = time.Now().Unix()
	} else {
		p.endTime, err = strconv.ParseInt(endTimeParam, 10, 64)
	}
	if startTimeParam == "" {
		p.startTime = time.Unix(p.endTime-int64((time.Hour*1).Seconds()), 0).Unix()
	} else {
		p.endTime, err = strconv.ParseInt(endTimeParam, 10, 64)
	}
	p.service = params.Get("service")
	p.level = params.Get("level")
	return
}
