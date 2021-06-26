// +build mage

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
)

// TODO: use github.com/magefile/mage/target to skip unnecessary steps

// ------------------------------------------------------------
// Constants and const-like variables

const mainCss = "resources/css/main.css"
const mainScss = "scss/main.scss"

// TODO: figure out how to document these and/or make them CL options
const envSkipTests = "ADLER_SKIP_TESTS"
const envSkipValidation = "ADLER_SKIP_VALIDATION"

var scssDir = filepath.Dir(mainScss)
var cssDir = filepath.Dir(mainCss)

// ------------------------------------------------------------
// Targets

// Build builds an executable, but does not install it (depends on: test)
//goland:noinspection GoUnusedExportedFunction
func Build() error {
	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println("Building")
	if mg.Verbose() {
		println(strings.Join(cmd.Args, " "))
	}

	return cmd.Run()
}

// Install builds and installs the executable (depends on: test)
//goland:noinspection GoUnusedExportedFunction
func Install() error {
	mg.Deps(Test, Assets.Compile)

	cmd := exec.Command("go", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println("Installing")
	if mg.Verbose() {
		println(strings.Join(cmd.Args, " "))
	}

	return cmd.Run()
}

// Test runs all tests
func Test() error {
	if skipTests() {
		warn("Skipping tests")
		return nil
	}

	cmd := exec.Command("go", "test", "./...")
	cmd.Stderr = os.Stderr

	println("Running tests")
	if mg.Verbose() {
		println(strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}

	var sb strings.Builder
	cmd.Stdout = &sb
	err := cmd.Run()
	if err != nil {
		print(sb.String())
	}
	return err
}

//goland:noinspection GoUnusedExportedType
type Assets mg.Namespace

// Validate validates SCSS (requires sass-lint: https://www.npmjs.com/package/sass-lint)
func (Assets) Validate() error {
	if skipValidation() {
		warn("Skipping validation")
		return nil
	}

	var errors []error

	//goland:noinspection GoUnhandledErrorResult
	filepath.Walk(scssDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".scss" && !info.IsDir() && !ignored(path) {
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

// Compile compiles SCSS (depends on: assets:validate)
func (Assets) Compile() error {
	mg.Deps(Assets.Validate)

	if mg.Verbose() {
		println("Ensuring output directory " + cssDir)
	}
	err := os.MkdirAll(cssDir, 0755)
	if err != nil {
		return err
	}

	sass := ensureCommand("sass", "sass not found; did you run brew install sass/sass/sass or npm install -g sass?")

	var sassQuietArg string
	if mg.Verbose() {
		sassQuietArg = "--no-quiet"
	} else {
		sassQuietArg = "--quiet"
	}

	cmd := exec.Command(sass, sassQuietArg, "--stop-on-error", scssDir+":"+cssDir)
	cmd.Stdout = os.Stdout
	if mg.Verbose() {
		cmd.Stderr = os.Stderr
		println(strings.Join(cmd.Args, " "))
	}
	return cmd.Run()
}

// ------------------------------------------------------------
// Helper functions

func skipTests() bool {
	return os.Getenv(envSkipTests) != ""
}

func skipValidation() bool {
	return os.Getenv(envSkipValidation) != ""
}

// TODO: put all this timestamp business in a struct, or find a utility library for it (does mage have one?)

var timeZero = time.Time{}

func ignored(path string) bool {
	gi, err := gitIgnore()
	if err != nil {
		panic(err)
	}
	return gi.MatchesPath(path)
}

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
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buf), nil
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
