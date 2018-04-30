package log

import (
	"fmt"
)

// Writer is responsible for writing blocks to disk for one service and level combination
type Writer struct {
	Service         string
	Level           string
	cache           *Cache
	InChannel       chan *Block
	shutdownChannel chan struct{}
}

// NewWriter creates a newWriter
func NewWriter(service, level string, cache *Cache) *Writer {
	w := &Writer{
		Service:         service,
		Level:           level,
		cache:           cache,
		InChannel:       make(chan *Block, 1),
		shutdownChannel: make(chan struct{}, 1),
	}
	w.Run()
	return w
}

//HashKey for identifying a writer
func (w *Writer) HashKey() string {
	return WriterKeyFor(w.Service, w.Level)
}

// Run the writer in a goroutine
func (w *Writer) Run() {
	go w.listen()
}

//Shutdown the writer
func (w *Writer) Shutdown() {
	w.shutdownChannel <- struct{}{}
}

func (w *Writer) listen() {
	for {
		select {
		case block := <-w.InChannel:
			w.handleNewBlock(block)
		}

	}
}

func (w *Writer) handleNewBlock(block *Block) {
	err := block.WriteToFile()
	if err != nil {
		fmt.Println(err)
	} else {
		w.cache.AddBlock(block)
	}
}

//WriterKeyFor generates an identifier for a writer
func WriterKeyFor(service, level string) string {
	return service + "/" + level
}
