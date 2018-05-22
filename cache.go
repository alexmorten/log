package log

import (
	"sync"
)

var cacheMessageCountLimit = 1000000

// Cache of blocks for instant access
type Cache struct {
	blocks          map[string]map[string][]*Block
	mutex           sync.Mutex
	inChannel       chan *Block
	shutdownChannel chan struct{}
	messageCounter  int
}

//NewCache ...
func NewCache() *Cache {
	cache := &Cache{
		blocks:          map[string]map[string][]*Block{},
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

//GetBlock for service and level
func (c *Cache) GetBlock(startTime, endTime int64, service, level string) *Block {
	blocks := []*Block{}
	for _, block := range c.blocks[service][level] {
		if block.IsInTimeRange(startTime, endTime) {
			blocks = append(blocks, block)
		}
	}
	if len(blocks) == 0 {
		return nil
	}

	blocks = sortBlocks(blocks)
	if len(blocks) > 0 {
		blocks[0].ReduceToTimeRange(startTime, endTime)
	}
	if len(blocks) > 1 {
		blocks[len(blocks)-1].ReduceToTimeRange(startTime, endTime)
	}

	mergedBlock := blocks[0].Copy()
	for i := 1; i < len(blocks); i++ {
		mergedBlock.Merge(blocks[i])
	}

	return mergedBlock
}

//GetLevels for a given service
func (c *Cache) GetLevels(service string) (levels []string) {
	for level := range c.blocks[service] {
		levels = append(levels, level)
	}
	return
}

//GetServices that have messages in the cache
func (c *Cache) GetServices() (services []string) {
	for serviceName := range c.blocks {
		services = append(services, serviceName)
	}
	return
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

	//make sure the inner map is initialized as well
	if c.blocks[b.Service] == nil {
		c.blocks[b.Service] = map[string][]*Block{}
	}
	current := c.blocks[b.Service][b.Level]

	c.blocks[b.Service][b.Level] = append(current, b)
	c.cleanCache()
}

func (c *Cache) cleanCache() {
	if c.messageCounter > cacheMessageCountLimit {
		for serviceName, levelToBlockMap := range c.blocks {
			for level, blocks := range levelToBlockMap {
				if len(blocks) == 0 {
					continue
				}
				newBlocks := []*Block{}
				c.messageCounter -= len(blocks[0].Messages)
				for i := 1; i < len(blocks); i++ {
					newBlocks = append(newBlocks, blocks[i])
				}
				c.blocks[serviceName][level] = newBlocks
			}
		}
	}
}
