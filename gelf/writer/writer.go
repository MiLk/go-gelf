package writer

import (
	"net"

	"github.com/Graylog2/go-gelf/gelf/message"
)

type Writer interface {
	Close() error
	Write([]byte) (int, error)
	WriteMessage(*message.Message) error
}

// Writer implements io.Writer and is used to send both discrete
// messages to a graylog2 server, or data from a stream-oriented
// interface (like the functions in log).
type GelfWriter struct {
	Conn net.Conn
	Host string
}

// Close connection and interrupt blocked Read or Write operations
func (w *GelfWriter) Close() error {
	return w.Conn.Close()
}

func (gw *GelfWriter) WriteTo(w Writer, p []byte) (n int, err error) {
	// 2 for the function that called our caller.
	file, line := getCallerIgnoringLogMulti(2)

	m := message.New(p, gw.Host, file, line)

	if err = w.WriteMessage(m); err != nil {
		return 0, err
	}

	return len(p), nil
}
