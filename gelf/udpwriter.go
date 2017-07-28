// Copyright 2012 SocialCode. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gelf

import "github.com/Graylog2/go-gelf/gelf/writer/udp"

// New returns a new GELF Writer.  This writer can be used to send the
// output of the standard Go log functions to a central GELF server by
// passing it to log.SetOutput()
func NewUDPWriter(addr string) (*udp.UDPWriter, error) {
	return udp.New(addr)
}
