// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

// Package streamy provides functions for streams.
package streamy

import (
	"bytes"
	"io"
)

// Index returns the index of the first instance of the given byte slice, number of bytes read and error if any.
func Index(r io.Reader, search []byte, readSize int) (index int64, read int64, err error) {
	if readSize == 0 {
		readSize = 4096
	}
	tailLen := len(search)
	b := make([]byte, readSize+tailLen)
	var i, n int
	//n, err = r.Read(b[tailLen:]) // err is checked below.
	for {
		read += int64(n)
		i = bytes.Index(b, search)
		if i > -1 {
			// Subtract the tail length since the bytes.Index method uses the entire byte slice.
			return index + int64(i-tailLen), read, nil
		} else if err != nil {
			if err == io.EOF {
				return -1, read, nil
			}
			return -1, read, err
		}
		copy(b, b[readSize:])        // Copy tail bytes to the beginning of the byte slice.
		index += int64(n)            // Update index before the read call.
		n, err = r.Read(b[tailLen:]) // err is checked above.
	}
}
