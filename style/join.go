package style

import (
	"bytes"
	"net"
	"strings"
)

// Join is the faster version of Join, just to demonstrate it
// You may want to use "List" as your base, as it's more applicable to real life!
func Join(prefixes []net.IPNet) bytes.Buffer {
	var buf bytes.Buffer

	// convert all the prefixes to strings
	prefixesConverted := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		prefixesConverted = append(prefixesConverted, prefix.String())
	}

	// join all the prefixes together, and write it out to the buffer
	buf.WriteString(strings.Join(prefixesConverted, "\n"))

	return buf
}
