package config

import "runtime"

func OSName() (n string) {
	n = runtime.GOOS
	if n == "darwin" {
		n = "osx"
	}
	return
}
