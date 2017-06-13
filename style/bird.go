package style

import (
	"fmt"
	"net"
)

// BIRD prefix list format
func BIRD(prefixes []net.IPNet, name string) {
	// header, with mandatory name
	if name == "" {
		fmt.Printf("prefixlist = [\n")
	} else {
		fmt.Printf("%s = [\n", name)
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
	fmt.Printf("];\n")
}
