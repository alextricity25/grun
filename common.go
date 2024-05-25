package main

import (
	"log"
	"os"
	"path/filepath"
)

type Options struct {
	Edit    bool
	List    bool
	Remove  bool
	Version bool
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(homeDir, ".grun.json")
}
