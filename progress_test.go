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
		reader    io.Reader
		totalSize int64
	}{
		{
			reader:    bytes.NewBufferString("foo"),
			totalSize: 3,
		},
		{
			reader:    bytes.NewBufferString("bar"),
			totalSize: 3,
		},
		{
			reader: bytes.NewBufferString(""),
		},
	}
	for _, v := range table {
		progress := streamy.Progress{}
		if v.totalSize > 0 {
			progress.SetTotalSize(v.totalSize, streamy.Byte)
		}
		written, err := io.Copy(io.Discard, io.TeeReader(v.reader, &progress))
		if err != nil {
			t.Errorf("got %v, want nil", err)
		} else if written != v.totalSize {
			t.Errorf("got %v, want %v", written, v.totalSize)
		}
	}
}

func TestProgressStats(t *testing.T) {
	delay := 10 * time.Millisecond
	table := []struct {
		reader              io.Reader
		totalSize           int64
		statsBytesPerSecond int64
		statsBytesWritten   int64
		statsPercentage     int
		statsRemaining      time.Duration
		statsTotalBytes     int64
		goEnableStats       bool
	}{
		{
			reader:              bytes.NewBufferString("foo"),
			statsBytesPerSecond: 3,
			statsBytesWritten:   3,
			statsPercentage:     100,
			statsRemaining:      0,
		},
		{
			reader:              bytes.NewBufferString("bar"),
			statsBytesPerSecond: 3,
			statsBytesWritten:   3,
			statsPercentage:     100,
			statsRemaining:      0,
			statsTotalBytes:     3,
			totalSize:           3,
		},
		{
			reader:              &slowReader{content: "foo bar baz", delay: 200 * time.Millisecond},
			totalSize:           11,
			statsBytesPerSecond: 6,
			statsBytesWritten:   11,
			statsPercentage:     100,
			statsRemaining:      0,
			statsTotalBytes:     11,
		},
		{
			reader:              &slowReader{content: "foo bar baz", delay: 10 * time.Millisecond},
			totalSize:           11,
			statsBytesPerSecond: 11,
			statsBytesWritten:   11,
			statsPercentage:     100,
			statsRemaining:      0,
			statsTotalBytes:     11,
			goEnableStats:       true,
		},
	}
	for _, v := range table {
		progress := streamy.Progress{}
		if v.totalSize > 0 {
			progress.SetTotalSize(v.totalSize, streamy.Byte)
		}
		if v.goEnableStats {
			go func() {
				time.Sleep(delay * 2)
				if err := progress.EnableStats(streamy.ProgressStatsModeSimple); err != nil {
					t.Errorf("got %v, want nil", err)
				}
			}()
		} else {
			if err := progress.EnableStats(streamy.ProgressStatsModeSimple); err != nil {
				t.Errorf("got %v, want nil", err)
			}
		}
		written, err := io.Copy(io.Discard, io.TeeReader(v.reader, &progress))
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
		if v.totalSize > 0 && written != v.totalSize {
			t.Errorf("got %v, want %v", written, v.totalSize)
		}
		stats := progress.Stats()
		if stats.BytesWritten != v.statsBytesWritten {
			t.Errorf("got %v, want %v", stats.BytesWritten, v.statsBytesWritten)
		} else if stats.BytesPerSecond != v.statsBytesPerSecond {
			t.Errorf("got %v, want %v", stats.BytesPerSecond, v.statsBytesPerSecond)
		} else if stats.Remaining != v.statsRemaining {
			t.Errorf("got %v, want %v", stats.Remaining, v.statsRemaining)
		}
		if v.totalSize > 0 {
			if v.totalSize > 0 && stats.TotalBytes != v.statsTotalBytes {
				t.Errorf("got %v, want %v", stats.TotalBytes, v.statsTotalBytes)
			} else if stats.Percentage != v.statsPercentage {
				t.Errorf("got %v, want %v", stats.Percentage, v.statsPercentage)
			}
		}
	}
}

func TestProgressControls(t *testing.T) {
	delay := 10 * time.Millisecond
	table := []struct {
		reader           io.Reader
		totalSize        int64
		goEnableControls bool
		stopCall         bool
		dupStopCall      bool
		enableStats      bool
	}{
		{
			reader:    &slowReader{content: "foo bar baz", delay: delay},
			totalSize: 11,
		},
		{
			reader:      &slowReader{content: "foo bar baz", delay: delay},
			totalSize:   11,
			enableStats: true,
		},
		{
			reader:           &slowReader{content: "foo bar baz", delay: delay},
			totalSize:        11,
			goEnableControls: true,
		},
		{
			reader:           &slowReader{content: "foo bar baz", delay: delay},
			totalSize:        11,
			goEnableControls: true,
			enableStats:      true,
		},
		{
			reader:      &slowReader{content: "foo bar baz", delay: delay},
			totalSize:   11,
			dupStopCall: true,
		},
		{
			reader:           &slowReader{content: "foo bar baz", delay: delay},
			totalSize:        11,
			goEnableControls: true,
			dupStopCall:      true,
		},
	}
	for _, v := range table {
		progress := streamy.Progress{}
		if v.enableStats {
			if err := progress.EnableStats(streamy.ProgressStatsModeSimple); err != nil {
				t.Errorf("got %v, want nil", err)
			}
		}
		if v.goEnableControls {
			go func() {
				time.Sleep(delay * 2)
				if err := progress.EnableControls(); err != nil {
					t.Errorf("got %v, want nil", err)
				}
			}()
		} else {
			if err := progress.EnableControls(); err != nil {
				t.Errorf("got %v, want nil", err)
			}
		}
		if v.stopCall {
			go func() {
				time.Sleep(delay * 4)
				progress.Stop()
			}()
			if v.dupStopCall {
				go func() {
					time.Sleep(delay * 4)
					progress.Stop()
				}()
			}
		}
		written, err := io.Copy(io.Discard, io.TeeReader(v.reader, &progress))
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
		if v.stopCall {
			if written <= 0 || written >= v.totalSize {
				t.Errorf("got %v, want >0 <%v", written, v.totalSize)
			}
		} else {
			if written != v.totalSize {
				t.Errorf("got %v, want %v", written, v.totalSize)
			}
		}
	}
}

func BenchmarkProgress(b *testing.B) {
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

func BenchmarkControlsEnabled(b *testing.B) {
	progress := streamy.Progress{}
	if err := progress.EnableControls(); err != nil {
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
	delay   time.Duration
	written int64
}

// Read implements the io.Reader interface.
func (sr *slowReader) Read(p []byte) (n int, err error) {
	defer time.Sleep(sr.delay)

	p[0] = []byte(sr.content[sr.written : sr.written+1])[0]
	sr.written++
	if sr.written == int64(len(sr.content)) {
		return 1, io.EOF
	}
	return 1, nil
}
