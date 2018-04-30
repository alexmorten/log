package log

import (
	"errors"
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
	Cache            *Cache
}

//NewServer creates a new Server and initializes its members
func NewServer() *Server {
	cache := NewCache()
	return &Server{
		Cache:            cache,
		WriterCollection: NewWriterCollection(cache),
	}
}

// StartServer starts a new Server
func StartServer() {
	s := NewServer()
	http.Handle("/", s)
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println(err)
	}
	s.Shutdown()
}

//Shutdown the server and all its components
func (s *Server) Shutdown() {
	s.WriterCollection.Shutdown()
	s.Cache.Shutdown()
}

//GetCache of server
func (s *Server) GetCache() *Cache {
	return s.Cache
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
	b := &Block{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = proto.Unmarshal(bytes, b)
	if err != nil || !b.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storageWriter := s.WriterCollection.GetWriter(b.Service, b.Level)
	storageWriter.InChannel <- b
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
	getParams, err := parseParams(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	blocks, err := GetBlocksInTimeRange(
		getParams.startTime,
		getParams.endTime,
		getParams.service,
		getParams.level,
		s,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messages := []*Message{}
	for _, block := range blocks {
		messages = append(messages, block.Messages...)
	}
	response := &GetResponse{Messages: messages}

	bytes, err := proto.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
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
	if p.service == "" {
		err = errors.New("No service param provided")
		return
	}
	p.level = params.Get("level")
	if p.level == "" {
		err = errors.New("No level param provided")
		return
	}

	return
}
