package main

import (
	"fmt"
	"runtime"
)

var (
	version = "dev"
	commit  = "I don't remember exactly"
	date    = "I don't remember exactly"
)

// displayVersion DisplayVersion version.
func displayVersion(appName string) {
	fmt.Printf(`%s:
 version     : %s
 commit      : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, appName, version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
