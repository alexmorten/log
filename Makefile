proto:
	protoc --go_out=./ protocol.proto

dep:
	glide install

test:
	go test ./... -timeout 10s

run:
	go run main/log.go
