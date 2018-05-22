# log

## run the server 

### cloning repo
 clone this repo and `go run cmd/server/main.go`

### using docker
- you should create a volume for this: `docker volume create log-volume`
- run it with `docker run -p 7654:7654 -d --name log --mount source=log-volume,target=/app/data --rm alexmorten/log`

## install the cli
- the command line is an easy way of seeing the logs that have been sent to the server
- `go get -u github.com/alexmorten/log/cmd/logcli`

### usage of the cli

`logcli [-service <service name>] [-level <level name> (needs service to be provided too)] [-url <url to the server>]`
