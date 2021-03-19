package main

import (
	"bufio"
	"bytes"
	"fmt"
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
		"TTL":  handleTTLCmd,
	}
}

func handleGETCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 1 {
		return nil, "not enough parameters", false
	}
	key := string(params[0])
	// TODO: handle other params
	value, code := shardedMapStore.Get(key)
	if code != store.Success {
		log.Printf("handler_get_cmd_failed | code=%v", code)
		return nil, fmt.Sprintf("NOT OK: %v", code), false
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
		code := shardedMapStore.Set(key, params[1])
		if code != store.Success {
			log.Printf("set_cmd_failed | code=%v", code)
			return nil, fmt.Sprintf("NOT OK: %v", code), false
		}
	} else if len(params) == 3 {
		timeout, err := strconv.Atoi(string(params[2]))
		if err != nil {
			log.Printf("parse_timeout_failed | msg=%v", err.Error())
			return nil, "NOT OK: invalid timeout", false
		}
		code := shardedMapStore.SetWithTimeout(key, params[1], time.Duration(timeout)*time.Second)
		if code != store.Success {
			log.Printf("set_with_timeout_cmd_failed | code=%v", code)
			return nil, fmt.Sprintf("NOT OK: %v", code), false
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
	code := shardedMapStore.Delete(key)
	if code != store.Success {
		log.Printf("handler_del_cmd_failed | code=%v", code)
		return nil, fmt.Sprintf("NOT OK: %v", code), false
	}
	return respOK, "", true
}

func handleINCRCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 1 {
		return nil, "not enough parameters", false
	}
	key := string(params[0])

	code := shardedMapStore.Increase(key)
	if code != store.Success {
		log.Printf("handler_increase_cmd_failed | code=%v", code)
		return nil, fmt.Sprintf("NOT OK: %v", code), false      // TODO: The returned err msg should be unified
	}
	return respOK, "", true
}

func handleDUMPALLCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	// ignore param?
	// TODO: this method should limit the number of keys?
	jsonStr, code := shardedMapStore.DumpAllJSON()
	if code != code {
		log.Printf("handler_dump_all_cmd_failed | code=%v", code)
		return nil, fmt.Sprintf("NOT OK: %v", code), false
	}
	return []byte(jsonStr), "", true
}


func handleTTLCmd(params... []byte) (resp []byte, errMsg string, ok bool) {
	if len(params) < 1 {
		return nil, "not enough parameters", false
	}
	key := string(params[0])

	ttl, code := shardedMapStore.GetTTL(key)
	if code != store.Success {
		log.Printf("handler_get_ttl_cmd_failed | code=%v", code)
		return nil, fmt.Sprintf("NOT OK: %v", code), false      // TODO: The returned err msg should be unified
	}
	return []byte(strconv.Itoa(int(ttl - time.Now().UnixNano()))), "", true
}