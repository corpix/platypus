package main

import (
	"runtime"

	"github.com/cryptounicorns/market-fetcher-http/cli"
)

func init() { runtime.GOMAXPROCS(runtime.NumCPU()) }
func main() { cli.Execute() }
