package style

import (
	"fmt"
	"net"
)

// OpenBGPD prefix list format
func OpenBGPD(prefixes []net.IPNet, name string) {
	// header, with optional name
	if name == "" {
		fmt.Printf("prefix { \\\n")
	} else {
		fmt.Printf("%s=\"prefix { \\\n", name)
	}

	// construct a new list
	for _, prefix := range prefixes {
		fmt.Printf("\t%s \\\n", prefix.String())
	}

	// footer, with optional name
	if name == "" {
		fmt.Printf("\t}\n")
	} else {
		fmt.Printf("\t}\"\n")
	}
}
