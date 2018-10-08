package main

import (
	"flag"
	"os"
)

var (
	appVersion = "dev"
	buildTime  = ""
)

type CmdOption struct {
	Host string
	Port int
}

func handleCmdArgs() *CmdOption {
	versionPtr := flag.Bool("v", false, "show version")
	o := &CmdOption{}
	flag.StringVar(&o.Host, "host", "0.0.0.0", "service listen host")
	flag.IntVar(&o.Port, "port", 2020, "service listen port")
	flag.Parse()
	if *versionPtr {
		//
		os.Exit(0)
	}
	return o
}
