package main

import (
	"flag"
)

func init() {

	flag.BoolVar(&opt.List, "list", false, "List Google Cloud Run Jobs")

	flag.Usage = func() {
		showBanner()
		showUsage()
	}
	flag.Parse()
}

func main() {
	cfgPath := GetConfigPath()

	cfg, err := Load(cfgPath)
}
