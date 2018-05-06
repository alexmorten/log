FROM golang:1.8 as builder
WORKDIR /go/src/github.com/alexmorten/log
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o log main/log.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/alexmorten/log/log .
CMD [ "./log" ]
EXPOSE 7654
