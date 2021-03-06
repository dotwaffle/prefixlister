package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	progName    = "prefixlister"
	progVersion = "0.1"
)

// regexps
var (
	reAS      = regexp.MustCompile("^[Aa][Ss]([1-9][0-9]+|[1-9])$")
	reASN     = regexp.MustCompile("^([1-9][0-9]+|[1-9])$")
	reASSet   = regexp.MustCompile("^[AaRr][Ss].+$")
	reSources = regexp.MustCompile("^[A-Za-z0-9,]+$")
)

// flags
var (
	debug         = flag.Bool("debug", false, "Enable debugging mode")
	whoisServer   = flag.String("host", "rr.ntt.net", "WHOIS server to query, irrd servers only (not RIPE)")
	whoisPort     = flag.String("port", "43", "WHOIS port to query")
	afi           = flag.String("afi", "4", "Address Family to query [4|6]")
	aggregate     = flag.Bool("aggregate", false, "Aggregate prefixes [BROKEN, SLOW, UNFINISHED, imagine implicit orlonger]")
	pipelineDepth = flag.Int("pipeline", -1, "Pipeline Depth")
	speedMode     = flag.Bool("speed-mode", false, "Activate speed mode [NOSORTING, NODEDUPE, NOAGGREGATE, THERESNOLIMIT]")
	displayStyle  = flag.String("style", "list", "Style of prefix-list to generate")
	displayName   = flag.String("name", "prefixlister", "Name of prefix-list to generate")
	sources       = flag.String("sources", "", "Sources (default: all, recommended: RADB,RIPE,APNIC)")
)

func main() {
	// parse options and validate
	flag.Parse()
	if *debug == true {
		log.SetLevel(log.DebugLevel)
	}
	if *pipelineDepth > 16384 || *pipelineDepth < -1 || *pipelineDepth == 0 {
		// some of these values might actually work, but let's prevent sillyness!
		log.WithFields(log.Fields{
			"pipeline": *pipelineDepth,
		}).Fatal("Too long of a pipeline")
	}
	if *afi != "4" && *afi != "6" {
		log.WithFields(log.Fields{
			"afi": *afi,
		}).Fatal("Only IPv4 and IPv6 supported")
	}
	if *sources != "" && !reSources.MatchString(*sources) {
		// prevent nefarious usage of sources line
		flag.Usage()
		log.WithFields(log.Fields{
			"sources": *sources,
		}).Fatal("Bad Sources line: letters, numbers, and commas only")
	}

	// get query
	remainingArgs := flag.Args()
	if len(remainingArgs) != 1 {
		flag.Usage()
		log.WithFields(log.Fields{
			"query": remainingArgs,
		}).Fatal("Bad query arguments: Must have only one query")
	}
	query := remainingArgs[0]

	// is query valid?
	var queryList []string
	var expand bool
	if reAS.MatchString(query) {
		queryList = append(queryList, query)
	} else if reASN.MatchString(query) {
		queryList = append(queryList, "AS"+query)
	} else if reASSet.MatchString(query) {
		expand = true
	} else {
		log.WithFields(log.Fields{
			"query": query,
		}).Fatal("Failed to understand input query")
	}

	// dial whois server
	timeConnOpen := time.Now()
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

	// if we need to set record sources, do so now
	if *sources != "" {
		whois.WriteString("!s" + *sources + "\n")
		if err := whois.Flush(); err != nil {
			log.Fatal("Connection failure mid-stream")
		}
		confirmation, err := whois.ReadString('\n')
		log.WithFields(log.Fields{
			"sources":      *sources,
			"confirmation": strings.TrimSuffix(confirmation, "\n"),
		}).Debug("Set Sources")
		if err != nil || confirmation != "C\n" {
			log.Fatal("Failed to set record sources")
		}
	}

	// if we need to expand the query, do so now
	if expand == true {
		whois.WriteString("!a\n")
		if err := whois.Flush(); err != nil {
			log.Fatal("Connection failure mid-stream")
		}
		confirmation, err := whois.ReadString('\n')
		if confirmation == "F Missing required set name for A query\n" {
			log.Debug("Server supports 'a' queries")
			queryList = append(queryList, query)
		} else {
			log.WithFields(log.Fields{
				"confirmation": confirmation,
				"query":        query,
			}).Debug("Expanding Query Set")
			queryList, err = expandASSet(whois, query)
			if err != nil {
				log.WithFields(log.Fields{
					"query": query,
					"err":   err,
				}).Fatal("Failed to get AS-SET result")
			}
			expand = false
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

	// run inverse lookup for each object requested, asynchronously
	go func() {
		log.WithFields(log.Fields{
			"afi":     *afi,
			"queries": len(queryList),
		}).Debug("Querying objects")
		for _, asn := range queryList {
			<-wait
			if err := lookupRecordKey(whois, *afi, asn, expand); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Fatal("Bad WHOIS AFI Lookup")
			}
			if err := whois.Flush(); err != nil {
				log.WithFields(log.Fields{
					"afi":       *afi,
					"asn":       asn,
					"available": whois.Available(),
					"expand":    expand,
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
	timeConnClose := time.Now()
	log.WithFields(log.Fields{
		"prefixes": len(results),
		"duration": timeConnClose.Sub(timeConnOpen),
	}).Debug("Connection Closed")

	// if we have no results, assume this isn't wanted and fail early
	if len(results) == 0 {
		log.WithFields(log.Fields{
			"query": query,
		}).Fatal("No prefixes returned")
	}

	// FIXME: temporary hack
	// eventually, speedMode will be removed =(
	if *speedMode {
		// print results out to stdout
		fmt.Println(strings.Join(results, "\n"))
	} else {
		// deduplication
		timeDedupeStart := time.Now()
		results = dedupePrefixes(results)
		timeDedupeFinish := time.Now()
		log.WithFields(log.Fields{
			"duration": timeDedupeFinish.Sub(timeDedupeStart),
			"prefixes": len(results),
		}).Debug("After deduplication")

		// conversion from strings to structs
		prefixes := make([]net.IPNet, 0, len(results))
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
		timeConvertFinish := time.Now()
		log.WithFields(log.Fields{
			"duration": timeConvertFinish.Sub(timeDedupeFinish),
			"prefixes": len(prefixes),
		}).Debug("After conversion")

		// sorting into order
		sort.Sort(ByPrefix(prefixes))
		timeSortingFinish := time.Now()
		log.WithFields(log.Fields{
			"duration": timeSortingFinish.Sub(timeConvertFinish),
			"prefixes": len(prefixes),
		}).Debug("After sorting")

		if *aggregate {
			prefixes = aggregatePrefixList(prefixes)
			timeAggregateFinish := time.Now()
			log.WithFields(log.Fields{
				"duration": timeAggregateFinish.Sub(timeSortingFinish),
				"prefixes": len(prefixes),
			}).Debug("After aggregation")
		}

		// print results out to stdout
		timeRenderStart := time.Now()
		displayPrefixes(prefixes, *displayStyle, *displayName)
		timeRenderFinish := time.Now()
		log.WithFields(log.Fields{
			"duration": timeRenderFinish.Sub(timeRenderStart),
			"prefixes": len(prefixes),
		}).Debug("Time to render")
	}

}
