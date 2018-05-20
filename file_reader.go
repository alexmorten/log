package log

import (
	"fmt"
	"io/ioutil"
)

var pathPrefix = "data"

//FileReader handles reading messages from the filesystem
type FileReader struct{}

//GetBlock for given service ,level and timerange
func (f *FileReader) GetBlock(startTime, endTime int64, service, level string) *Block {
	blocks := []*Block{}
	fileNames := getFileNames(service, level)
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
func (f *FileReader) GetLevels(service string) (levels []string) {
	dirInfos, err := ioutil.ReadDir(levelPath(service))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range dirInfos {
		levels = append(levels, info.Name())
	}
	return
}

//GetServices that have messages in the cache
func (f *FileReader) GetServices() (services []string) {
	dirInfos, err := ioutil.ReadDir(servicePath())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range dirInfos {
		services = append(services, info.Name())
	}
	return
}

//Shutdown for the Store interface
func (f *FileReader) Shutdown() {}

func getFileNames(service, level string) (files []string) {
	fileInfos, err := ioutil.ReadDir(BlockPath(service, level))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range fileInfos {
		files = append(files, info.Name())
	}
	return
}

func servicePath() string {
	return pathPrefix
}

func levelPath(service string) string {
	return fmt.Sprintf("%v/%v", servicePath(), service)
}
