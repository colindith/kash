package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/colindith/kash/store"
)

const (
	connHost = "localhost"
	connPort = "3333"
	connType = "tcp"
)

var (
	respOK = []byte("OK\n")
)

var shardedMapStore store.Store

func main() {
	StartKashServer()
}

// StartKashServer is the entry point of the kash server
func StartKashServer() {
	initRouter()

	initStore()
	defer closeStore()

	go runTCPServer()
}

func runTCPServer() {
	l, err := net.Listen(connType, connHost + ":" + connPort)
	if err != nil {
		log.Fatalf("net_listen_error | err=%v", err.Error())
		return
	}
	defer l.Close()
	log.Printf("kash_server_listen_at | %v", connHost + ":" + connPort)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("net_accept_error | err=%v", err.Error())
			return
		}
		go handleConnection(c)
	}
}

func initStore() {
	shardedMapStore = store.GetShardedMapStore()
}

func closeStore() {

}

func handleConnection(c net.Conn) {
	log.Printf("serving_connection | addr=%v", c.RemoteAddr().String())
	defer c.Close()
	for {
		netData, err := bufio.NewReader(c).ReadBytes('\n')
		if err != nil {
			log.Fatalf("bufio_read_bytes_error | err=%v", err.Error())
			return
		}

		args := bytes.Split(netData[:len(netData)-1], []byte{' '})
		if string(args[0]) == "STOP" {
			break
		}

		result := respOK
		var errMsg string
		handler, ok := cmdHandlerRouter[string(args[0])]
		if !ok {
			result = []byte("cmd not recognized")
		} else {
			result, errMsg, ok = handler(args[1:]...)
			if !ok {
				result = []byte(errMsg)
			}
		}

		_, err = c.Write(result)
		if err != nil {
			log.Fatalf("net_connection_write_error | err=%v", err.Error())
		}
	}

}

type handlerFunc func(params... []byte) (resp []byte, errMsg string, ok bool)

var cmdHandlerRouter map[string]handlerFunc

func initRouter() {
	cmdHandlerRouter = map[string]handlerFunc{
		"GET": handleGETCmd,
		"SET": handleSETCmd,
	}
}

func handleGETCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	key := string(params[0])
	// TODO: handle other params
	value, err := shardedMapStore.Get(key)
	if err != nil {
		log.Printf("handler_get_cmd_failed | msg=%v", err.Error())
		return nil, err.Error(), false
	}
	resp = append(value.([]byte), byte('\n'))
	return resp, "", true
}

func handleSETCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 2 {
		return nil, "not enough parameters", false
	}
	key := string(params[0])
	if len(params) == 2 {
		err := shardedMapStore.Set(key, params[1])
		if err != nil {
			log.Printf("handler_set_cmd_failed | msg=%v", err.Error())
			return nil, err.Error(), false
		}
	} else if len(params) == 3 {
		timeout, err := strconv.Atoi(string(params[2]))
		if err != nil {
			log.Printf("handler_set_cmd_failed | msg=%v", err.Error())
			return nil, err.Error(), false
		}
		err = shardedMapStore.SetWithTimeout(key, params[1], time.Duration(timeout)*time.Second)
		if err != nil {
			log.Printf("handler_set_cmd_failed | msg=%v", err.Error())
			return nil, err.Error(), false
		}
	}
	// TODO: handle other params


	return respOK, "", true
}
