package style

import (
	"bytes"
	"fmt"
	"net"
)

// JUNOSPrefixList prefix list format
func JUNOSPrefixList(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// header, with mandatory name
	buf.WriteString(fmt.Sprintf("policy-options {\n"))
	buf.WriteString(fmt.Sprintf("replace:\n"))
	if name == "" {
		buf.WriteString(fmt.Sprintf("\tprefix-list prefixlist {\n"))
	} else {
		buf.WriteString(fmt.Sprintf("\tprefix-list %s {\n", name))
	}

	// construct a new list
	for _, prefix := range prefixes {
		buf.WriteString(fmt.Sprintf("\t\t%s;\n", prefix.String()))
	}

	// footer
	buf.WriteString(fmt.Sprintf("\t}\n}\n"))

	return buf
}

// JUNOSRouteFilter prefix list format
func JUNOSRouteFilter(prefixes []net.IPNet, name string) bytes.Buffer {
	var buf bytes.Buffer

	// header, with mandatory name
	buf.WriteString(fmt.Sprintf("policy-options {\n"))
	if name == "" {
		buf.WriteString(fmt.Sprintf("\tpolicy-statement prefixlist {\n"))
	} else {
		buf.WriteString(fmt.Sprintf("\tpolicy-statement %s {\n", name))
	}
	buf.WriteString(fmt.Sprintf("replace:\n"))
	buf.WriteString(fmt.Sprintf("\tfrom {\n"))

	// construct a new list
	for _, prefix := range prefixes {
		buf.WriteString(fmt.Sprintf("\t\troute-filter %s exact;\n", prefix.String()))
	}

	// footer
	buf.WriteString(fmt.Sprintf("\t\t}\n\t}\n}\n"))

	return buf
}
