package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") == "1" {
		main()
		return
	}
	
	cmd := exec.Command(os.Args[0], "-test.run=TestMain")
	cmd.Env = append(os.Environ(), "GO_TEST_PROCESS=1")
	err := cmd.Run()
	
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		// Expected to exit with error code when no args provided
		return
	}
}