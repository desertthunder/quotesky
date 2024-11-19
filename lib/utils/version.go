package utils

import (
	"fmt"
	"runtime/debug"
)

func Get() string {
	r := ""
	m := false

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				r = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					m = true
				}
			}
		}
	}

	if r == "" {
		return "unavailable"
	}

	if m {
		return fmt.Sprintf("%s-dirty", r)
	}

	return r
}
