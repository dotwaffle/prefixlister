package main

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dotwaffle/prefixlister/style"
)

func displayPrefixes(prefixes []net.IPNet, displayStyle string, displayName string) {
	// to prevent capitalisation issues, just lowercase everything
	switch strings.ToLower(displayStyle) {
	case "list":
		style.List(prefixes)
	case "cisco-ios", "ciscoios", "cisco", "ios":
		style.CiscoIOS(prefixes, displayName)
	case "cisco-xr", "cisco-ios-xr", "ios-xr", "xr":
		style.CiscoIOSXR(prefixes, displayName)
	case "openbgpd":
		style.OpenBGPD(prefixes, displayName)
	case "bird":
		style.BIRD(prefixes, displayName)
	case "juniper", "junos", "juniper-prefix-list", "junos-prefix-list":
		style.JUNOSPrefixList(prefixes, displayName)
	case "juniper-route-filter", "junos-route-filter":
		style.JUNOSRouteFilter(prefixes, displayName)
	case "brocade":
		style.Brocade(prefixes, displayName)
	case "force10":
		style.Force10(prefixes, displayName)
	case "quagga":
		style.Quagga(prefixes, displayName)
	case "redback":
		style.Redback(prefixes, displayName)
	default:
		log.WithFields(log.Fields{
			"style": displayStyle,
			"name":  displayName,
		}).Debug("Display Style Not Found, using default (list) format")
		style.List(prefixes)
	}
}
