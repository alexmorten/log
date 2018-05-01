package log

import (
	"sync"
)

var cacheMessageCountLimit = 1000000

// Cache of blocks for instant access
type Cache struct {
	blocks          map[string][]*Block
	mutex           sync.Mutex
	inChannel       chan *Block
	shutdownChannel chan struct{}
	messageCounter  int
}

//NewCache ...
func NewCache() *Cache {
	cache := &Cache{
		blocks:          map[string][]*Block{},
		mutex:           sync.Mutex{},
		inChannel:       make(chan *Block),
		shutdownChannel: make(chan struct{}),
	}
	go cache.listenForBlocks()

	return cache
}

//InChannel is the input Channel for the cache
func (c *Cache) InChannel() chan *Block {
	return c.inChannel
}

//AddBlock to Cache
func (c *Cache) AddBlock(b *Block) {
	c.InChannel() <- b
}

//GetBlocks for service and level
func (c *Cache) GetBlocks(service, level string) []*Block {
	return c.blocks[BlockPath(service, level)]
}

//Shutdown the cache
func (c *Cache) Shutdown() {
	c.shutdownChannel <- struct{}{}
}

func (c *Cache) listenForBlocks() {
loop:
	for {
		select {
		case block := <-c.inChannel:
			c.handleAddBlock(block)
		case <-c.shutdownChannel:
			break loop
		}
	}
}

func (c *Cache) handleAddBlock(b *Block) {
	c.messageCounter += len(b.Messages)

	key := b.path()
	current := c.blocks[key]

	c.blocks[key] = append(current, b)
	c.cleanCache()
}

func (c *Cache) cleanCache() {
	if c.messageCounter > cacheMessageCountLimit {
		for key, blocks := range c.blocks {
			if len(blocks) == 0 {
				continue
			}
			newBlocks := []*Block{}
			c.messageCounter -= len(blocks[0].Messages)
			for i := 1; i < len(blocks); i++ {
				newBlocks = append(newBlocks, blocks[i])
			}
			if len(newBlocks) > 0 {
				c.blocks[key] = newBlocks
			} else {
				c.blocks[key] = []*Block{}
			}
		}
	}
}
