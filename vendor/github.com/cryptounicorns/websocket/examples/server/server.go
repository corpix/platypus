package main

import (
	"github.com/cryptounicorns/websocket/examples/server/cli"
)

func main() {
	var (
		err = cli.Run()
	)

	if err != nil {
		panic(err)
	}
}
