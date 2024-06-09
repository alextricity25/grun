package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
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
	// cfgPath := GetConfigPath()

	// cfg, err := Load(cfgPath)
	p := tea.NewProgram(
		newModel(defaultTime),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
