all:
	rm -rf bin/server
	rm -rf bin/client
	go build -o bin/server tcp/server/server.go
	go build -o bin/client tcp/client/client.go