package gelf

import "io"

// All NewWhateverWriters must return a struct implementing io.WriteCloser.
// Use this to ensure you do.
//
func assertIsWriteCloser(_ io.WriteCloser) {}
