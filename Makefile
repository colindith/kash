all:
	rm -rf bin/server
	rm -rf bin/client
	go build -o bin/server server/server.go
	go build -o bin/client client/client.go