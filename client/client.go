package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/alexmorten/log"
	"github.com/gogo/protobuf/proto"
)

//Client for communicating to server
type Client struct {
	Config          *Config
	Cache           *Cache
	shutdownChannel chan struct{}
}

//NewClient with default config
func NewClient() *Client {
	c := &Client{
		Config: NewConfig(),
		Cache:  NewCache(),
	}
	go c.pushMessagesPeriodically()
	return c
}

//NewClientWithConfig with given config
func NewClientWithConfig(config Config) *Client {
	c := &Client{
		Config: &config,
		Cache:  NewCache(),
	}
	go c.pushMessagesPeriodically()
	return c
}

//LogMessage writes log message on any given level
func (c *Client) LogMessage(level, message string) {
	if level == "" || message == "" {
		panic("both level and message should be set when writing a message")
	}
	m := &log.Message{
		Text:      message,
		Timestamp: time.Now().Unix(),
	}
	c.Cache.AddMessage(level, m)
}

//Log standard Message
func (c *Client) Log(messageArgs ...interface{}) {
	c.LogMessage("standard", constructMessage(messageArgs...))
}

//LogError Message
func (c *Client) LogError(messageArgs ...interface{}) {
	c.LogMessage("error", constructMessage(messageArgs...))
}

//LogWarn standard Message
func (c *Client) LogWarn(messageArgs ...interface{}) {
	c.LogMessage("warning", constructMessage(messageArgs...))
}

//Commit the cache
func (c *Client) Commit() {
	c.pushMessages()
}

func (c *Client) pushMessagesPeriodically() {
	ticker := time.NewTicker(c.Config.SyncTime)
loop:
	for {
		select {
		case <-ticker.C:
			c.pushMessages()
		case <-c.shutdownChannel:
			break loop
		}
	}
	//release ticker resources
	ticker.Stop()
}

//Shutdown the Client
func (c *Client) Shutdown() {
	c.shutdownChannel <- struct{}{}
	c.Commit()
}

func (c *Client) pushMessages() {
	messagesMap := c.Cache.GetCachedMessagesAndReset()
	if len(messagesMap) == 0 {
		return
	}
	blocks := []*log.Block{}
	for level, messageArray := range messagesMap {
		if len(messageArray) == 0 {
			continue
		}
		block := &log.Block{
			Messages:  messageArray,
			StartTime: messageArray[0].Timestamp,
			EndTime:   messageArray[len(messageArray)-1].Timestamp,
			Service:   c.Config.ServiceName,
			Level:     level,
		}
		blocks = append(blocks, block)
	}

	if len(blocks) == 0 {
		return
	}
	request := &log.PostRequest{
		Blocks: blocks,
	}
	byteArr, err := proto.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	reader := bytes.NewReader(byteArr)
	resp, err := http.Post(c.Config.URL, "application/proto", reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode, "was returned") //TODO make this a proper log
	}
}

// we want spaces between each argument, but no new line
func constructMessage(args ...interface{}) string {
	message := ""
	for index, item := range args {
		message += fmt.Sprint(item)
		if index != len(args)-1 {
			message += " "
		}
	}
	return message
}
