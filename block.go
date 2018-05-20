package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/gogo/protobuf/proto"
)

// Merge b2 into b1, b2 should be a block that's younger
func (b *Block) Merge(b2 *Block) {
	b.Messages = append(b.Messages, b2.Messages...)
	b.EndTime = b2.EndTime
}

//Copy copies the block, the values inside the block still point to the same instances
func (b *Block) Copy() *Block {
	return &Block{
		StartTime: b.StartTime,
		EndTime:   b.EndTime,
		Level:     b.Level,
		Service:   b.Service,
		Messages:  b.Messages,
	}
}

//WriteToFile writes the block to disk
func (b *Block) WriteToFile() error {
	os.MkdirAll(b.path(), os.ModePerm)
	f, err := os.Create(b.path() + "/" + b.fileName())
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := proto.Marshal(b)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// Valid checks if the Block is valid
func (b *Block) Valid() bool {
	if b.StartTime > b.EndTime {
		return false
	}
	if len(b.Messages) == 0 {
		return false
	}
	if b.Service == "" || b.Level == "" {
		return false
	}
	return true

}

// ReadFromFile uses the start_time and end_time of itself to read the appropriate file and fill itself with the stored info
func (b *Block) ReadFromFile() (err error) {
	byteArray, err := ioutil.ReadFile(b.path() + "/" + b.fileName())
	if err != nil {
		return
	}
	err = proto.Unmarshal(byteArray, b)
	return
}

//BlockPath returns the Path where blocks for a given service and level are stored
func BlockPath(service, level string) string {
	return fmt.Sprintf("%v/%v", levelPath(service), level)
}

const maxUint = ^uint64(0)
const minUint = 0
const maxInt = int64(maxUint >> 1)
const minInt = -maxInt - 1

//ReduceToTimeRange modifies the block to only include messages in the provided timeRange
func (b *Block) ReduceToTimeRange(startTime, endTime int64) {
	messages := []*Message{}
	highestTimestamp := int64(0)
	lowestTimestamp := maxInt
	for _, message := range b.Messages {
		if message.IsInTimeRange(startTime, endTime) {
			messages = append(messages, message)

			if highestTimestamp < message.Timestamp {
				highestTimestamp = message.Timestamp
			}
			if lowestTimestamp > message.Timestamp {
				lowestTimestamp = message.Timestamp
			}
		}
	}
	b.StartTime = lowestTimestamp
	b.EndTime = highestTimestamp
	b.Messages = messages
}

//IsInTimeRange ...
func (b *Block) IsInTimeRange(startTime, endTime int64) bool {
	return b.StartTime <= endTime && b.EndTime >= startTime
}

//ParseFileNameIntoBlock creates a block with the start and end time from the filename.
//Use ReadFromFile to get the rest of the info
func ParseFileNameIntoBlock(filename string) (b *Block, err error) {
	pattern := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := pattern.FindStringSubmatch(filename)
	if len(matches) != 3 {
		err = errors.New("filename was invalid")
		return
	}
	startTime, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return
	}
	endTime, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return
	}
	b = &Block{
		StartTime: startTime,
		EndTime:   endTime,
	}
	return
}

func (b *Block) path() string {
	return BlockPath(b.Service, b.Level)
}

func (b *Block) fileName() string {
	return fmt.Sprintf("%v-%v", b.StartTime, b.EndTime)
}

func (b *Block) toPlainMessageStack() *MessageContainerStack {
	stack := &MessageContainerStack{}
	for _, message := range b.Messages {
		container := &PlainMessage{
			Message: message,
		}
		stack.PutMessageContainer(container)
	}
	return stack
}

func (b *Block) toServiceMessageStack() *MessageContainerStack {
	stack := &MessageContainerStack{}
	for _, message := range b.Messages {
		container := &ServiceMessage{
			Message: message,
			Level:   b.Level,
		}
		stack.PutMessageContainer(container)
	}
	return stack
}

func (b *Block) toCompleteMessageStack() *MessageContainerStack {
	stack := &MessageContainerStack{}
	for _, message := range b.Messages {
		container := &CompleteMessage{
			Message: message,
			Level:   b.Level,
			Service: b.Service,
		}
		stack.PutMessageContainer(container)
	}
	return stack
}
