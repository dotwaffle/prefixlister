package style

import (
	"fmt"
	"net"
)

// CiscoIOS is one of the most common prefix-list formats around
func CiscoIOS(prefixes []net.IPNet, name string) {
	// delete the old prefix-list
	fmt.Printf("no ip prefix-list %s\n", name)

	// construct a new list
	for _, prefix := range prefixes {
		fmt.Printf("ip prefix-list %s permit %s\n", name, prefix.String())
	}
}
