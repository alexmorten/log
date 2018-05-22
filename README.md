# log

## run the server 

- clone this repo
- `go run cmd/server/main.go`

## install the cli
the command line is an easy way of seeing the logs that have been sent to the server
`go get -u github.com/alexmorten/log/cmd/logcli`

### usage of the cli

`logcli [-service <service name>] [-level <level name> (needs service to be provided too)] [-url <url to the server>]`
