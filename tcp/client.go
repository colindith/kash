package tcp

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

var pool sync.Map

func getConn(addr string) (*net.Conn, error) {
	if value, ok := pool.Load(addr); ok {
		if c, ok := value.(*net.Conn); ok {
			// TODO: how to know the conn is still alive?
			return c, nil
		}
	}
	conn, err := net.DialTimeout("tcp", addr, 10 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("err connect to host: %v", addr)
	}
	pool.Store(addr, &conn)
	return &conn, nil
}

func CloseConn() {
	pool.Range(func(key, value interface{}) bool {
		if c, ok := value.(*net.Conn); ok {
			(*c).Close()
		}
		return true
	})
}

// SendTCPCmd send the cmd to the remote tcp server and close the connection immediately
func SendTCPCmd(host string, port string, cmd string) string {
	// TODO: The conn pool can be configured with "max_active", "min_active", "active_timeout"
	conn, err := getConn(host+ ":" + port)
	if err != nil {
		//log.Printf("dail_to_tcp_server_err | err=%v", err.Error())
		return "error..."
	}

	_, err = fmt.Fprintf(*conn, cmd + "\n")
	if err != nil {
		//log.Printf("tcp_write_err | err=%v", err.Error())
		return "error..."
	}
	resp, err := bufio.NewReader(*conn).ReadString('\n')
	if err != nil {
		//log.Printf("tcp_read_err | err=%v", err.Error())
		return "error..."
	}
	//if resp != "OK\n" {
	//	log.Printf("NOT_OK | msg=%v", resp)
	//}
	return resp
}
