package main

import (
	"net"
	"strings"

	"github.com/dotwaffle/prefixlister/style"
)

func displayPrefixes(prefixes []net.IPNet, displayStyle string, displayName string) {
	// to prevent capitalisation issues, just lowercase everything
	switch strings.ToLower(displayStyle) {
	case "list":
		style.List(prefixes)
	case "cisco-ios", "ciscoios", "cisco", "ios":
		style.CiscoIOS(prefixes, displayName)
	default:
		style.List(prefixes)
	}
}
