# log

## run the server 

### by cloning this repo
 clone this repo and `go run cmd/server/main.go`

### by using docker
- you should create a volume for this: `docker volume create log-volume`
- run it with `docker run -p 7654:7654 -d --name log --mount source=log-volume,target=/app/data --rm alexmorten/log`

## client library 

### usage
```go

import (
	"github.com/alexmorten/log/client"
)

func main() {
	log := client.NewClient()
	log.Config.ServiceName = "some_service_name"
	log.Config.URL = "<server location url>"

	log.Log("Nothing special")
	log.LogWarn("This could potentially be dangerous")
	log.LogError("It was.")
	log.LogMessage("custom", "something not fitting into any predefined level")

	log.Commit() // can be used to commit the log messages early to the server
	log.Shutdown() //remember to shut the logger down, otherwise some messages could be lost on an ungraceful shutdown
}

```

## install the cli
- the command line is an easy way of seeing the logs that have been sent to the server
- `go get -u github.com/alexmorten/log/cmd/logcli`

### usage of the cli

`logcli [-service <service name>] [-level <level name> (needs service to be provided too)] [-url <url to the server>]`

## TODO
### server 
- [ ] split blocks per year/month/day for faster access over long periods of time
- [ ] merge blocks before/ after being written to disk, to handle more reasonable sizes (keep it close to a multiple of os block size)
- [ ] add a json endpoint to get messages from the browser for example ( or solve this through a proxy service?)

### cli
- [ ] parse/serialize human readable times 
- [ ] add interactive component, allowing to scroll through logs easiliy (with arrow keys for example)
- [ ] add color for different log levels
