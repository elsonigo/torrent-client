package handshake

import (
	"fmt"
	"io"
)

// A Handshake is a special message that a peer uses to identify itself
type Handshake struct {
	Pstr string
	ReservedBytes
	InfoHash [20]byte
	PeerID   [20]byte
}

// The 8 ReservedBytes tell the client, which additional BitTorrent behaviour it supports/implements.
type ReservedBytes [8]byte

// supportsExtended checks if the 20th bit from the right (starting at 0)
// from the ReservedBytes is set to 1. This client supports the "extended" message (BEP10).
func (rb *ReservedBytes) supportsExtended() bool {
	return rb[5]&16 == 16
}

// New creates a new handshake with the standard pstr
func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

// Serialize serializes the handshake to a buffer
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], h.ReservedBytes[:])
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

// Read parses a handshake from a stream
func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lengthBuf[0])

	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	// 8 reserved, 20 infoHash, 20 peerID = 48
	handshakeBuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var reserved ReservedBytes
	copy(reserved[:], handshakeBuf[pstrlen:pstrlen+len(reserved)])

	var infoHash, peerID [20]byte
	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	h := Handshake{
		Pstr:          string(handshakeBuf[0:pstrlen]),
		ReservedBytes: reserved,
		InfoHash:      infoHash,
		PeerID:        peerID,
	}

	return &h, nil
}
