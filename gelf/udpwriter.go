// Copyright 2012 SocialCode. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gelf

import (
	"compress/flate"
	"net"
	"os"

	"github.com/Graylog2/go-gelf/gelf/codec/packet"
	"github.com/Graylog2/go-gelf/gelf/message"
)

type UDPWriter struct {
	GelfWriter
	Pw *packet.PacketWriter
}

// New returns a new GELF Writer.  This writer can be used to send the
// output of the standard Go log functions to a central GELF server by
// passing it to log.SetOutput()
func NewUDPWriter(addr string) (*UDPWriter, error) {
	var err error
	w := new(UDPWriter)
	w.Pw = packet.New()
	w.Pw.Compressor.CompressionLevel = flate.BestSpeed

	if w.conn, err = net.Dial("udp", addr); err != nil {
		return nil, err
	}
	if w.host, err = os.Hostname(); err != nil {
		return nil, err
	}

	return w, nil
}

// WriteMessage sends the specified message to the GELF server
// specified in the call to New().  It assumes all the fields are
// filled out appropriately.  In general, clients will want to use
// Write, rather than WriteMessage.
func (w *UDPWriter) WriteMessage(m *message.Message) error {
	return w.Pw.WriteMessage(w.conn, m)
}

// Write encodes the given string in a GELF message and sends it to
// the server specified in New().
func (w *UDPWriter) Write(p []byte) (n int, err error) {
	// 1 for the function that called us.
	file, line := getCallerIgnoringLogMulti(1)

	m := message.New(p, w.host, file, line)

	if err = w.WriteMessage(m); err != nil {
		return 0, err
	}

	return len(p), nil
}
