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
	"github.com/get-woke/go-gitignore"
	"github.com/magefile/mage/mg"
)

// TODO: clean up output
// TODO: use github.com/magefile/mage/target to skip unnecessary steps

// ------------------------------------------------------------
// Constants and const-like variables

const projectName = "adler"
const mainScssPath = "scss/main.scss"

var scssDir = filepath.Dir(mainScssPath)
var gitIgnore = compileGitIgnore()

// ------------------------------------------------------------
// Targets

//goland:noinspection GoUnusedExportedType
type Assets mg.Namespace

// embeds static assets (depends on: assets:compile; requires statik: https://github.com/rakyll/statik)
func (Assets) Embed() error {
	mg.Deps(Assets.Compile)

	includes := strings.Join([]string{
		"*.css",
		"*.ico",
		"*.md",
		"*.png",
		"*.tmpl",
		"*.woff",
		"*.woff2",
	}, ",")

	var statik = ensureCommand("statik", "statik not found; did you run go get github.com/rakyll/statik?")
	cmd := exec.Command(statik, "-Z", "-src", "resources", "-include", includes, "-ns", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println(strings.Join(cmd.Args, " "))

	return cmd.Run()
}

// validates SCSS (requires sass-lint: https://www.npmjs.com/package/sass-lint)
func (Assets) Validate() error {
	var errors []error

	//goland:noinspection GoUnhandledErrorResult
	filepath.Walk(scssDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".scss" && !info.IsDir() && !gitIgnore.MatchesPath(path) {
			if err := sassLint(path); err != nil {
				errors = append(errors, err)
			}
		}
		return nil
	})

	if len(errors) > 0 {
		for _, err := range errors {
			warn(err.Error())
		}
		return errors[len(errors)-1]
	}
	return nil
}

// compiles SCSS (depends on: assets:validate)
func (Assets) Compile() error {
	mg.Deps(Assets.Validate)

	println("Reading SCSS from " + mainScssPath)
	mainScss, err := readFileAsString(mainScssPath)
	if err != nil {
		return err
	}

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

func sassLint(scssFile string) error {
	var sassLint = ensureCommand("sass-lint", "sass-lint not found; did you run npm install -g sass-lint?")
	cmd := exec.Command(sassLint, "-v", "--max-warnings", "0", "-c", "scss/.sass-lint.yml", scssFile)
	cmd.Stdout = os.Stdout
	if mg.Verbose() {
		cmd.Stderr = os.Stderr
		println(strings.Join(cmd.Args, " "))
	}
	return cmd.Run()
}

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
		warn(failureMsg)
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

func warn(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

func compileGitIgnore() *ignore.GitIgnore {
	gi, err := ignore.CompileIgnoreFile(".gitignore")
	if err != nil {
		panic(err)
	}
	return gi
}

type closeable interface {
	Close() error
}

func closeQuietly(cl closeable) {
	if cl != nil {
		err := cl.Close()
		if err != nil {
			msg := fmt.Sprintf("Error closing %v: %v", cl, err)
			warn(msg)
		}
	}
}
