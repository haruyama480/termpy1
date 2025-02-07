package main

import (
	"flag"
	"time"

	"github.com/haruyama480/termpy1/game"
)

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "seed for tsumo")

	g := game.TokoConsole{}
	g.Run(seed)
}
