// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy_test

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/devfacet/streamy"
)

func TestTeeReaderN(t *testing.T) {
	table := []struct {
		arg0 io.Reader
		arg1 *bytes.Buffer
		arg2 int64
		out  string
	}{
		{bytes.NewBufferString("foo"), &bytes.Buffer{}, 3, "foo"},
		{bytes.NewBufferString("foo"), &bytes.Buffer{}, 2, "fo"},
		{bytes.NewBufferString("foo"), &bytes.Buffer{}, 1, "f"},
		{bytes.NewBufferString("foo"), &bytes.Buffer{}, 0, ""},
		{bytes.NewBufferString("foo"), &bytes.Buffer{}, -1, ""},
		{bytes.NewBufferString("foo"), &bytes.Buffer{}, 4, "foo"},
	}
	for _, v := range table {
		r := streamy.TeeReaderN(v.arg0, v.arg1, v.arg2)
		io.Copy(io.Discard, r)
		s := v.arg1.String()
		if s != v.out {
			t.Errorf("got %v, want %v", s, v.out)
		}
	}
}

func BenchmarkTeeReaderN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		streamy.TeeReaderN(bytes.NewBufferString("foo"), &bytes.Buffer{}, 2)
	}
}

func TestReaderOnly(t *testing.T) {
	f, _ := os.Open("have_test.go")

	table := []struct {
		arg0   io.Reader
		seeker bool
		writer bool
	}{
		{f, true, false},
		{bufio.NewReadWriter(&bufio.Reader{}, &bufio.Writer{}), false, true},
	}
	for _, v := range table {
		if v.seeker {
			if _, ok := v.arg0.(io.Seeker); !ok {
				t.Errorf("got no io.Seeker, want io.Seeker")
			} else if _, ok := streamy.ReaderOnly(v.arg0).(io.Seeker); ok {
				t.Errorf("got io.Seeker, want no io.Seeker")
			}
		}
		if v.writer {
			if _, ok := v.arg0.(io.Writer); !ok {
				t.Errorf("got no io.Writer, want io.Writer")
			} else if _, ok := streamy.ReaderOnly(v.arg0).(io.Writer); ok {
				t.Errorf("got io.Writer, want no io.Writer")
			}
		}
	}
}
