/*
	Copyright (c) 2020 by Eric Jacksch VE3XEJ

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package bytes

import (
	"encoding/binary"
)

type Bytes []byte

// Convenience function
func New() Bytes {
	return Bytes{}
}

// Select bytes by location
func (buf *Bytes) Get(start int, length int) Bytes {

	// Check for nil pointer
	if buf == nil {
		return []byte{}
	}

	// Dereference
	tmp := *buf

	// Return desired portion
	return tmp[start : start+length]
}

// Select bytes by location and return as a string
func (buf *Bytes) GetString(start int, length int) string {

	// Check for nil pointer
	if buf == nil {
		return ""
	}

	// Return desired portion
	return string(buf.Get(start, length))
}

// Select bytes by location and return as a Uint32
func (buf *Bytes) GetUint32(start int, length int) uint32 {

	// Check for nil pointer or length too big for a uint32
	if buf == nil || length > 4 {
		return 0
	}

	b := buf.Get(start, length)
	return b.Uint32()
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

// Convert string to bytes and append
func (buf *Bytes) AppendString(s string) {

	// Check for nil pointer
	if buf == nil {
		return
	}

	buf.Append([]byte(s))
}

// Convert uint32 to bytes and append
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

// Convert bytes to uint32
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

// Compare bytes and determine if identical
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

// Compare to start of bytes
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

// Compare string to start of bytes
func (buf *Bytes) MatchStartString(m string) bool {
	return buf.MatchStart([]byte(m))
}
