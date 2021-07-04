package main

import (
	"flag"
	"fmt"
	"os"
)

var flagDebug = flag.Bool("debug", false, "")

func Debugf(msg string, args ...interface{}) {
	if !*flagDebug {
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] "+msg+"\n", args...)
}
