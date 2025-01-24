package testing

import "flag"

// Testing returns true if the code is running under go test
func Testing() bool {
	return flag.Lookup("test.v") != nil
}
