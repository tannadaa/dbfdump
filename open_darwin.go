package main

import "os/exec"

func open(dir string) {
	cmd := exec.Command("open", dir)
	cmd.Run()
}
