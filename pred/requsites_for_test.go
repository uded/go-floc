package pred

import "github.com/uded/go-floc"

func alwaysTrue(floc.Context) bool {
	return true
}

func alwaysFalse(floc.Context) bool {
	return false
}
