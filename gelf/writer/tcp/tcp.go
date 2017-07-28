package tcp

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Graylog2/go-gelf/gelf/message"
	"github.com/Graylog2/go-gelf/gelf/writer"
)

const (
	DefaultMaxReconnect   = 3
	DefaultReconnectDelay = 1
)

type TCPWriter struct {
	writer.GelfWriter
	addr           string
	mu             sync.Mutex
	MaxReconnect   int
	ReconnectDelay time.Duration
}

func New(addr string) (*TCPWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &TCPWriter{
		GelfWriter: writer.GelfWriter{
			Conn: conn,
			Host: host,
		},
		MaxReconnect:   DefaultMaxReconnect,
		ReconnectDelay: DefaultReconnectDelay,
		addr:           addr,
	}, nil
}

// WriteMessage sends the specified message to the GELF server
// specified in the call to New().  It assumes all the fields are
// filled out appropriately.  In general, clients will want to use
// Write, rather than WriteMessage.
func (w *TCPWriter) WriteMessage(m *message.Message) (err error) {
	messageBytes, err := m.ToBytes()
	if err != nil {
		return err
	}

	messageBytes = append(messageBytes, 0)

	n, err := w.writeToSocketWithReconnectAttempts(messageBytes)
	if err != nil {
		return err
	}
	if n != len(messageBytes) {
		return fmt.Errorf("bad write (%d/%d)", n, len(messageBytes))
	}

	return nil
}

func (w *TCPWriter) Write(p []byte) (int, error) {
	return w.WriteTo(w, p)
}

func (w *TCPWriter) writeToSocketWithReconnectAttempts(zBytes []byte) (n int, err error) {
	var errConn error

	w.mu.Lock()
	for i := 0; n <= w.MaxReconnect; i++ {
		errConn = nil

		n, err = w.Conn.Write(zBytes)
		if err != nil {
			time.Sleep(w.ReconnectDelay * time.Second)
			w.Conn, errConn = net.Dial("tcp", w.addr)
		} else {
			break
		}
	}
	w.mu.Unlock()

	if errConn != nil {
		return 0, fmt.Errorf("Write Failed: %s\nReconnection failed: %s", err, errConn)
	}
	return n, nil
}
