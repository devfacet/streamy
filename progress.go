// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy

import (
	"errors"
	"math"
	"sync"
	"time"
)

// ErrProgressStopped means that a progress stopped (i.e. by calling Stop method).
var ErrProgressStopped = errors.New("stopped")

// Progress implements the io.Writer interface for tracking bytes.
type Progress struct {
	bytesWritten     int64
	totalBytes       int64
	statsEnabled     bool
	statsMode        ProgressStatsMode
	statsBytes       [][]int64
	statsBytesLimit  int
	statsFrom        int64
	statsTo          int64
	controlsEnabled  bool
	controlsClosedCh <-chan struct{}
	controlsStopCh   chan struct{}
	rwMu             sync.RWMutex
	stopped          bool
}

// Write implements the io.Writer interface.
func (progress *Progress) Write(p []byte) (n int, err error) {
	// Init vars
	n = len(p)
	ni := int64(n)

	progress.rwMu.RLock()
	defer progress.rwMu.RUnlock()

	// Update written bytes
	progress.bytesWritten += ni

	// Update stats
	if progress.statsEnabled {
		unixNano := time.Now().UnixNano()
		if progress.statsFrom == 0 {
			// Note that initial start time (from writer) might be earlier.
			progress.statsFrom = unixNano
		}

		// Shift bytes
		if len(progress.statsBytes) >= progress.statsBytesLimit {
			progress.statsBytes = append(progress.statsBytes[1:progress.statsBytesLimit], []int64{unixNano, ni})
		} else {
			progress.statsBytes = append(progress.statsBytes, []int64{unixNano, ni})
		}
	}

	// Check controls
	if progress.controlsEnabled {
		select {
		case <-progress.controlsClosedCh:
		case _, ok := <-progress.controlsStopCh:
			if !ok {
				progress.stopped = true
				return 0, ErrProgressStopped
			}
		}
	}

	// Update stats
	if progress.statsEnabled {
		// Any delay in this block should be added to the stats.
		progress.statsTo = time.Now().UnixNano()
	}

	return n, nil
}

// BytesWritten returns the number of bytes written.
func (progress *Progress) BytesWritten() int64 {
	return progress.bytesWritten
}

// SetTotalSize sets the total size by the given size and binary unit.
func (progress *Progress) SetTotalSize(size int64, unit BinaryUnit) {
	progress.totalBytes = size * int64(unit)
}

// TotalBytes returns the total number of bytes.
func (progress *Progress) TotalBytes() int64 {
	return progress.totalBytes
}

// EnableStats enables the progress stats.
func (progress *Progress) EnableStats(mode ProgressStatsMode) error {
	progress.rwMu.Lock()
	defer progress.rwMu.Unlock()
	if !progress.statsEnabled {
		switch mode {
		case ProgressStatsModeSimple:
			progress.statsMode = ProgressStatsModeSimple
			progress.statsBytesLimit = 10
		default:
			return errors.New("invalid mode")
		}
	}
	progress.statsEnabled = true
	return nil
}

// DisableStats disables the progress stats.
func (progress *Progress) DisableStats() {
	progress.rwMu.Lock()
	defer progress.rwMu.Unlock()
	if progress.statsEnabled {
		progress.statsBytes = nil
		progress.statsBytesLimit = 0
		progress.statsFrom = 0
		progress.statsTo = 0
	}
	progress.statsEnabled = false
}

// EnableControls enables the progress controls such as pause and stop.
func (progress *Progress) EnableControls() error {
	progress.rwMu.Lock()
	defer progress.rwMu.Unlock()
	if !progress.controlsEnabled {
		// Ref: https://go.dev/ref/spec#Receive_operator
		//			https://go.dev/ref/spec#Send_statements
		//			https://go.dev/ref/spec#Close
		ch := make(chan struct{})
		close(ch)
		progress.controlsClosedCh = ch
		progress.controlsStopCh = make(chan struct{})
	}
	progress.controlsEnabled = true
	return nil
}

// Stop stops the the progress writer.
func (progress *Progress) Stop() {
	progress.rwMu.Lock()
	defer progress.rwMu.Unlock()
	if !progress.controlsEnabled || progress.stopped {
		return
	}
	close(progress.controlsStopCh)
}

// Stats returns the progress stats.
func (progress *Progress) Stats() ProgressStats {
	if !progress.statsEnabled {
		// Do not return anything to avoid confusion
		return ProgressStats{}
	}

	// Init stats
	stats := ProgressStats{
		BytesWritten: progress.bytesWritten,
		TotalBytes:   progress.totalBytes,
		Took:         time.Duration(progress.statsTo - progress.statsFrom),
	}

	// Calculate the percentage
	// Note that totalBytes is given by the user and it can be 0 (see progress.SetTotalSize).
	if progress.totalBytes > 0 {
		if progress.totalBytes == progress.bytesWritten {
			stats.Percentage = 100
		} else {
			stats.Percentage = int((float64(progress.bytesWritten) / float64(progress.totalBytes)) * 100)
		}
	}

	// Calculate the bytes per second
	switch progress.statsMode {
	case ProgressStatsModeSimple:
		var first int64 = 0
		var last int64 = 0
		var total int64 = 0
		for i, l := 0, len(progress.statsBytes); i < l; i++ {
			if i == 0 {
				first = progress.statsBytes[i][0]
			}
			if i+1 == l { // The length might be 1 so don't use else.
				last = progress.statsBytes[i][0]
			}
			total += progress.statsBytes[i][1]
		}
		diff := last - first
		if stats.Took.Seconds() > 1 {
			stats.BytesPerSecond = int64(math.Round(float64(total) / time.Duration(diff).Seconds()))
		} else {
			stats.BytesPerSecond = stats.BytesWritten
		}
	}

	// Calculate remaining seconds
	// Note that totalBytes is given by the user and it can be 0 (see progress.SetTotalSize).
	if progress.totalBytes > 0 {
		if progress.totalBytes == progress.bytesWritten {
			stats.Remaining = 0
		} else if stats.BytesPerSecond > 0 {
			stats.Remaining = time.Duration(math.Ceil(float64(progress.totalBytes-progress.bytesWritten)/float64(stats.BytesPerSecond))) * time.Second
		}
	}

	return stats
}

// ProgressStats represents the progress stats.
type ProgressStats struct {
	TotalBytes     int64
	BytesWritten   int64
	BytesPerSecond int64
	Took           time.Duration
	Remaining      time.Duration
	Percentage     int
}

// ProgressStatsMode represents a progress stats mode.
type ProgressStatsMode struct {
	mode uint8
}

var (
	// ProgressStatsModeSimple represents the simple progress stats mode.
	ProgressStatsModeSimple = ProgressStatsMode{mode: 0}
)
