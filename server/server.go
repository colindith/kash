package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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
	StartKashServer(connPort)
}

// StartKashServer is the entry point of the kash server
func StartKashServer(port string) {
	initRouter()

	initStore()

	runTCPServer(port)
}

func runTCPServer(port string) {
	l, err := net.Listen(connType, connHost + ":" + port)
	if err != nil {
		log.Fatalf("net_listen_error | err=%v", err.Error())
		return
	}
	defer l.Close()
	defer closeStore()
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
	shardedMapStore.Close()
}

func handleConnection(c net.Conn) {
	log.Printf("serving_connection | addr=%v", c.RemoteAddr().String())
	defer c.Close()
	for {
		netData, err := bufio.NewReader(c).ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("bufio_read_bytes_error | err=%v", err.Error())    // TODO: This is client input problem. Should not be error
			return
		}

		args := bytes.Split(bytes.Trim(netData[:len(netData)-1], " "), []byte{' '})
		if string(args[0]) == "STOP" {
			break
		}

		result := respOK
		var errMsg string
		handler, ok := cmdHandlerRouter[strings.ToUpper(string(args[0]))]
		if !ok {
			result = []byte("cmd not recognized")
		} else {
			result, errMsg, ok = handler(args[1:]...)
			if !ok {
				result = []byte(errMsg)
			}
		}
		result = append(result, byte('\n'))

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
		"GET":  handleGETCmd,
		"SET":  handleSETCmd,
		"DEL":  handleDELCmd,
		"INCR": handleINCRCmd,
		"DUMP": handleDUMPALLCmd,
	}
}

func handleGETCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 1 {
		return nil, "not enough parameters", false
	}
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
			log.Printf("set_cmd_failed | msg=%v", err.Error())
			return nil, err.Error(), false
		}
	} else if len(params) == 3 {
		timeout, err := strconv.Atoi(string(params[2]))
		if err != nil {
			log.Printf("parse_timeout_failed | msg=%v", err.Error())
			return nil, "invalid timeout", false
		}
		err = shardedMapStore.SetWithTimeout(key, params[1], time.Duration(timeout)*time.Second)
		if err != nil {
			log.Printf("set_with_timeout_cmd_failed | msg=%v", err.Error())
			return nil, err.Error(), false
		}
	}
	// TODO: handle other params


	return respOK, "", true
}

func handleDELCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 1 {
		return nil, "not enough parameters", false
	}
	key := string(params[0])
	// TODO: handle other params
	err := shardedMapStore.Delete(key)
	if err != nil {
		log.Printf("handler_del_cmd_failed | msg=%v", err.Error())
		return nil, err.Error(), false
	}
	return respOK, "", true
}

func handleINCRCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 1 {
		return nil, "not enough parameters", false
	}
	key := string(params[0])

	err := shardedMapStore.Increase(key)
	if err != nil {
		log.Printf("handler_increase_cmd_failed | msg=%v", err.Error())
		return nil, err.Error(), false      // TODO: The returned err msg should be unified
	}
	return respOK, "", true
}

func handleDUMPALLCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	// ignore param?
	// TODO: this method should limit the number of keys?
	jsonStr, err := shardedMapStore.DumpAllJSON()
	if err != nil {
		log.Printf("handler_dump_all_cmd_failed | msg=%v", err.Error())
		return nil, err.Error(), false
	}
	return []byte(jsonStr), "", true
}