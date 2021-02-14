package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func Test_runTCPServer(t *testing.T) {
	StartKashServer()

	conn, err := net.DialTimeout("tcp", connHost+ ":" +connPort, 10 * time.Second)
	if err != nil {
		// handle error
		t.Errorf("dail_to_tcp_server_err | err=%v", err.Error())
	}

	_, err = fmt.Fprintf(conn, "SET key1 123456\n")
	if err != nil {
		t.Errorf("tcp_write_err | err=%v", err.Error())
	}
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if resp != "OK\n" {
		t.Errorf("get_incorrect_set_resp | resp=%v, want=%v", resp, "OK\n")
	}

	_, err = fmt.Fprintf(conn, "GET key1\n")
	if err != nil {
		t.Errorf("tcp_write_err | err=%v", err.Error())
	}

	resp, err = bufio.NewReader(conn).ReadString('\n')
	if resp != "123456\n" {         // TODO: the \n is used as delimiter. This is should be handled by tcp client
		t.Errorf("get_incorrect_cached_data | data=%v, want=%v", resp, "123456")
	}
}