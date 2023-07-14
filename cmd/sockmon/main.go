package main

import (
	"math/rand"
	"os"
	"time"

	// app "github.com/wide-vsix/sockmon/pkg/sockmon"
	app "sockmon/pkg/sockmon"

)

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := app.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}