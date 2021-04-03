package main

import (
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "build", "-o", "out/bin/lol", ".")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
