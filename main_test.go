package main

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"kaleidoscopickitten": main,
	})
}

var update = flag.Bool("u", false, "update testscript output files")

func TestScript(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir:                 filepath.Join("frontmatter", "testdata", "script"),
		UpdateScripts:       *update,
		RequireExplicitExec: true,
	})
}
