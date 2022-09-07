/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package bytes

import (
	"encoding/binary"
)

// Bytes is just an array of bytes to allow cleaner code
type Bytes []byte

// New is a convenience function to create an array of bytes
func New() Bytes {
	return Bytes{}
}

// Copy returns a deep copy of the data
func (buf *Bytes) Copy() Bytes {
	ret := New()
	for _, b := range *buf {
		ret = append(ret, b)
	}
	return ret
}

// Get select bytes by location
func (buf *Bytes) Get(start int, length int) Bytes {

	// Check for nil pointer
	if buf == nil {
		return []byte{}
	}

	// Extract data we want
	tmp := New()
	if length < 1 {
		tmp = (*buf)[start:]
	} else {
		tmp = (*buf)[start : start+length]
	}

	// Return a deep copy
	return tmp.Copy()
}

// GetString gets select bytes by location and return as a string
func (buf *Bytes) GetString(start int, length int) string {

	// Check for nil pointer
	if buf == nil {
		return ""
	}

	// Return desired portion
	return string(buf.Get(start, length))
}

// GetUint32 gets selected bytes by location and return as an Uint32
func (buf *Bytes) GetUint32(start int, length int) uint32 {

	// Check for nil pointer or length too big for an uint32
	if buf == nil || length > 4 {
		return 0
	}

	b := buf.Get(start, length)
	return b.Uint32()
}

// Put supplied bytes in location
func (buf *Bytes) Put(start int, length int, b Bytes) {

	// Check for nil pointer
	if buf == nil {
		return
	}

	// Check for sufficient data
	if len(b) < length {
		return
	}

	for x := 0; x < length; x++ {
		(*buf)[start+x] = b[x]
	}

	return
}

// Append bytes
func (buf *Bytes) Append(slice Bytes) {

	// Check for nil pointer
	if buf == nil {
		return
	}

	// Append
	for _, b := range slice {
		*buf = append(*buf, b)
	}
}

// AppendString converts string to bytes and appends
func (buf *Bytes) AppendString(s string) {

	// Check for nil pointer
	if buf == nil {
		return
	}

	buf.Append([]byte(s))
}

// AppendUint32 converts uint32 to bytes and appends to buf
func (buf *Bytes) AppendUint32(i uint32) {

	// Check for nil pointer
	if buf == nil {
		return
	}

	// Convert uint32 to byte slice
	// BigEndian for network byte order
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)

	// append
	buf.Append(bytes)
}

// Uint32 converts bytes to uint32
func (buf *Bytes) Uint32() uint32 {

	// Check for nil pointer
	if buf == nil {
		return 0
	}

	// Maximum of four bytes
	if len(*buf) > 4 {
		return 0
	}

	// BigEndian from network byte order
	var ret uint32 = 0
	for _, b := range *buf {
		ret = (ret * 256) + uint32(b)
	}
	return ret
}

// Equal compares bytes and returns True if identical
func (buf *Bytes) Equal(m Bytes) bool {

	// Check for nil pointer
	if buf == nil {
		return false
	}

	// For exact match, lengths must be equal
	if len(*buf) != len(m) {
		return false
	}

	// Compare
	for i, b := range *buf {
		if b != m[i] {
			return false
		}
	}
	return true
}

// MatchStart compares to the start of bytes and returns true if they match
func (buf *Bytes) MatchStart(m Bytes) bool {

	// Check for nil pointer
	if buf == nil {
		return false
	}

	// if desired string is longer, it can not match
	if len(*buf) < len(m) {
		return false
	}

	// Iterate through *buf because pointer can't be indexed
	// End when difference found or end of m is reached
	for i, b := range *buf {
		if i >= len(m) {
			break
		}
		if b != m[i] {
			return false
		}
	}
	return true
}

// MatchStartString Compares a string to start of bytes and returns true if they match
func (buf *Bytes) MatchStartString(m string) bool {
	return buf.MatchStart([]byte(m))
}
