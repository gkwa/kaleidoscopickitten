package main

import (
	"os"

	"github.com/gkwa/kaleidoscopickitten/cmd"
	logging "gopkg.in/op/go-logging.v1"
)

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendLeveled)
}

func main() {
	cmd.Execute()
}
