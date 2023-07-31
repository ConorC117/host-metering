package main

import (
	"flag"
	"fmt"

	"redhat.com/milton/config"
	"redhat.com/milton/daemon"
	"redhat.com/milton/hostinfo"
)

func main() {
	prometheusUrl := flag.String("prometheus-url", "http://localhost:9090/api/v1/write", "URL of the Prometheus server")
	tick := flag.Uint("tick", 600, "Report every tick seconds")
	certPath := flag.String("cert", "", "Host certificate path")
	keyPath := flag.String("key", "", "Host certificate key path")

	flag.NewFlagSet("help", flag.ExitOnError)
	flag.NewFlagSet("daemon", flag.ExitOnError)
	flag.NewFlagSet("once", flag.ExitOnError)
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Error: no subcommand specified")
		printUsage()
		return
	}

	command := args[0]
	switch command {
	case "help":
		printUsage()
		return
	case "daemon", "once":
		cfg := config.NewConfig(*prometheusUrl, *tick, *certPath, *keyPath)
		cfg.Print()
		hostInfo, err := hostinfo.LoadHostInfo(cfg)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		hostInfo.Print()
		d := daemon.NewDaemon(cfg, hostInfo)

		if command == "once" {
			d.RunOnce()
			return
		}
		d.Run()
		return
	default:
		fmt.Println("Error: unknown subcommand", command)
		printUsage()
		return
	}
}

func printUsage() {
	fmt.Println("Usage: milton [OPTIONS] SUBCOMMAND")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println("Subcommands:")
	fmt.Println("  daemon    Run in daemon mode")
	fmt.Println("  once      Execute once")
	fmt.Println("  help      Print this help message")
}