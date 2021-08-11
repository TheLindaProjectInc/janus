package main

import (
	"log"

	"github.com/TheLindaProjectInc/janus/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	cli.Run()
}
