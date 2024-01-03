package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sketch/pkg/router"
	"sketch/pkg/version"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	term = make(chan os.Signal, 1)
)

var (
	v   = flag.Bool("v", false, "show the binary build version")
	ver = flag.Bool("version", false, "show the binary build version")
)

func init() {
	// * Register Prometheus Metrics Collector
	prometheus.MustRegister(version.NewCollector(version.AppName))
}

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
	slog.Info("server start listen on", "addr", ":9000")

	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Printf("Received SIGTERM, exiting gracefully...")
	}
}

func showVersion() {
	if *v || *ver {
		fmt.Println(version.Print())
		os.Exit(0)
	}
}
