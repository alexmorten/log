package log

import (
	"io/ioutil"
	"sort"
)

//CacheAccessor is something that ownes and return a cache
type CacheAccessor interface {
	GetCache() *Cache
}

func getFileNames(service, level string) (files []string, err error) {
	fileInfos, err := ioutil.ReadDir(BlockPath(service, level))
	for _, info := range fileInfos {
		files = append(files, info.Name())
	}
	return
}

//GetBlocksInTimeRange returns the blocks in the timerange
func GetBlocksInTimeRange(startTime, endTime int64, service, level string, a CacheAccessor) (blocks []*Block, err error) {
	cachedBlocks := a.GetCache().GetBlocks(service, level)

	if len(cachedBlocks) > 0 {
		for _, block := range cachedBlocks {
			if block.IsInTimeRange(startTime, endTime) {
				blocks = append(blocks, block)
			}
		}
	} else {
		fileNames, err := getFileNames(service, level)
		if err != nil {
			return nil, err
		}
		for _, fileName := range fileNames {
			if b, e := ParseFileNameIntoBlock(fileName); e == nil {
				if b.IsInTimeRange(startTime, endTime) {
					b.Service = service
					b.Level = level
					if e := b.ReadFromFile(); e == nil {
						blocks = append(blocks, b)
					}
				}
			}
		}
	}

	blocks = sortBlocks(blocks)
	if len(blocks) > 0 {
		blocks[0].ReduceToTimeRange(startTime, endTime)
	}
	if len(blocks) > 1 {
		blocks[len(blocks)-1].ReduceToTimeRange(startTime, endTime)
	}
	return
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
