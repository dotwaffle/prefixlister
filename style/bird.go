package style

import (
	"bytes"
	"fmt"
	"net"
)

// BIRD prefix list format
func BIRD(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// header, with mandatory name
	if name == "" {
		buf.WriteString(fmt.Sprintf("prefixlist = [\n"))
	} else {
		buf.WriteString(fmt.Sprintf("%s = [\n", name))
	}

	// construct a new list
	// the last entry does not have a comma
	for pos, prefix := range prefixes {
		if pos == len(prefixes)-1 {
			buf.WriteString(fmt.Sprintf("\t%s\n", prefix.String()))
		} else {
			buf.WriteString(fmt.Sprintf("\t%s,\n", prefix.String()))
		}
	}

	// footer
	buf.WriteString(fmt.Sprintf("];\n"))

	return buf
}
