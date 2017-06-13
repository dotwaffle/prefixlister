package style

import (
	"bytes"
	"fmt"
	"net"
)

// Brocade prefix list format
func Brocade(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// delete the old prefix-list
	if name == "" {
		buf.WriteString(fmt.Sprintf("no ip prefix-list prefixlist\n"))
	} else {
		buf.WriteString(fmt.Sprintf("no ip prefix-list %s\n", name))
	}

	// construct a new list
	for _, prefix := range prefixes {
		if name == "" {
			buf.WriteString(fmt.Sprintf("ip prefix-list %s permit %s\n", name, prefix.String()))
		} else {
			buf.WriteString(fmt.Sprintf("ip prefix-list prefixlist permit %s\n", prefix.String()))
		}
	}

	return buf
}
