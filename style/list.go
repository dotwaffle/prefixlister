package style

import (
	"bytes"
	"fmt"
	"net"
)

// List is the basic style type, provided as an example
func List(prefixes []net.IPNet) bytes.Buffer {
	var buf bytes.Buffer

	// if you wanted to put a header, you'd put it here

	// iterate over the list
	for _, prefix := range prefixes {
		// net.IPNet contains two fields:
		// * IP (type net.IP, []byte, can be v4 or v6)
		// * Mask (type IPMask, []byte, can be v4 or v6)
		// the String() method does "the right thing" for the IP version
		buf.WriteString(fmt.Sprintf("%s\n", prefix.String()))
	}

	// if you wanted to put a footer, you'd put it here

	return buf
}
