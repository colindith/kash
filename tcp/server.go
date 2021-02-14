package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
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
	respOK = []byte("OK")
	respNotOK = []byte("NOT OK")
)

var shardedMapStore store.Store

func main() {
	initStore()
	defer closeStore()

	runTCPServer()
}

func runTCPServer() {
	//arguments := os.Args
	//if len(arguments) == 1 {
	//	fmt.Println("Please provide a port number!")
	//	return
	//}



	l, err := net.Listen(connType, connHost + ":" + connPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

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

		args := bytes.Split(netData, []byte{' '})

		//args := strings.Split(strings.TrimSpace(netData), " ")
		switch string(args[0]) {
		case "STOP":
			break
		case "GET":
			var err error
			resp, ok := handleGETCmd(args[1:]...)
			if !ok {
				resp = respNotOK
			}
			_, err = c.Write(resp)

			if err != nil {
				// TODO: log the error
			}
		case "SET":
			var err error
			resp := respOK
			ok := handleSETCmd(args[1:]...)
			if !ok {
				resp = respNotOK
			}
			_, err = c.Write(resp)

			if err != nil {
				// TODO: log the error
			}
		// TODO: handle other cmd
		}


		result := "OK\n"
		c.Write([]byte(result))
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
