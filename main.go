package main

import (
	"log"
	"os"

	"github.com/gkwa/kaleidoscopickitten/cmd"
	logging "gopkg.in/op/go-logging.v1"
)

func init() {
	// Configure go-logging
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	format := logging.MustStringFormatter(`%{message}`)
	backendFormatted := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatted)
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendLeveled)

	// Configure standard log package (used by yq)
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}
