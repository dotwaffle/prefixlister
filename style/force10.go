package style

import (
	"fmt"
	"net"
)

// Force10 prefix list format
func Force10(prefixes []net.IPNet, name string) {
	// delete the old prefix-list
	if name == "" {
		fmt.Printf("no ip prefix-list prefixlist\n")
		fmt.Printf("ip prefix-list prefixlist\n")
	} else {
		fmt.Printf("no ip prefix-list %s\n", name)
		fmt.Printf("ip prefix-list %s\n", name)
	}

	// construct a new list
	for pos, prefix := range prefixes {
		fmt.Printf("\tseq %d permit %s\n", pos*10, prefix.String())
	}
}
