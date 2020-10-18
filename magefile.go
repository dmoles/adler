// +build mage

package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

const projectName = "adler"

var statik = func() string {
	result, err := which("statik")
	if err != nil {
		os.Stderr.WriteString("statik not installed; did you run go get github.com/rakyll/statik?")
		os.Exit(1)
	}
	return result
}()

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

// ------------------------------------------------------------
// Targets

//goland:noinspection GoUnusedExportedType
type Assets mg.Namespace

// Embeds static assets
func (Assets) Embed() error {
	// statik -src=resources
	cmd := exec.Command(statik, "-Z", "-src", "resources", "-ns", projectName)

	println(strings.Join(cmd.Args, " "))

	return cmd.Run()
}
