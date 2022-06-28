package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sketch/pkg/router"
	"sketch/pkg/version"
)

var (
	term = make(chan os.Signal, 1)
)

var (
	ver = flag.Bool("version", false, "show the binary build version")
)

//go:generate swag init -g sketch.go -o ../../pkg/swagger

// @title Sketch API
// @version 1.0
// @BasePath /v1
// @license.name Apache 2.0
// @description Sketch API Server
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	flag.Parse()
	showVersion()

	go router.Start()
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Printf("Received SIGTERM, exiting gracefully...")
	}
}

func showVersion() {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			fmt.Println(version.Print())
			os.Exit(0)
		}
	}
}
