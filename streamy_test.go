// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/devfacet/streamy"
)

func TestIndex(t *testing.T) {
	table := []struct {
		reader   io.Reader
		search   []byte
		readSize int
		index    int64
		read     int64
	}{
		{
			reader:   bytes.NewBufferString("tes"),
			search:   []byte("test"),
			readSize: 4,
			index:    -1,
			read:     3,
		},
		{
			reader:   bytes.NewBufferString("test"),
			search:   []byte("test"),
			readSize: 4,
			index:    0,
			read:     4,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 0,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 1,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 2,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 3,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 4,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 5,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 14,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("test"),
			readSize: 15,
			index:    10,
			read:     14,
		},
		{
			reader:   bytes.NewBufferString("this is a test"),
			search:   []byte("foo"),
			readSize: 0,
			index:    -1,
			read:     14,
		},
	}
	for _, v := range table {
		index, read, err := streamy.Index(v.reader, v.search, v.readSize)
		if err != nil {
			t.Errorf("got %v, want nil", err)
		} else if index != v.index {
			t.Errorf("got %v, want %v", index, v.index)
		} else if read != v.read {
			t.Errorf("got %v, want %v", read, v.read)
		}
	}
}

func BenchmarkIndexS2N1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		streamy.Index(bytes.NewBuffer([]byte{0x86, 0xc8, 0x63, 0xbf, 0xd2, 0x02, 0x96, 0x49, 0x49, 0x96}), []byte{0x49, 0x49}, 1)
	}
}

func BenchmarkIndexS2N2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		streamy.Index(bytes.NewBuffer([]byte{0x86, 0xc8, 0x63, 0xbf, 0xd2, 0x02, 0x96, 0x49, 0x49, 0x96}), []byte{0x49, 0x49}, 2)
	}
}

func BenchmarkIndexS2N10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		streamy.Index(bytes.NewBuffer([]byte{0x86, 0xc8, 0x63, 0xbf, 0xd2, 0x02, 0x96, 0x49, 0x49, 0x96}), []byte{0x49, 0x49}, 10)
	}
}

func BenchmarkIndexS2N0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		streamy.Index(bytes.NewBuffer([]byte{0x86, 0xc8, 0x63, 0xbf, 0xd2, 0x02, 0x96, 0x49, 0x49, 0x96}), []byte{0x49, 0x49}, 0)
	}
}
