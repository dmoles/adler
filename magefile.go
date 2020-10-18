// +build mage

package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// ------------------------------------------------------------
// Constants and const-like variables

const projectName = "adler"

var statik = ensureCommand("statik", "statik not found; did you run go get github.com/rakyll/statik?")

// ------------------------------------------------------------
// Targets

//goland:noinspection GoUnusedExportedType
type Assets mg.Namespace

// Embeds static assets
func (Assets) Embed() error {
	cmd := exec.Command(statik, "-Z", "-src", "resources", "-ns", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println(strings.Join(cmd.Args, " "))

	return cmd.Run()
}

// ------------------------------------------------------------
// Helper functions

func ensureCommand(cmdName, failureMsg string) string {
	result, err := which(cmdName)
	if err != nil {
		os.Stderr.WriteString(failureMsg)
		os.Exit(1)
	}
	return result
}

func which(command string) (string, error) {
	var stdout bytes.Buffer

	cmd := exec.Command("which", command)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := stdout.String()
	return strings.TrimSpace(result), nil
}

