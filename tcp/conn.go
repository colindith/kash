package tcp

import (
	"bufio"
	"context"
	"net"
	"sync"
	"time"
)

// DialOption specifies an option for dialing a Redis server.
type DialOption struct {
	f func(*dialOptions)
}

type dialOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	dialer       *net.Dialer
	dialFunc     func(ctx context.Context, network, addr string) (net.Conn, error)
}

// SetDialFunc set a customized dial function for creating TCP connection.
func SetDialFunc(dial func(network, addr string) (net.Conn, error)) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dial(network, addr)
		}
	}}
}

func setReadTimeout(readTimeout time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.readTimeout = readTimeout
	}}
}

func setWriteTimeout(writeTimeout time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.writeTimeout = writeTimeout
	}}
}


// DialContext connects to the Redis server at the given network and
// address using the specified options and context.
func DialContext(ctx context.Context, network, address string, options ...DialOption) (Conn, error) {
	do := dialOptions{
		dialer: &net.Dialer{
			Timeout:   time.Second * 30,
			KeepAlive: time.Minute * 5,
		},
	}
	for _, option := range options {
		option.f(&do)
	}
	if do.dialFunc == nil {
		do.dialFunc = do.dialer.DialContext
	}

	netConn, err := do.dialFunc(ctx, network, address)
	if err != nil {
		return nil, err
	}

	c := &conn{
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  do.readTimeout,
		writeTimeout: do.writeTimeout,
	}

	//if do.clientName != "" {
	//	// send client name to remote server
	//}


	return c, nil
}

// Conn represents a connection to a Kash server.
type Conn interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value when the connection is not usable.
	//Err() error

	// Do sends a command to the server and returns the received reply.
	Do(commandName string, args ...interface{}) (reply interface{}, err error)

	// Send writes the command to the client's output buffer.
	Send(commandName string, args ...interface{}) error

	// Flush flushes the output buffer to the Redis server.
	Flush() error

	// Receive receives a single reply from the Redis server
	Receive() (reply interface{}, err error)
}

var _ Conn = (*conn)(nil)

// conn implements Conn
type conn struct {
	// Shared
	mu      sync.Mutex
	conn    net.Conn

	// Read
	readTimeout time.Duration
	br          *bufio.Reader

	// Write
	writeTimeout time.Duration
	bw           *bufio.Writer
}

func (c *conn) Close() error {
	// TODO: implement this
	return nil
}

func (c *conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	// TODO: implement this
	return nil, nil
}

func (c *conn) Send(commandName string, args ...interface{}) error {
	// TODO: implement this
	return nil
}

func (c *conn) Flush() error {
	// TODO: implement this
	return nil
}

func (c *conn) Receive() (reply interface{}, err error) {
	// TODO: implement this
	return nil, nil
}