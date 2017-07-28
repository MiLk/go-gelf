package packet

import (
	"compress/zlib"

	"io"

	"github.com/pkg/errors"
)

// What compression type the writer should use when sending messages
// to the graylog2 server
type CompressType int

const (
	CompressGzip CompressType = iota
	CompressZlib
	CompressNone
)

type Compressor struct {
	CompressionLevel int // one of the consts from compress/flate
	CompressionType  CompressType
}

func (c *Compressor) NewWriter(w io.Writer) (io.WriteCloser, error) {
	switch c.CompressionType {
	case CompressGzip:
		return zlib.NewWriterLevel(w, c.CompressionLevel)
	case CompressZlib:
		return zlib.NewWriterLevel(w, c.CompressionLevel)
	case CompressNone:
		return nil, nil
	default:
		return nil, errors.Errorf("unknown compression type %d", c.CompressionType)
	}
}
