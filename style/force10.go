package style

import (
	"bytes"
	"fmt"
	"net"
)

// Force10 prefix list format
func Force10(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// delete the old prefix-list
	if name == "" {
		buf.WriteString(fmt.Sprintf("no ip prefix-list prefixlist\n"))
		buf.WriteString(fmt.Sprintf("ip prefix-list prefixlist\n"))
	} else {
		buf.WriteString(fmt.Sprintf("no ip prefix-list %s\n", name))
		buf.WriteString(fmt.Sprintf("ip prefix-list %s\n", name))
	}

	// construct a new list
	for pos, prefix := range prefixes {
		buf.WriteString(fmt.Sprintf("\tseq %d permit %s\n", pos*10, prefix.String()))
	}

	return buf
}
