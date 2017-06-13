package style

import (
	"bytes"
	"fmt"
	"net"
)

// Quagga is one of the most common prefix-list formats around
func Quagga(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// determine if IPv4 or IPv6 prefix list
	afi := "ip" // assume IPv4
	if ok := prefixes[0].IP.To4(); ok == nil {
		// can't convert to IPv4 address, must be IPv6
		afi = "ipv6"
	}

	// delete the old prefix-list
	if name == "" {
		buf.WriteString(fmt.Sprintf("no %s prefix-list prefixlist\n", afi))
	} else {
		buf.WriteString(fmt.Sprintf("no %s prefix-list %s\n", afi, name))
	}

	// construct a new list
	for pos, prefix := range prefixes {
		if name == "" {
			buf.WriteString(fmt.Sprintf("%s prefix-list %s seq %d permit %s\n", afi, name, pos*10, prefix.String()))
		} else {
			buf.WriteString(fmt.Sprintf("%s prefix-list prefixlist seq %d permit %s\n", afi, pos*10, prefix.String()))
		}
	}

	return buf
}
