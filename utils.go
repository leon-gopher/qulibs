package qulibs

import "os"

const (
	ciEnvKey = "PEDESTAL_CI"
)

func IsCI() bool {
	switch {
	case os.Getenv(ciEnvKey) == "yes":
		return true
	}

	return false
}
