package log

import (
	"sort"
)

//CacheAccessor is something that ownes and return a cache
type CacheAccessor interface {
	GetCache() *Cache
}

//Store handles the retrival of blocks
type Store interface {
	GetBlock(startTime, endTime int64, service, level string) *Block

	GetLevels(service string) (levels []string)

	GetServices() (services []string)

	Shutdown()
}

//Reader reads blocks from the given Store levels
type Reader struct {
	Stores []Store
}

//NewReader creates a reader that reads from the provided Store levels
func NewReader(stores ...Store) *Reader {
	return &Reader{
		Stores: stores,
	}
}

//GetServiceLevelMessagesInTimeRange returns the blocks in the timerange
//if no blocks are found in the first Store level, the next ones are tried in order
func (r *Reader) GetServiceLevelMessagesInTimeRange(startTime, endTime int64, service, level string) (messages []*PlainMessage) {
	var block *Block
	for _, store := range r.Stores {
		block = store.GetBlock(startTime, endTime, service, level)
		if block != nil {
			break
		}
	}
	plainMessageStack := block.toPlainMessageStack()
	plainMessageStack.Flip()
	for plainMessageStack.Peek() != nil {
		messages = append(messages, plainMessageStack.PopMessageContainer().(*PlainMessage))
	}
	return
}

//GetServiceMessagesInTimeRange returns the blocks in the timerange
//if no blocks are found in the first Store level, the next ones are tried in order
func (r *Reader) GetServiceMessagesInTimeRange(startTime, endTime int64, service string) (messages []*ServiceMessage) {
	stackPerLevel := []*MessageContainerStack{}
	for _, store := range r.Stores {
		levels := store.GetLevels(service)
		for _, level := range levels {
			block := store.GetBlock(startTime, endTime, service, level)
			if block != nil {
				stackPerLevel = append(stackPerLevel, block.toServiceMessageStack())
			}
		}
		if len(stackPerLevel) > 0 {
			break
		}
	}

	mergedStack := mergeOrderedMessageStacks(stackPerLevel)
	for !mergedStack.Empty() {
		messages = append(messages, mergedStack.PopMessageContainer().(*ServiceMessage))
	}
	return
}

//GetCompleteMessagesInTimeRange returns the blocks in the timerange
//if no blocks are found in the first Store level, the next ones are tried in order
func (r *Reader) GetCompleteMessagesInTimeRange(startTime, endTime int64) (messages []*CompleteMessage) {
	stackPerServiceAndLevel := []*MessageContainerStack{}
	for _, store := range r.Stores {
		services := store.GetServices()
		for _, service := range services {
			levels := store.GetLevels(service)
			for _, level := range levels {
				block := store.GetBlock(startTime, endTime, service, level)
				if block != nil {
					stackPerServiceAndLevel = append(stackPerServiceAndLevel, block.toCompleteMessageStack())
				}
			}
		}

		if len(stackPerServiceAndLevel) > 0 {
			break
		}
	}

	mergedStack := mergeOrderedMessageStacks(stackPerServiceAndLevel)
	for !mergedStack.Empty() {
		messages = append(messages, mergedStack.PopMessageContainer().(*CompleteMessage))
	}
	return
}

//Shutdown the stores
func (r *Reader) Shutdown() {
	for _, store := range r.Stores {
		store.Shutdown()
	}
}

type blockCollection []*Block

func sortBlocks(blocks []*Block) []*Block {
	var collection blockCollection = blocks
	sort.Sort(collection)
	return collection
}

func (c blockCollection) Len() int {
	return len(c)
}

func (c blockCollection) Less(i, j int) bool {
	return c[i].StartTime < c[j].StartTime
}

func (c blockCollection) Swap(i, j int) {
	tempPointer := c[i]
	c[i] = c[j]
	c[j] = tempPointer
}

//expects stacks not to be empty
//stacks should have newest messages on top
func mergeOrderedMessageStacks(orderedContainerStacks []*MessageContainerStack) *MessageContainerStack {
	mergedStack := &MessageContainerStack{}
	if len(orderedContainerStacks) == 0 {
		return mergedStack
	}
	indexOfHighest := 0
	var highestTimestamp, secondHighestTimestamp int64

loop:
	for {
		allEmpty := true
		for index, stack := range orderedContainerStacks {
			if !stack.Empty() {
				allEmpty = false

				m := stack.PeekMessageContainer().GetLogMessage()
				if m.Timestamp >= highestTimestamp {
					secondHighestTimestamp = highestTimestamp
					highestTimestamp = m.Timestamp
					indexOfHighest = index
				} else if m.Timestamp > secondHighestTimestamp {
					secondHighestTimestamp = m.Timestamp
				}
			}
		}
		if allEmpty {
			break loop
		}
		stackOfHighest := orderedContainerStacks[indexOfHighest]
		mergedStack.PutMessageContainer(stackOfHighest.PopMessageContainer())

		// we know the second highest timestamp,
		//  so we can check against the next message in the stack, that we just popped a message of
		for !stackOfHighest.Empty() &&
			stackOfHighest.PeekMessageContainer().GetLogMessage().Timestamp >= secondHighestTimestamp {
			mergedStack.PutMessageContainer(stackOfHighest.PopMessageContainer())
		}

		highestTimestamp = 0
		secondHighestTimestamp = 0
	}
	return mergedStack
}
