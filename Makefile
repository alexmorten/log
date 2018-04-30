proto:
	protoc --go_out=./ protocol.proto

test:
	go test ./... -timeout 10s

run:
	go run main/log.go
