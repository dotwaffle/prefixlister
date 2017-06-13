package main

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// ByPrefix (for sorting IPNets)
type ByPrefix []net.IPNet

// Len (sorter for IPNet)
func (s ByPrefix) Len() int {
	return len(s)
}

// Swap (sorter for IPNet)
func (s ByPrefix) Swap(left, right int) {
	s[left], s[right] = s[right], s[left]
}

// Less (sorter for IPNet)
func (s ByPrefix) Less(left, right int) bool {
	// compare the addresses a byte at a time
	for i := 0; i < len(s[left].IP); i++ {
		switch {
		case s[left].IP[i] < s[right].IP[i]:
			return true
		case s[left].IP[i] > s[right].IP[i]:
			return false
		case s[left].IP[i] == s[right].IP[i]:
			continue
		}
	}

	// compare the netmask if everything is equal so far!
	for i := 0; i < len(s[left].Mask); i++ {
		switch {
		case s[left].Mask[i] < s[right].Mask[i]:
			return true
		case s[left].Mask[i] > s[right].Mask[i]:
			return false
		case s[left].Mask[i] == s[right].Mask[i]:
			continue
		}
	}

	// they're identical...
	return true
}

func dedupePrefixes(s []string) []string {
	// prevent crazy expansion phase by providing appropriately lengthed slice to start with
	seen := make(map[string]struct{}, len(s))

	// location of the write head
	var pos int

	for _, v := range s {
		// have we already seen the prefix?
		if _, ok := seen[v]; ok {
			continue
		}

		// mark the prefix seen
		seen[v] = struct{}{}

		// store the new entry in the result slice
		s[pos] = v

		// move the write head along one
		pos++
	}

	return s[:pos]
}

func aggregatePrefixList(prefixes []net.IPNet) []net.IPNet {
	// if we don't make any changes on an aggregation run, consider it finished
	var changes = true
	for changes == true {
		// new run, no changes made yet
		changes = false

		for i := range prefixes {
			// does the next prefix and this prefix summarise to the same CIDR? If so, aggregate!
			if i < (len(prefixes) - 2) {
				shortened := shortenPrefixes(prefixes[i], prefixes[i+1])
				if len(shortened) == 1 {
					prefixes = append(shortened, prefixes[i+2:]...)
					changes = true
					break
				}
				merged := mergePrefixes(prefixes[i], prefixes[i+1])
				if len(merged) == 1 {
					prefixes = append(merged, prefixes[i+2:]...)
					changes = true
					break
				}
			}
		}

	}

	return prefixes
}

func shortenPrefixes(left, right net.IPNet) []net.IPNet {
	// firstly, does either left fit in right, or right fit in left?
	if left.Contains(right.IP) || right.Contains(left.IP) {
		// whichever has the shorter mask is by default the winner!
		leftMask, leftBits := left.Mask.Size()
		rightMask, rightBits := right.Mask.Size()
		if leftMask < rightMask && leftBits == rightBits {
			return []net.IPNet{left}
		} else if rightMask > leftMask {
			return []net.IPNet{right}
		} else {
			log.Fatal("Identical Mask Conundrum!")
		}
	}

	// they didn't contain each other, so can't be aggregated!
	return []net.IPNet{left, right}
}

func mergePrefixes(left, right net.IPNet) []net.IPNet {
	// extract the netmasks
	leftMask, leftBits := left.Mask.Size()
	rightMask, rightBits := right.Mask.Size()

	// are the masks are the same size? if so, they're a candidate for merging
	if leftMask == rightMask && leftBits == rightBits {
		// does one size shorter fit both prefixes?
		shorterNet := net.IPNet{IP: left.IP, Mask: net.CIDRMask((leftMask - 1), leftBits)}
		if shorterNet.Contains(right.IP) {
			prefixes := make([]net.IPNet, 0)
			return append(prefixes, shorterNet)
		}
	}

	// they didn't merge with each other, so can't be aggregated!
	return []net.IPNet{left, right}
}
