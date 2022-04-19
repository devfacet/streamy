// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/devfacet/streamy"
)

func TestProgress(t *testing.T) {
	table := []struct {
		reader              io.Reader
		totalSize           int64
		statsEnabled        bool
		statsBytesWritten   int64
		statsBytesPerSecond int64
		statsTotalBytes     int64
		statsRemaining      time.Duration
		statsPercentage     int
	}{
		{
			reader:              bytes.NewBufferString("foo"),
			totalSize:           3,
			statsEnabled:        false,
			statsBytesWritten:   0,
			statsBytesPerSecond: 0,
			statsTotalBytes:     0,
			statsRemaining:      0,
			statsPercentage:     0,
		},
		{
			reader:              bytes.NewBufferString("bar"),
			totalSize:           3,
			statsEnabled:        true,
			statsBytesWritten:   3,
			statsBytesPerSecond: 3,
			statsTotalBytes:     3,
			statsRemaining:      0,
			statsPercentage:     100,
		},
		{
			reader:              &slowReader{content: "foo bar baz"},
			totalSize:           11,
			statsEnabled:        true,
			statsBytesWritten:   11,
			statsBytesPerSecond: 6,
			statsTotalBytes:     11,
			statsRemaining:      0,
			statsPercentage:     100,
		},
	}
	for _, v := range table {
		progress := streamy.Progress{}
		if v.statsEnabled {
			if err := progress.EnableStats(streamy.ProgressStatsModeSimple); err != nil {
				t.Errorf("got %v, want nil", err)
			}
		}
		if v.totalSize > 0 {
			progress.SetTotalSize(v.totalSize, streamy.Byte)
		}
		written, err := io.Copy(io.Discard, io.TeeReader(v.reader, &progress))
		if err != nil {
			t.Errorf("got %v, want nil", err)
		} else if written != v.totalSize {
			t.Errorf("got %v, want %v", written, v.totalSize)
		}
		stats := progress.Stats()
		if stats.BytesWritten != v.statsBytesWritten {
			t.Errorf("got %v, want %v", stats.BytesWritten, v.statsBytesWritten)
		} else if stats.BytesPerSecond != v.statsBytesPerSecond {
			t.Errorf("got %v, want %v", stats.BytesPerSecond, v.statsBytesPerSecond)
		} else if stats.TotalBytes != v.statsTotalBytes {
			t.Errorf("got %v, want %v", stats.TotalBytes, v.statsTotalBytes)
		} else if stats.Remaining != v.statsRemaining {
			t.Errorf("got %v, want %v", stats.Remaining, v.statsRemaining)
		} else if stats.Percentage != v.statsPercentage {
			t.Errorf("got %v, want %v", stats.Percentage, v.statsPercentage)
		}
	}
}

func BenchmarkProgressStatsDisabled(b *testing.B) {
	progress := streamy.Progress{}
	for i := 0; i < b.N; i++ {
		io.TeeReader(bytes.NewBufferString("foo"), &progress)
	}
}

func BenchmarkProgressStatsEnabled(b *testing.B) {
	progress := streamy.Progress{}
	if err := progress.EnableStats(streamy.ProgressStatsModeSimple); err != nil {
		b.Errorf("got %v, want nil", err)
	}
	for i := 0; i < b.N; i++ {
		io.TeeReader(bytes.NewBufferString("foo"), &progress)
	}
}

func BenchmarkProgressStats(b *testing.B) {
	progress := streamy.Progress{}
	if err := progress.EnableStats(streamy.ProgressStatsModeSimple); err != nil {
		b.Errorf("got %v, want nil", err)
	}
	for i := 0; i < b.N; i++ {
		io.TeeReader(bytes.NewBufferString("foo"), &progress)
		progress.Stats()
	}
}

// slowReader implements a slow reader.
type slowReader struct {
	content string
	written int64
}

// Read implements the io.Reader interface.
func (sr *slowReader) Read(p []byte) (n int, err error) {
	defer time.Sleep(200 * time.Millisecond) // slow it down

	p[0] = []byte(sr.content[sr.written : sr.written+1])[0]
	sr.written++
	if sr.written == int64(len(sr.content)) {
		return 1, io.EOF
	}
	return 1, nil
}
