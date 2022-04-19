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

func TestHaveSeeker(t *testing.T) {
	f, _ := os.Open("have_test.go")

	table := []struct {
		arg0 interface{}
		want bool
	}{
		{bytes.Buffer{}, false},
		{bufio.NewWriter(io.Discard), false},
		{f, true},
	}
	for _, v := range table {
		if b := streamy.HaveSeeker(v.arg0); b != v.want {
			t.Errorf("got %v, want %v", b, v.want)
		}
	}
}

func TestHaveStdin(t *testing.T) {
	table := []struct {
		stdin bool
		want  bool
		err   error
	}{
		{false, false, nil},
	}
	for _, v := range table {
		if v.stdin {
			os.Stdin.WriteString("foo")
			os.Stdin.Sync()
		}
		b, err := streamy.HaveStdin()
		if b != v.want {
			t.Errorf("got %v, want %v", b, v.want)
		} else if err != v.err {
			t.Errorf("got %v, want %v", err, v.err)
		}
	}
}
