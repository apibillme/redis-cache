# build from vendor directory
FROM golang:1.11
WORKDIR /go/src/github.com/apibillme/rediscache/
ADD . .
RUN set -ex && \      
  CGO_ENABLED=0 go build \
        -tags netgo \
        -v -a \
        -ldflags '-extldflags "-static"' \
        -o rediscache

# use 2 step process
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/apibillme/rediscache/rediscache .
ENTRYPOINT ["./rediscache"]
CMD [ "--ttl", "360", "--keycap", "500", "--tcp", "8000", "--address", "localhost:6379" ]
EXPOSE 8000
