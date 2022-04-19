// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy

import (
	"io"
)

// TeeReaderN returns an io.Reader that writes to the given writer what it reads from the given reader.
// It's similar to io.TeeReader but takes an additional argument that controls the number of bytes to write.
func TeeReaderN(r io.Reader, w io.Writer, n int64) io.Reader {
	return &teeReaderN{r: r, w: w, limit: n}
}

// teeReaderN represents a TeeReaderN entity.
type teeReaderN struct {
	r       io.Reader
	w       io.Writer
	limit   int64
	written int64
}

// Read implements the io.Reader interface.
func (t *teeReaderN) Read(p []byte) (n int, err error) {
	read, err := t.r.Read(p)
	if read > 0 && t.limit > t.written {
		high := int64(read)
		if (t.written + high) > t.limit {
			high = t.limit - t.written
		}
		written, err := t.w.Write(p[:high])
		if err != nil {
			return read, err
		}
		t.written += int64(written)
	}
	return read, err
}

// ReaderOnly takes any interface that implements io.Reader and returns just an io.Reader.
// It's useful for converting "advanced" readers (i.e. io.Seeker) to streams.
func ReaderOnly(r io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		if _, err := io.Copy(pw, r); err != nil {
			panic(err)
		}
	}()
	return pr
}
