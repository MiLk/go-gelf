package packet

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
)

// Used to control GELF chunking.  Should be less than (MTU - len(UDP
// header)).
//
// TODO: generate dynamically using Path MTU Discovery?
const (
	ChunkSize        = 1420
	ChunkedHeaderLen = 12
	ChunkedDataLen   = ChunkSize - ChunkedHeaderLen
)

var (
	MagicChunked = []byte{0x1e, 0x0f}
)

// numChunks returns the number of GELF chunks necessary to transmit
// the given compressed buffer.
func numChunks(b []byte) int {
	lenB := len(b)
	if lenB <= ChunkSize {
		return 1
	}
	return len(b)/ChunkedDataLen + 1
}

// writes the gzip compressed byte array to the connection as a series
// of GELF chunked messages.  The format is documented at
// http://docs.graylog.org/en/2.1/pages/gelf.html as:
//
//     2-byte magic (0x1e 0x0f), 8 byte id, 1 byte sequence id, 1 byte
//     total, chunk-data
func writeChunked(w io.Writer, zBytes []byte) error {
	b := make([]byte, 0, ChunkSize)
	buf := bytes.NewBuffer(b)
	nChunksI := numChunks(zBytes)
	if nChunksI > 128 {
		return fmt.Errorf("msg too large, would need %d chunks", nChunksI)
	}
	nChunks := uint8(nChunksI)
	// use urandom to get a unique message id
	msgId := make([]byte, 8)
	n, err := io.ReadFull(rand.Reader, msgId)
	if err != nil || n != 8 {
		return fmt.Errorf("rand.Reader: %d/%s", n, err)
	}

	bytesLeft := len(zBytes)
	for i := uint8(0); i < nChunks; i++ {
		buf.Reset()
		// manually write header.  Don't care about
		// host/network byte order, because the spec only
		// deals in individual bytes.
		buf.Write(MagicChunked) //magic
		buf.Write(msgId)
		buf.WriteByte(i)
		buf.WriteByte(nChunks)
		// slice out our chunk from zBytes
		chunkLen := ChunkedDataLen
		if chunkLen > bytesLeft {
			chunkLen = bytesLeft
		}
		off := int(i) * ChunkedDataLen
		chunk := zBytes[off : off+chunkLen]
		buf.Write(chunk)

		// write this chunk, and make sure the write was good
		n, err := w.Write(buf.Bytes())
		if err != nil {
			return fmt.Errorf("Write (chunk %d/%d): %s", i,
				nChunks, err)
		}
		if n != len(buf.Bytes()) {
			return fmt.Errorf("Write len: (chunk %d/%d) (%d/%d)",
				i, nChunks, n, len(buf.Bytes()))
		}

		bytesLeft -= chunkLen
	}

	if bytesLeft != 0 {
		return fmt.Errorf("error: %d bytes left after sending", bytesLeft)
	}
	return nil
}
