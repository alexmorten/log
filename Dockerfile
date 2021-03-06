FROM golang:1.8 as builder
WORKDIR /go/src/github.com/alexmorten/log
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/alexmorten/log/server app
WORKDIR app
CMD ["./server"]
EXPOSE 7654
