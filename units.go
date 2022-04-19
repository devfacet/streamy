// Streamy
// For the full copyright and license information, please view the LICENSE.txt file.

package streamy

const (
	// Byte represent the byte unit type.
	Byte BinaryUnit = 1

	// KB represent the KB unit type.
	KB BinaryUnit = 1000
	// MB represent the MB unit type.
	MB BinaryUnit = 1000 * KB
	// GB represent the GB unit type.
	GB BinaryUnit = 1000 * MB
	// TB represent the TB unit type.
	TB BinaryUnit = 1000 * GB
	// PB represent the PB unit type.
	PB BinaryUnit = 1000 * TB

	// KiB represent the KiB unit type.
	KiB BinaryUnit = 1024
	// MiB represent the MiB unit type.
	MiB BinaryUnit = 1024 * KiB
	// GiB represent the GiB unit type.
	GiB BinaryUnit = 1024 * MiB
	// TiB represent the TiB unit type.
	TiB BinaryUnit = 1024 * GiB
	// PiB represent the PiB unit type.
	PiB BinaryUnit = 1024 * TiB
)

// BinaryUnit represents the binary unit type.
type BinaryUnit int64
