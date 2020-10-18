// +build mage

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bep/golibsass/libsass"
	"github.com/magefile/mage/mg"
)

// TODO: clean up output
// TODO: use github.com/magefile/mage/target to skip unnecessary steps

// ------------------------------------------------------------
// Constants and const-like variables

const projectName = "adler"
const mainScssPath = "scss/main.scss"

var statik = ensureCommand("statik", "statik not found; did you run go get github.com/rakyll/statik?")
var sassLint = ensureCommand("sass-lint", "sass-lint not found; did you run npm install -g sass-lint?")

// ------------------------------------------------------------
// Targets

//goland:noinspection GoUnusedExportedType
type Assets mg.Namespace

// Embeds static assets
func (Assets) Embed() error {
	mg.Deps(Assets.Compile)

	cmd := exec.Command(statik, "-Z", "-src", "resources", "-ns", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println(strings.Join(cmd.Args, " "))

	return cmd.Run()
}

// Validates SCSS
func (Assets) Validate() error {
	cmd := exec.Command(sassLint, "-v", "--max-warnings", "0", "-c", "scss/.sass-lint.yml", mainScssPath)
	cmd.Stdout = os.Stdout
	if mg.Verbose() {
		cmd.Stderr = os.Stderr
	}

	println(strings.Join(cmd.Args, " "))

	return cmd.Run()
}

// Compiles SCSS
func (Assets) Compile() error {
	mg.Deps(Assets.Validate)

	println("Reading SCSS from " + mainScssPath)
	mainScss, err := readFileAsString(mainScssPath)
	if err != nil {
		return err
	}

	scssDir := filepath.Dir(mainScssPath)
	//println("Scanning " + scssDir + " for includes")
	//
	//includes, err := filepath.Glob(scssDir + "/_*.scss")
	//if err != nil {
	//	return err
	//}
	//msg := fmt.Sprintf("Found includes: %v", strings.Join(includes, ", "))
	//println(msg)

	println("Initializing transpiler")
	transpiler, _ := libsass.New(libsass.Options{
		IncludePaths: []string{scssDir},
		OutputStyle:  libsass.ExpandedStyle,
	})

	println("Transpiling")
	result, err := transpiler.Execute(mainScss)
	if err != nil {
		return err
	}

	outputDir := "resources/css"
	println("Ensuring output directory: " + outputDir)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "main.css")
	println("Writing output file: " + outputFile)
	return ioutil.WriteFile(outputFile, []byte(result.CSS), 0644)
}

// ------------------------------------------------------------
// Helper functions

func readFileAsString(path string) (string, error) {
	f, err := os.Open(path)
	defer closeQuietly(f)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if _, err := io.Copy(&sb, f); err != nil {
		return "", err
	}
	return sb.String(), nil
}

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

type closeable interface {
	Close() error
}

func closeQuietly(cl closeable) {
	if cl != nil {
		err := cl.Close()
		if err != nil {
			msg := fmt.Sprintf("Error closing %v: %v\n", cl, err)
			os.Stderr.WriteString(msg)
		}
	}
}
