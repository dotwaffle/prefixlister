package style

import (
	"fmt"
	"net"
)

// JUNOSPrefixList prefix list format
func JUNOSPrefixList(prefixes []net.IPNet, name string) {
	// header, with mandatory name
	fmt.Printf("policy-options {\n")
	fmt.Printf("replace:\n")
	if name == "" {
		fmt.Printf("\tprefix-list prefixlist {\n")
	} else {
		fmt.Printf("\tprefix-list %s {\n", name)
	}

	// construct a new list
	for _, prefix := range prefixes {
		fmt.Printf("\t\t%s;\n", prefix.String())
	}

	// footer
	fmt.Printf("\t}\n}\n")
}

// JUNOSRouteFilter prefix list format
func JUNOSRouteFilter(prefixes []net.IPNet, name string) {
	// header, with mandatory name
	fmt.Printf("policy-options {\n")
	if name == "" {
		fmt.Printf("\tprefix-statement prefixlist {\n")
	} else {
		fmt.Printf("\tprefix-statement %s {\n", name)
	}
	fmt.Printf("replace:\n")
	fmt.Printf("\tfrom {\n")

	// construct a new list
	for _, prefix := range prefixes {
		fmt.Printf("\t\troute-filter %s exact;\n", prefix.String())
	}

	// footer
	fmt.Printf("\t\t}\n\t}\n}\n")
}
