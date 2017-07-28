// Copyright 2012 SocialCode. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gelf

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
	conn net.Conn
	host string
}

// Close connection and interrupt blocked Read or Write operations
func (w *GelfWriter) Close() error {
	return w.conn.Close()
}
