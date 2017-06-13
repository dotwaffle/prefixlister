package style

import (
	"fmt"
	"net"
)

// Quagga is one of the most common prefix-list formats around
func Quagga(prefixes []net.IPNet, name string) {
	// determine if IPv4 or IPv6 prefix list
	afi := "ip" // assume IPv4
	if ok := prefixes[0].IP.To4(); ok == nil {
		// can't convert to IPv4 address, must be IPv6
		afi = "ipv6"
	}

	// delete the old prefix-list
	if name == "" {
		fmt.Printf("no %s prefix-list prefixlist\n", afi)
	} else {
		fmt.Printf("no %s prefix-list %s\n", afi, name)
	}

	// construct a new list
	for pos, prefix := range prefixes {
		if name == "" {
			fmt.Printf("%s prefix-list %s seq %d permit %s\n", afi, name, pos*10, prefix.String())
		} else {
			fmt.Printf("%s prefix-list prefixlist seq %d permit %s\n", afi, pos*10, prefix.String())
		}
	}
}
