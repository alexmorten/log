proto:
	protoc --go_out=./ protocol.proto

install:
	curl https://glide.sh/get | sh

dep:
	glide install

test:
	go test ./... -timeout 10s

run:
	go run cmd/server/main.go

image-build:
	docker build -t alexmorten/log .
