package main

import (
	"net"
	"os"
	"path/filepath"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type outputFormat struct {
	Name     string
	AFI      int
	Prefixes []string
}

func displayPrefixes(prefixes []net.IPNet, displayStyle string, displayName string) {
	// determine if this is an IPv4 or IPv6 prefix-list
	var afi int
	if ok := prefixes[0].IP.To4(); ok == nil {
		afi = 6
	} else {
		afi = 4
	}
	if !(afi == 4 || afi == 6) {
		log.Fatal("Invalid prefix list returned")
	}

	// format output suitable for display
	var data outputFormat
	data.Name = displayName
	data.AFI = afi
	data.Prefixes = make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		data.Prefixes = append(data.Prefixes, prefix.String())
	}

	// import the template
	tmpl, err := template.ParseFiles(filepath.Join("templates", displayStyle))
	if err != nil {
		log.Fatalf("IMPORT TEMPLATE FAIL: %s", err)
	}

	// apply the template, print it to stdout
	if err := tmpl.Execute(os.Stdout, data); err != nil {
		log.Fatalf("TEMPLATE EXECUTION FAIL: %s", err)
	}
}
