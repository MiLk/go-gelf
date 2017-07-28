package packet

import (
	"fmt"
	"io"

	"bytes"

	"github.com/Graylog2/go-gelf/gelf/message"
)

type PacketWriter struct {
	Compressor Compressor
}

type Option func(*PacketWriter)

func New() *PacketWriter {
	return &PacketWriter{}
}

func (pw *PacketWriter) WriteMessage(w io.Writer, m *message.Message) error {
	mBytes, err := m.Bytes()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	compWrt, err := pw.Compressor.NewWriter(&b)
	if err != nil {
		return err
	}
	if compWrt != nil {
		compWrt.Write(mBytes)
		compWrt.Close()
		mBytes = b.Bytes()
	}

	if numChunks(mBytes) > 1 {
		return writeChunked(w, mBytes)
	}
	n, err := w.Write(mBytes)
	if err != nil {
		return err
	}
	if n != len(mBytes) {
		return fmt.Errorf("bad write (%d/%d)", n, len(mBytes))
	}

	return nil
}
