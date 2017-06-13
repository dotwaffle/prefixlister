package style

import (
	"bytes"
	"fmt"
	"net"
)

// OpenBGPD prefix list format
func OpenBGPD(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// header, with optional name
	if name == "" {
		buf.WriteString(fmt.Sprintf("prefix { \\\n"))
	} else {
		buf.WriteString(fmt.Sprintf("%s=\"prefix { \\\n", name))
	}

	// construct a new list
	for _, prefix := range prefixes {
		buf.WriteString(fmt.Sprintf("\t%s \\\n", prefix.String()))
	}

	// footer, with optional name
	if name == "" {
		buf.WriteString(fmt.Sprintf("\t}\n"))
	} else {
		buf.WriteString(fmt.Sprintf("\t}\"\n"))
	}

	return buf
}
