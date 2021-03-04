package tcp

import (
	"bufio"
	"log"
	"net"
	"testing"
)

func Test_SendTCPCmd(t *testing.T) {
	l ,err := testTCPServer("localhost", "3333")
	if err != nil {
		t.Errorf("test_tcp_server_err | err=%v", err.Error())
	}
	SendTCPCmd("localhost", "3333", "11223344")
	l.Close()
}

// testTCPServer run a tcp server and and only serve one request from one client
func testTCPServer(host string, port string) (l net.Listener, err error) {
	l, err = net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatalf("net_listen_error | err=%v", err.Error())
		return l, err
	}
	log.Printf("kash_server_listen_at | %v", host + ":" + port)

	go func() {
		//for {
			c, err := l.Accept()
			if err != nil {
				log.Fatalf("net_accept_error | err=%v", err.Error())
				return
			}

			//for {
				netData, err := bufio.NewReader(c).ReadBytes('\n')
				if err != nil {
					log.Fatalf("bufio_read_bytes_error | err=%v", err.Error())
					return
				}

				if string(netData) != "11223344\n" {
					//return errors.New(fmt.Sprintf("receiving_data_incorrect | received=%v | want=%v", string(netData), "11223344\n"))
					log.Fatalf("receiving_data_incorrect | received=%v | want=%v", string(netData), "11223344\n")
					return
				}

				// Always response "OK\n"
				_, err = c.Write([]byte("OK\n"))
				if err != nil {
					log.Fatalf("net_connection_write_error | err=%v", err.Error())
				}
			//}
			c.Close()
		//}
	}()
	return
}