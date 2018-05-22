package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/alexmorten/log"
	"github.com/gogo/protobuf/proto"
)

var service, level, serverURL string

func main() {
	flag.StringVar(&service, "service", "", "restrict output to log messages from the provided service")
	flag.StringVar(&level, "level", "", "restrict output to log messages on the provided level (can only be used together with a service)")
	flag.StringVar(&serverURL, "url", "http://localhost:7654", "url of the log server")
	flag.Parse()
	u, err := url.Parse(serverURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	params := u.Query()
	if service != "" {
		params.Set("service", service)
		if level != "" {
			params.Set("level", level)
		}
	} else if level != "" {
		fmt.Println("you can only use the level flag if you also provide a service!")
		return
	}
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Status)
	if service != "" && level != "" {
		handleServiceLevelResponse(resp)
	} else if service != "" {
		handleServiceResponse(resp)
	} else {
		handleCompleteResponse(resp)
	}
}

func handleCompleteResponse(r *http.Response) {
	response := &log.GetResponse{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := proto.Unmarshal(bytes, response); err != nil {
		fmt.Println(err)
		return
	}

	for _, message := range response.Messages {
		fmt.Printf("%v | %v | %v : %v \n", message.Message.Timestamp, message.Service, message.Level, message.Message.Text)
	}
}

func handleServiceResponse(r *http.Response) {
	response := &log.GetServiceResponse{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := proto.Unmarshal(bytes, response); err != nil {
		fmt.Println(err)
		return
	}

	for _, message := range response.Messages {
		fmt.Printf("%v | %v : %v \n", message.Message.Timestamp, message.Level, message.Message.Text)
	}
}

func handleServiceLevelResponse(r *http.Response) {
	response := &log.GetServiceLevelResponse{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := proto.Unmarshal(bytes, response); err != nil {
		fmt.Println(err)
		return
	}

	for _, message := range response.Messages {
		fmt.Printf("%v : %v \n", message.Message.Timestamp, message.Message.Text)
	}
}
