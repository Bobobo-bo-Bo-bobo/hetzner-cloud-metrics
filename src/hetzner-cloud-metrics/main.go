package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func usageHeader() {
	fmt.Printf(versionText, name, version, name, runtime.Version())
}

func usage() {
	usageHeader()
	fmt.Printf(helpText)
}

func main() {
	var tokenFile = flag.String("token-file", "", "File containing the API access token")
	var help = flag.Bool("help", false, "Show help")
	var version = flag.Bool("version", false, "Show version")
	var step = flag.Int64("step", 60, "step size")
	var cpu = flag.Bool("cpu", false, "Query CPU")
	var disk = flag.Bool("disk", false, "Query disk")
	var network = flag.Bool("network", false, "Quer network")
	var parsedServers HetznerAllServer
	var useProxy string
	var metrics []string

	flag.Usage = usage
	flag.Parse()
	trailing := flag.Args()

	if *help {
		usage()
		os.Exit(0)
	}

	if *version {
		usageHeader()
		os.Exit(0)
	}

	if *tokenFile == "" {
		fmt.Fprintf(os.Stderr, "Error: No token file provided\n\n")
		usage()
		os.Exit(1)
	}

	token, err := readSingleLine(*tokenFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *step <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Step must be greater than 0\n\n")
		os.Exit(1)
	}

	if !(*cpu || *disk || *network) {
		fmt.Fprintf(os.Stderr, "Error: No metric data selected\n\n")
		usage()
		os.Exit(1)
	}

	if *cpu {
		metrics = append(metrics, hetznerMetricsTypeCPU)
	}
	if *disk {
		metrics = append(metrics, hetznerMetricsTypeDisk)
	}
	if *network {
		metrics = append(metrics, hetznerMetricsTypeNetwork)
	}

	environment := getEnvironment()

	// Like curl the lower case variant has precedence.
	httpsproxy, lcfound := environment["https_proxy"]
	httpsProxy, ucfound := environment["HTTPS_PROXY"]

	if lcfound {
		useProxy = httpsproxy
	} else if ucfound {
		useProxy = httpsProxy
	}

	// get list of instances
	servers, err := httpRequest(hetznerAllServersURL, "GET", nil, nil, useProxy, token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if servers.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, servers.Status)
		os.Exit(1)
	}

	err = json.Unmarshal(servers.Content, &parsedServers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(trailing) > 0 {
		// remove all servers not in list
		parsedServers.Server = limitServerList(parsedServers.Server, trailing)
	}

	formatMetrics(parsedServers, *step, *cpu, *disk, *network, useProxy, token)
}
