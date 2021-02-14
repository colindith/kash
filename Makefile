all:
	rm -r cmd/server
	rm -r cmd/client
	go build tcp/server.go -o cmd/server
	./cmd/server