package style

import (
	"fmt"
	"net"
)

// CiscoIOS is one of the most common prefix-list formats around
func CiscoIOS(prefixes []net.IPNet, name string) {
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
	for _, prefix := range prefixes {
		if name == "" {
			fmt.Printf("%s prefix-list %s permit %s\n", afi, name, prefix.String())
		} else {
			fmt.Printf("%s prefix-list prefixlist permit %s\n", afi, prefix.String())
		}
	}
}
