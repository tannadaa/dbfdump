package main

import (
	"os/exec"
	"path/filepath"
)

func open(dir string) {
	cmd := exec.Command("explorer", "/select,", filepath.Join(dir, outFile))
	cmd.Run()
}
