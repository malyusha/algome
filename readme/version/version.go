package version

import (
	"fmt"
)

var (
	name    = "algome"
	version = "dev"
	date    = "n/a"
)

// String returns a string representation of the build.
func String() string {
	return fmt.Sprintf("%s %s [date: %s]", name, version, date)
}
