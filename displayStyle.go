package main

import (
	"net"

	"github.com/dotwaffle/prefixlister/style"
)

func displayPrefixes(prefixes []net.IPNet, displayStyle string) {
	switch displayStyle {
	case "list":
		style.List(prefixes)
	default:
		style.List(prefixes)
	}
}
