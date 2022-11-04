dev: 
	$(GOPATH)/bin/air -c .air.conf

server_bot:
	go build  -gcflags=all="-l -B" -ldflags="-w -s" -o server_bot ./cmd/server_bot

test:
	go test -v ./...

build: server_bot
	xz server_bot

start: server_bot
	GIN_MODE=release ./server_bot
