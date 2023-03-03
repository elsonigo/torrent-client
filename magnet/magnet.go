package magnet

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

type Magnet struct {
	InfoHash    [20]byte
	DisplayName string
	Trackers    []string
	Params      url.Values
}

const (
	errCouldNotParseURI        = "could not parse uri"
	errUnexpectedURLScheme     = "unexpected url scheme"
	errMissingExactTopic       = "magnet link is missing the 'exact topic' parameter"
	errMissingDisplayNameParam = "magnet link is missing 'display name' parameter"
	errMissingAddressTrackers  = "magnet link is missing 'address trackers' parameter"
	errMalformedExactTopic     = "malformed 'exact topic'"
	errDecodeExactTopic        = "error decoding 'exact topic'"
	errUnsupportedExactTopic   = "unsupported encoding of 'exact topic'"
)

const exactTopicPrefix = "urn:btih:"

func ParseMagnet(uri string) (magnet Magnet, err error) {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		err = fmt.Errorf(errCouldNotParseURI)
		return
	}

	if parsedUrl.Scheme != "magnet" {
		err = fmt.Errorf(errUnexpectedURLScheme)
		return
	}

	query := parsedUrl.Query()

	magnet.InfoHash, err = decodeExactTopic(query)
	if err != nil {
		return
	}

	magnet.Trackers, err = extractAddressTrackers(query)
	if err != nil {
		return
	}

	magnet.DisplayName, err = extractDisplayName(query)
	if err != nil {
		return
	}

	magnet.Params = query

	return
}

func dropProcessedParametersFromQuery(query url.Values, keys ...string) {
	for _, key := range keys {
		query.Del(key)
	}
}

func extractDisplayName(query url.Values) (dn string, err error) {
	if !query.Has("dn") {
		err = fmt.Errorf(errMissingDisplayNameParam)
		return
	}

	dn = query.Get("dn")

	dropProcessedParametersFromQuery(query, "dn")

	return
}

// Extracts all the address trackers (&tr) parameters into a slice of strings
func extractAddressTrackers(query url.Values) (trackers []string, err error) {
	if !query.Has("tr") {
		err = fmt.Errorf(errMissingAddressTrackers)
		return
	}

	trackers = query["tr"]

	dropProcessedParametersFromQuery(query, "tr")

	return
}

// Extracts the 'exact topic' (&xt) parameter from the magnet link
// and decodes it to the [20]byte info hash of the torrent.
func decodeExactTopic(query url.Values) (infoHash [20]byte, err error) {
	if !query.Has("xt") {
		err = fmt.Errorf(errMissingExactTopic)
		return
	}

	xt := query.Get("xt")

	if !strings.HasPrefix(xt, exactTopicPrefix) {
		err = fmt.Errorf(errMalformedExactTopic)
		return
	}

	hashString := xt[len(exactTopicPrefix):]

	bytesWritten, err := decode(infoHash[:], []byte(hashString))
	if err != nil {
		err = fmt.Errorf(errDecodeExactTopic)
		return
	}

	if bytesWritten != 20 {
		err = fmt.Errorf(errDecodeExactTopic)
		return
	}

	dropProcessedParametersFromQuery(query, "xt")

	return
}

// For backwards compatibility with existing links, clients
// should also allow Base32 encoded versions of the hash.
// https://en.wikipedia.org/wiki/Magnet_URI_scheme
func decode(dst []byte, src []byte) (int, error) {
	switch len(src) {
	case 40:
		return hex.Decode(dst, src)
	case 32:
		return base32.StdEncoding.Decode(dst, src)
	}

	return 0, fmt.Errorf(errUnsupportedExactTopic)
}
