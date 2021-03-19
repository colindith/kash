package main

import (
	"log"

	"github.com/colindith/kash/cmd_reader"
	"github.com/colindith/kash/tcp"
)

// TODO: The client package should start a cmd line process and let users type the cmd in the cmd line.
// TODO: The cmd line should be provide some basic functions like using upper arrow key to find the history and and using left/right arrow to move the cursor

const (           // TODO: the host & port should be passed in from outside
	connHost = "localhost"
	connPort = "3333"
)

func main() {
	cfg := &cmd_reader.Config{}
	cfg.SetPromptStr("kash>")    // TODO: The prompt should be get from the remote domain/IP address
	if err := cfg.RegistryHandler(handler); err != nil {
		log.Fatalf("registor cmd reader handler fail | err=%v", err.Error())
		return
	}

	// TODO: should connect to remote server before running the cli
	defer tcp.CloseConn()
	cmd_reader.Run(cfg)
}

var handler = &cmd_reader.Handler{Serv: func(cmd string) (result string, err error) {
	result = tcp.SendTCPCmd(connHost, connPort, cmd)
	return result, nil
}}
