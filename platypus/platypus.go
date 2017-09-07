package main

import (
	"runtime"

	"github.com/cryptounicorns/platypus/cli"
)

func init() { runtime.GOMAXPROCS(runtime.NumCPU()) }
func main() { cli.Execute() }
