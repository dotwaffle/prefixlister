package style

import (
	"fmt"
	"net"
)

// Brocade prefix list format
func Brocade(prefixes []net.IPNet, name string) {
	// delete the old prefix-list
	if name == "" {
		fmt.Printf("no ip prefix-list prefixlist\n")
	} else {
		fmt.Printf("no ip prefix-list %s\n", name)
	}

	// construct a new list
	for _, prefix := range prefixes {
		if name == "" {
			fmt.Printf("ip prefix-list %s permit %s\n", name, prefix.String())
		} else {
			fmt.Printf("ip prefix-list prefixlist permit %s\n", prefix.String())
		}
	}
}
