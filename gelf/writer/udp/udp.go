package udp

import (
	"compress/flate"
	"net"
	"os"

	"github.com/Graylog2/go-gelf/gelf/codec/packet"
	"github.com/Graylog2/go-gelf/gelf/message"
	"github.com/Graylog2/go-gelf/gelf/writer"
)

type UDPWriter struct {
	writer.GelfWriter
	Pw packet.PacketWriter
}

// New returns a new GELF Writer.  This writer can be used to send the
// output of the standard Go log functions to a central GELF server by
// passing it to log.SetOutput()
func New(addr string) (*UDPWriter, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &UDPWriter{
		GelfWriter: writer.GelfWriter{
			Conn: conn,
			Host: host,
		},
		Pw: packet.PacketWriter{
			Compressor: packet.Compressor{
				CompressionLevel: flate.BestSpeed,
			},
		},
	}, nil
}

// WriteMessage sends the specified message to the GELF server
// specified in the call to New().  It assumes all the fields are
// filled out appropriately.  In general, clients will want to use
// Write, rather than WriteMessage.
func (w *UDPWriter) WriteMessage(m *message.Message) error {
	return w.Pw.WriteMessage(w.Conn, m)
}

// Write encodes the given string in a GELF message and sends it to
// the server specified in New().
func (w *UDPWriter) Write(p []byte) (n int, err error) {
	return w.WriteTo(w, p)
}
