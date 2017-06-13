package main

import (
	"bufio"
	"errors"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ERRORS
var (
	ErrKeyNotFound         = errors.New("Key Not Found")
	ErrMultipleCopies      = errors.New("Multiple Copies Found")
	ErrOtherError          = errors.New("Other whois error discovered")
	ErrBadWhoisLength      = errors.New("Bad WHOIS Length returned")
	ErrBadWhoisData        = errors.New("Bad WHOIS Data")
	ErrIncompleteWhoisData = errors.New("Incomplete WHOIS Data")
	ErrBadAFI              = errors.New("Bad AFI requested")
)

func expandASSet(whois *bufio.ReadWriter, set string) ([]string, error) {
	// send the query, recursively expanding as we go
	whois.WriteString("!i" + set + ",1\n")
	if err := whois.Flush(); err != nil {
		log.Fatal("Connection failure mid-stream")
	}

	return whoisResponseRead(whois)
}

func lookupASN(whois *bufio.ReadWriter, afi string, asn string) error {
	// work out the query string, and send it
	switch afi {
	case "4":
		whois.WriteString("!g" + asn + "\n")
	case "6":
		whois.WriteString("!6" + asn + "\n")
	default:
		return ErrBadAFI
	}

	return nil
}

func whoisResponseRead(whois *bufio.ReadWriter) ([]string, error) {
	// the first line should be a data length
	var msgLength int
	var data []string
	line, err := whois.ReadString('\n')
	switch {
	case strings.HasPrefix(line, "A"):
		data = append(data, line)
		msgLength, err = strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(data[0], "\n"), "A"))
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"data": data,
			}).Debug("Bad strconv?")
			return data, ErrBadWhoisLength
		}
	case strings.HasPrefix(line, "C"):
		return data, nil
	case strings.HasPrefix(line, "D"):
		return data, ErrKeyNotFound
	case strings.HasPrefix(line, "E"):
		return data, ErrMultipleCopies
	case strings.HasPrefix(line, "F"):
		return data, ErrOtherError
	default:
		log.WithFields(log.Fields{
			"line": line,
		}).Debug("Bad returned WHOIS Data")
		return data, ErrBadWhoisData
	}

	// read data back until we've read at least how much it said to expect
	var read int
	for read < msgLength {
		line, err := whois.ReadString('\n')
		if err != nil {
			log.WithFields(log.Fields{
				"line": line,
			}).Debug("Bad returned WHOIS Data")
			return nil, ErrBadWhoisData
		}
		read = read + len(line)
		data = append(data, line)
	}

	// read the confirmation code
	line, err = whois.ReadString('\n')
	data = append(data, line)

	// is the confirmation code good?
	if data[len(data)-1] != "C\n" {
		// did not receive a successful completion character
		return nil, ErrIncompleteWhoisData
	}

	// did we return enough data to actually parse?
	if len(data) < 3 {
		// did not return any useful data
		return nil, ErrBadWhoisData
	}

	// was the data we received the same length as we were told?
	reportedLength, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(data[0], "\n"), "A"))
	if err != nil {
		return nil, ErrBadWhoisData
	}

	// how much data did we actually get?
	if reportedLength != len(data[1]) {
		log.WithFields(log.Fields{
			"reported": reportedLength,
			"actual":   len(data),
		}).Debug("Data length mismatch")
		return nil, ErrBadWhoisLength
	}

	// remove the newline, split all the results out, and return them
	return strings.Split(strings.TrimSuffix(data[1], "\n"), " "), nil

}
