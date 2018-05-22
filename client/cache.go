package client

import (
	"sync"

	"github.com/alexmorten/log"
)

//Cache for client,
// we don't  need to send a request for each message since messages
// tend to happen in bunches, with mixed log levels
type Cache struct {
	sync.Mutex
	messages map[string][]*log.Message
}

//NewCache with defaults
func NewCache() *Cache {
	return &Cache{
		messages: make(map[string][]*log.Message),
	}
}

//AddMessage to Cache
func (c *Cache) AddMessage(level string, message *log.Message) {
	c.Lock()
	defer c.Unlock()

	c.messages[level] = append(c.messages[level], message)
}

//GetCachedMessagesAndReset gets the cache and resets it to an empty cache
func (c *Cache) GetCachedMessagesAndReset() map[string][]*log.Message {
	c.Lock()
	defer c.Unlock()
	messages := c.messages
	c.messages = make(map[string][]*log.Message)
	return messages
}
