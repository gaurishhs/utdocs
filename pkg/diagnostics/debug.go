package diagnostics

import "os"

func Debug(fn func()) {
	if os.Getenv("DEBUG") != "" {
		fn()
	}
}
