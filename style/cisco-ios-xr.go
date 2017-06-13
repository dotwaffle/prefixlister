package style

import (
	"fmt"
	"net"
)

// CiscoIOSXR prefix list format
func CiscoIOSXR(prefixes []net.IPNet, name string) {
	// header, with mandatory name
	if name == "" {
		fmt.Printf("no prefix-set prefixlist\n")
		fmt.Printf("prefix-set prefixlist\n")
	} else {
		fmt.Printf("no prefix-set %s\n", name)
		fmt.Printf("prefix-set %s\n", name)
	}

	// construct a new list
	// the last entry does not have a comma
	for pos, prefix := range prefixes {
		if pos == len(prefixes)-1 {
			fmt.Printf("\t%s\n", prefix.String())
		} else {
			fmt.Printf("\t%s,\n", prefix.String())
		}
	}

	// footer
	fmt.Printf("end-set\n")
}
