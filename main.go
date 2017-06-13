package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	progName    = "prefixlister"
	progVersion = "0.1"
)

// regexps
var (
	reAS    = regexp.MustCompile("AS([1-9][0-9]+|[1-9])")
	reASN   = regexp.MustCompile("([1-9][0-9]+|[1-9])")
	reASSet = regexp.MustCompile("(AS-|RS-).+")
)

// flags
var (
	debug         = flag.Bool("debug", false, "Enable debugging mode")
	whoisServer   = flag.String("host", "whois.radb.net", "WHOIS server to query")
	whoisPort     = flag.String("port", "43", "WHOIS port to query")
	afi           = flag.String("afi", "4", "Address Family to query [4|6]")
	aggregate     = flag.Bool("aggregate", false, "Aggregate prefixes [BROKEN, SLOW, UNFINISHED, imagine implicit orlonger]")
	pipelineDepth = flag.Int("pipeline", -1, "Pipeline Depth")
	speedMode     = flag.Bool("speed-mode", false, "Activate speed mode [NOSORTING, NODEDUPE, NOAGGREGATE, THERESNOLIMIT]")
)

func main() {
	// parse options and validate
	flag.Parse()
	if *debug == true {
		log.SetLevel(log.DebugLevel)
	}
	query := flag.Arg(0)

	// is query valid?
	var queryList []string
	var expand bool
	if reASN.MatchString(query) {
		queryList = append(queryList, "AS"+query)
	} else if reAS.MatchString(query) {
		queryList = append(queryList, query)
	} else if reASSet.MatchString(query) {
		expand = true
	} else {
		log.WithFields(log.Fields{
			"query": query,
		}).Fatal("Failed to understand input query")
	}

	// dial whois server
	whoisConn, err := net.Dial("tcp", net.JoinHostPort(*whoisServer, *whoisPort))
	if err != nil {
		log.WithFields(log.Fields{
			"host": *whoisServer,
			"port": *whoisPort,
		}).Fatal("Failed to connect to whois server")
	}
	defer whoisConn.Close()
	whois := bufio.NewReadWriter(bufio.NewReader(whoisConn), bufio.NewWriter(whoisConn))
	log.WithFields(log.Fields{
		"host": *whoisServer,
		"port": *whoisPort,
	}).Debug("Connected")

	// keep whois connection connection open for multiple queries
	whois.WriteString("!!\n")
	if err := whois.Flush(); err != nil {
		log.Fatal("Connection failure mid-stream")
	}

	// identify ourselves to the whois server
	whois.WriteString("!n" + progName + "-" + progVersion + "\n")
	if err := whois.Flush(); err != nil {
		log.Fatal("Connection failure mid-stream")
	}
	confirmation, err := whois.ReadString('\n')
	log.WithFields(log.Fields{
		"identity":     progName + "-" + progVersion,
		"confirmation": strings.TrimSuffix(confirmation, "\n"),
	}).Debug("Set Identity")
	if err != nil || confirmation != "C\n" {
		log.Fatal("Failed to set tool name for statistics/logging purposes")
	}

	// if we need to expand the query, do so now
	if expand == true {
		log.WithFields(log.Fields{
			"query": query,
		}).Debug("Expanding")
		queryList, err = expandASSet(whois, query)
		if err != nil {
			log.WithFields(log.Fields{
				"query": query,
				"err":   err,
			}).Fatal("Failed to get AS-SET result")
		}
	}

	// add the initial queries to the pipeline
	// if we set -1, just do all of them at once
	var chanSize int
	if *pipelineDepth == -1 {
		chanSize = len(queryList)
	} else {
		chanSize = *pipelineDepth
	}
	wait := make(chan bool, chanSize)
	for i := 0; i < chanSize; i++ {
		wait <- true
	}

	// run inverse lookup for each ASN requested, asynchronously
	go func() {
		for _, asn := range queryList {
			<-wait
			log.WithFields(log.Fields{
				"afi": *afi,
				"asn": asn,
			}).Debug("Querying ASN")
			err := lookupASN(whois, *afi, asn)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Fatal("Bad WHOIS AFI Lookup")
			}
			if err := whois.Flush(); err != nil {
				log.WithFields(log.Fields{
					"afi":       *afi,
					"asn":       asn,
					"available": whois.Available(),
					"trigger":   "writer",
				}).Debug("Write Buffer Statistics")
			}
		}
	}()

	// read the results off the wire
	var results []string
	for i := 0; i < len(queryList); i++ {
		// release another spot in the pipeline
		wait <- true

		result, err := whoisResponseRead(whois)
		if err == ErrKeyNotFound {
			continue
		} else if err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"result": result,
			}).Fatal("Bad WHOIS ASN Lookup")
		}
		results = append(results, result...)

	}

	// close whois connection
	whois.WriteString("!q\n")

	if *speedMode {
		// print results out to stdout
		fmt.Println(strings.Join(results, "\n"))
	} else {
		// dedupe, then sort nicely
		log.WithFields(log.Fields{
			"prefixes": len(results),
		}).Debug("Before deduplication")
		results = dedupePrefixes(results)

		var prefixes []net.IPNet
		for _, result := range results {
			_, prefix, err := net.ParseCIDR(result)
			if err != nil {
				log.WithFields(log.Fields{
					"cidr": result,
					"err":  err,
				}).Fatal("Bad CIDR returned from WHOIS")
			}
			prefixes = append(prefixes, *prefix)
		}
		sort.Sort(ByPrefix(prefixes))
		log.WithFields(log.Fields{
			"prefixes": len(prefixes),
		}).Debug("After deduplication")

		if *aggregate {
			prefixes = aggregatePrefixList(prefixes)
			log.WithFields(log.Fields{
				"prefixes": len(prefixes),
			}).Debug("After aggregation")
		}

		// print results out to stdout
		for _, prefix := range prefixes {
			fmt.Println(prefix.String())
		}
	}

}
