package gelf

import "github.com/Graylog2/go-gelf/gelf/writer/tcp"

func NewTCPWriter(addr string) (*tcp.TCPWriter, error) {
	return tcp.New(addr)
}
