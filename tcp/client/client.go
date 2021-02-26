package main

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/colindith/kash/cmd_reader"
)

// TODO: The client package should start a cmd line process and let users type the cmd in the cmd line.
// TODO: The cmd line should be provide some basic functions like using upper arrow key to find the history and and using left/right arrow to move the cursor

func main() {
	cmd_reader.Run("kash> ")
}

// SendTCPCmd send the cmd to the remote tcp server and close the connection immediately
func SendTCPCmd(host string, port string, cmd string) {
	// TODO: should maintain a TCP connection pool, instead of create new connection everytime
	// TODO: The conn pool can be configured with "max_active", "min_active", "active_timeout"
	conn, err := net.DialTimeout("tcp", host+ ":" + port, 10 * time.Second)
	if err != nil {
		fmt.Println("dail_to_tcp_server_err | err=", err.Error())
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, cmd + "\n")
	if err != nil {
		fmt.Println("tcp_write_err | err=", err.Error())
		return
	}
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if resp != "OK\n" {
		fmt.Println("NOT_OK | msg=", resp)
	}
}