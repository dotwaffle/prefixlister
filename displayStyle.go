package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dotwaffle/prefixlister/style"
)

func displayPrefixes(prefixes []net.IPNet, displayStyle string, displayName string) {
	// things move far quicker in buffers than with thousands of syscalls!
	var buf bytes.Buffer

	// to prevent capitalisation issues, just lowercase everything
	switch strings.ToLower(displayStyle) {
	case "list":
		buf = style.List(prefixes)
	case "join":
		buf = style.Join(prefixes)
	case "cisco-ios", "ciscoios", "cisco", "ios":
		buf = style.CiscoIOS(prefixes, displayName)
	case "cisco-xr", "cisco-ios-xr", "ios-xr", "xr":
		buf = style.CiscoIOSXR(prefixes, displayName)
	case "openbgpd":
		buf = style.OpenBGPD(prefixes, displayName)
	case "bird":
		buf = style.BIRD(prefixes, displayName)
	case "juniper", "junos", "juniper-prefix-list", "junos-prefix-list":
		buf = style.JUNOSPrefixList(prefixes, displayName)
	case "juniper-route-filter", "junos-route-filter":
		buf = style.JUNOSRouteFilter(prefixes, displayName)
	case "brocade":
		buf = style.Brocade(prefixes, displayName)
	case "force10":
		buf = style.Force10(prefixes, displayName)
	case "quagga":
		buf = style.Quagga(prefixes, displayName)
	case "redback":
		buf = style.Redback(prefixes, displayName)
	default:
		log.WithFields(log.Fields{
			"style": displayStyle,
			"name":  displayName,
		}).Debug("Display Style Not Found, using default (list) format")
		buf = style.List(prefixes)
	}

	// print the buffer to the screen
	fmt.Println(buf.String())
}
