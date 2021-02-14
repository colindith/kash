package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	respNotOK = []byte("NOT OK\n")
)

var shardedMapStore store.Store

func main() {
	initStore()
	defer closeStore()

	go runTCPServer()
}

func runTCPServer() {
	l, err := net.Listen(connType, connHost + ":" + connPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	fmt.Println("listen at: " + connHost + ":" + connPort)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
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
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()
	for {
		//netData, err := bufio.NewReader(c).ReadString('\n')
		netData, err := bufio.NewReader(c).ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		args := bytes.Split(netData[:len(netData)-1], []byte{' '})

		result := respOK
		switch string(args[0]) {
		case "STOP":
			break
		case "GET":
			var err error
			resp, ok := handleGETCmd(args[1:]...)
			if !ok {
				result = respNotOK
			} else {
				result = append(resp, byte('\n'))
			}

			if err != nil {
				// TODO: log the error
			}
		case "SET":
			var err error
			ok := handleSETCmd(args[1:]...)
			if !ok {
				result = respNotOK
			}

			if err != nil {
				// TODO: log the error
			}
		// TODO: handle other cmd
		}
		_, err = c.Write(result)
		if err != nil {
			// TODO: log the error
		}
	}

}

func handleGETCmd(params... []byte) (data []byte, ok bool) {
	key := string(params[0])
	// TODO: handle other params
	value, err := shardedMapStore.Get(key)
	if err != nil {
		// TODO: log the error
		return nil, false
	}
	data = value.([]byte)
	return data, true
}

func handleSETCmd(params... []byte) (ok bool) {
	if len(params) < 2 {
		return false
	}
	key := string(params[0])
	if len(params) == 2 {
		err := shardedMapStore.Set(key, params[1])
		if err != nil {
			// TODO: log the error
			return false
		}
	} else if len(params) == 3 {
		timeout, err := strconv.Atoi(string(params[2]))
		if err != nil {
			// TODO: log the error
			// TODO: should also return error msg
			return false
		}
		err = shardedMapStore.SetWithTimeout(key, params[1], time.Duration(timeout)*time.Second)
		if err != nil {
			// TODO: log the error
			return false
		}
	}
	// TODO: handle other params


	return true
}
