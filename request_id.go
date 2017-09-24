package gogo

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

const (
	RequestIDBytes = 12
)

type Macid struct {
	mac     []byte
	counter uint32
}

func NewMacid() *Macid {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic("net.Interfaces(): " + err.Error())
	}

	mac := ""
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && bytes.Compare(iface.HardwareAddr, nil) != 0 {
			// Don't use random as we have a real address
			mac = iface.HardwareAddr.String()
			break
		}
	}
	if mac == "" {
		var sum [RequestIDBytes]byte

		id := sum[:]
		_, err := io.ReadFull(rand.Reader, id)
		if err != nil {
			panic(fmt.Errorf("cannot get random string: %v", err))
		}

		mac = string(id)
	}

	hash := md5.New()
	hash.Write([]byte(mac))

	return &Macid{
		mac:     hash.Sum(nil),
		counter: 0,
	}
}

// New returns a new unique RequestID.
func (id *Macid) New() RequestID {
	var b [RequestIDBytes]byte

	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))

	// MAC, first 4 bytes of md5(MAC)
	b[4] = id.mac[0]
	b[5] = id.mac[1]
	b[6] = id.mac[2]
	b[7] = id.mac[3]

	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&id.counter, 1)
	b[8] = byte(i >> 24)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	return RequestID(b[:])
}

// NewWithTime returns a dummy RequestID with the timestamp part filled
// with the provided number of seconds from epoch UTC, and all other parts
// filled with zeroes. It's not safe to insert a document with an id generated
// by this method, it is useful only for queries to find documents with ids
// generated before or after the specified timestamp.
func (id *Macid) NewWithTime(t time.Time) RequestID {
	var b [RequestIDBytes]byte

	binary.BigEndian.PutUint32(b[:4], uint32(t.Unix()))

	return RequestID(string(b[:]))
}

// RequestID is a unique ID identifying a unique value. It must be exactly 12 bytes
// long.
//
type RequestID string

// IsRequestIDHex returns whether s is a valid hex representation of
// an RequestID. See the RequestIDHex function.
func IsRequestIDHex(s string) bool {
	if len(s) != RequestIDBytes*2 {
		return false
	}

	_, err := hex.DecodeString(s)
	return err == nil
}

// RequestIDHex returns an RequestID from the provided hex representation.
// Calling this function with an invalid hex representation will
// cause a runtime panic. See the IsRequestIDHex function.
func RequestIDHex(s string) RequestID {
	d, err := hex.DecodeString(s)
	if err != nil || len(d) != RequestIDBytes {
		panic(fmt.Sprintf("Invalid input to RequestIDHex: %q", s))
	}

	return RequestID(d)
}

// Hex returns a hex representation of the RequestID.
func (id RequestID) Hex() string {
	return hex.EncodeToString([]byte(id))
}

// Valid returns true if id is valid. A valid id must contain exactly 12 bytes.
func (id RequestID) Valid() bool {
	return len(id) == RequestIDBytes
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id RequestID) Time() time.Time {
	// First 4 bytes of RequestID is 32-bit big-endian seconds from epoch.
	secs := int64(binary.BigEndian.Uint32(id.byteSlice(0, 4)))

	return time.Unix(secs, 0)
}

// Machine returns the 4-byte machine mac part of the md5 mac address.
// It's a runtime error to call this method with an invalid id.
func (id RequestID) Machine() []byte {
	return id.byteSlice(4, 8)
}

// Counter returns the incrementing value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id RequestID) Counter() uint32 {
	b := id.byteSlice(8, 12)

	// Counter is stored as big-endian 3-bytes value
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// String returns a hex string representation of the id.
// Example: RequestIDHex("4d88e15b60f486e428412dc9").
func (id RequestID) String() string {
	return fmt.Sprintf(`RequestIDHex("%x")`, string(id))
}

// byteSlice returns byte slice of id from start to end.
// Calling this function with an invalid id will cause a runtime panic.
func (id RequestID) byteSlice(start, end int) []byte {
	if len(id) != RequestIDBytes {
		panic(fmt.Sprintf("Invalid RequestID: %q", string(id)))
	}

	return []byte(string(id)[start:end])
}
