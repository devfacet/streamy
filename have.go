// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy

import (
	"io"
	"os"
	"strings"
)

// HaveSeeker checks whether the given interface implements io.Seeker or not.
func HaveSeeker(i interface{}) bool {
	if _, ok := i.(io.Seeker); ok {
		if f, ok := i.(*os.File); ok && strings.HasPrefix(f.Name(), "/dev/std") {
			return false
		}
		return true
	}
	return false
}

// HaveStdin checks whether there is data in the stdin or not.
func HaveStdin() (bool, error) {
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
		return true, nil
	}
	return false, err
}
