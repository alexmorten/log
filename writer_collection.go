package log

import (
	"sync"
)

//WriterCollection handles the writers for the
type WriterCollection struct {
	writers map[string]*Writer
	mutex   sync.Mutex
	cache   *Cache
}

//NewWriterCollection creates a new thread safe collection of writers
func NewWriterCollection(cache *Cache) *WriterCollection {
	return &WriterCollection{
		writers: map[string]*Writer{},
		mutex:   sync.Mutex{},
		cache:   cache,
	}
}

//GetWriter thread-safely returns a writer for the service and level out of the collection or adds one if necessary.
func (c *WriterCollection) GetWriter(service, level string) *Writer {

	//no mutex here, so we can read already created writers in parallel
	storedWriter := c.getStoredWriter(service, level)
	if storedWriter != nil {
		return storedWriter
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	//make sure that we don't create a writer twice when we don't get a writer above
	storedWriter = c.getStoredWriter(service, level)
	if storedWriter != nil {
		return storedWriter
	}

	// Now we can be sure that we don't have a Writer in the collection
	writer := NewWriter(service, level, c.cache)
	c.writers[writer.HashKey()] = writer
	return writer
}

//Shutdown all writers
func (c *WriterCollection) Shutdown() {
	for _, writer := range c.writers {
		writer.Shutdown()
	}
}

func (c *WriterCollection) getCache() *Cache {
	return c.cache
}

func (c *WriterCollection) getStoredWriter(service, level string) *Writer {
	key := WriterKeyFor(service, level)
	return c.writers[key]
}
