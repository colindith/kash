package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// SendTCPCmd send the cmd to the remote tcp server and close the connection immediately
func SendTCPCmd(host string, port string, cmd string) string {
	// TODO: should maintain a TCP connection pool, instead of create new connection everytime
	// TODO: The conn pool can be configured with "max_active", "min_active", "active_timeout"
	conn, err := net.DialTimeout("tcp", host+ ":" + port, 10 * time.Second)
	if err != nil {
		log.Printf("dail_to_tcp_server_err | err=%v", err.Error())
		return "error..."
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, cmd + "\n")
	if err != nil {
		log.Printf("tcp_write_err | err=%v", err.Error())
		return "error..."
	}
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if resp != "OK\n" {
		log.Printf("NOT_OK | msg=%v", resp)
	}
	return resp
}
