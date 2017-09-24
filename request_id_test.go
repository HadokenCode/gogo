package gogo

import (
	"math"
	"testing"
	"time"

	"github.com/golib/assert"
)

func Test_Macid_New(t *testing.T) {
	assertion := assert.New(t)
	id := NewMacid()

	// Generate 10 ids
	ids := make([]RequestID, 10)
	for i := 0; i < 10; i++ {
		ids[i] = id.New()
	}

	for i := 1; i < 10; i++ {
		id := ids[i]
		prevID := ids[i-1]

		// Test for uniqueness among all other 9 generated ids
		for j, tid := range ids {
			if j != i {
				assertion.NotEqual(id, tid)
			}
		}

		// Check that timestamp was incremented and is within 30 seconds of the previous one
		assertion.InDelta(prevID.Time().Second(), id.Time().Second(), 0.1)

		// Check that machine ids are the same
		assertion.Equal(prevID.Machine(), id.Machine())

		// Test for proper increment
		assertion.Equal(1, int(id.Counter()-prevID.Counter()))
	}
}

func Test_Macid_NewWithTime(t *testing.T) {
	assertion := assert.New(t)
	id := NewMacid()
	ts := time.Unix(12345678, 0)

	rid := id.NewWithTime(ts)
	assertion.Equal(ts, rid.Time())
	assertion.Equal([]byte{0x00, 0x00, 0x00, 0x00}, rid.Machine())
	assertion.EqualValues(0, rid.Counter())
}

func Test_IsRequestIDHex(t *testing.T) {
	assertion := assert.New(t)

	testCases := []struct {
		id    string
		valid bool
	}{
		{"59c741c4e6edc1faffffffff", true},
		{"59c741c4e6edc1fafffffff", false},
		{"59c741c4e6edc1faffffffff0", false},
		{"59c741c4e6edc1fafffffffx", false},
	}

	for _, testCase := range testCases {
		assertion.Equal(testCase.valid, IsRequestIDHex(testCase.id))
	}
}

func Test_RequestIDHex(t *testing.T) {
	assertion := assert.New(t)
	s := "59c741c4e6edc1faffffffff"

	id := RequestIDHex(s)
	assertion.True(id.Valid())
	assertion.Equal(s, id.Hex())
	assertion.EqualValues(1506230724, id.Time().Unix())
	assertion.Equal([]byte{0xe6, 0xed, 0xc1, 0xfa}, id.Machine())
	assertion.EqualValues(math.MaxUint32, id.Counter())
}
