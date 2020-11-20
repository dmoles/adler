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
	"strconv"
	"strings"
	"time"

	"github.com/bep/golibsass/libsass"
	"github.com/get-woke/go-gitignore"
	"github.com/magefile/mage/mg"
)

// TODO: use github.com/magefile/mage/target to skip unnecessary steps

// ------------------------------------------------------------
// Constants and const-like variables

const projectName = "adler"
const mainCss = "resources/css/main.css"
const mainScss = "scss/main.scss"
const resourceDir = "resources"
const statikData = "statik/statik.go"

var scssDir = filepath.Dir(mainScss)
var gitIgnore = compileGitIgnore()

// ------------------------------------------------------------
// Targets

// builds an executable, but does not install it (depends on: test)
func Build() error {
	mg.Deps(Assets.Embed)

	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println("Building")
	if mg.Verbose() {
		println(strings.Join(cmd.Args, " "))
	}

	return cmd.Run()
}

// builds and installs the executable (depends on: test)
func Install() error {
	mg.Deps(Test)

	cmd := exec.Command("go", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println("Installing")
	if mg.Verbose() {
		println(strings.Join(cmd.Args, " "))
	}

	return cmd.Run()
}

// runs all tests (depends on: assets:embed)
func Test() error {
	mg.Deps(Assets.Embed)

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

// embeds static assets (depends on: assets:compile; requires statik: https://github.com/rakyll/statik)
func (Assets) Embed() error {
	mg.Deps(Assets.Compile)

	if !anyNewerThan(resourceDir, statikData) {
		println("Assets are up to date") // TODO: more consistent output
		return nil
	}

	includes := strings.Join([]string{
		"*.css",
		"*.ico",
		"*.md",
		"*.png",
		"*.tmpl",
		"*.webmanifest",
		"*.woff",
		"*.woff2",
	}, ",")

	var statik = ensureCommand("statik", "statik not found; did you run go get github.com/rakyll/statik?")
	cmd := exec.Command(statik, "-Z", "-src", resourceDir, "-include", includes, "-ns", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println("Embedding static assets")
	if mg.Verbose() {
		println(strings.Join(cmd.Args, " "))
	}

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

	if !anyNewerThan(resourceDir, mainCss) {
		println("CSS is up to date") // TODO: more consistent output
		return nil
	}

	libsassOptions := libsass.Options{
		IncludePaths: []string{scssDir},
		OutputStyle:  libsass.ExpandedStyle,
	}

	println(fmt.Sprintf("Transpiling %v to %v", mainScss, mainCss))
	if mg.Verbose() {
		msg := fmt.Sprintf("libsass options: %#v", libsassOptions)
		println(msg)
	}

	if mg.Verbose() {
		println("Reading " + mainScss)
	}
	mainScss, err := readFileAsString(mainScss)
	if err != nil {
		return err
	}

	transpiler, _ := libsass.New(libsassOptions)
	result, err := transpiler.Execute(mainScss)
	if err != nil {
		return err
	}

	if mg.Verbose() {
		println("Ensuring output directory " + ("resources/css"))
	}
	err = os.MkdirAll("resources/css", 0755)
	if err != nil {
		return err
	}

	if mg.Verbose() {
		println("Writing " + mainCss)
	}
	return ioutil.WriteFile(mainCss, []byte(result.CSS), 0644)
}

// ------------------------------------------------------------
// Helper functions

// TODO: put all this timestamp business in a struct, or find a utility library for it (does mage have one?)

var timeZero = time.Time{}

func anyNewerThan(sourceDir string, targetFile string) bool {
	targetModTime, err := effectiveModTime(targetFile)
	if err != nil {
		return true
	}
	if mg.Verbose() {
		fmt.Printf("%v modified %v\n", targetFile, targetModTime)
	}

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		modTime, err := effectiveModTime(path)
		if err != nil {
			return err
		}
		if modTime.After(targetModTime) {
			if mg.Verbose() {
				fmt.Printf("%v modified %v\n", path, modTime)
			}
			return io.EOF
		}
		return nil
	})
	if err == nil {
		return false
	}
	if err == io.EOF {
		return true
	}
	panic(err)
}

func effectiveModTime(path string) (time.Time, error) {
	if !tracked(path) || changed(path) {
		return modTime(path)
	}
	return commitTime(path)
}

func tracked(path string) bool {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", path)
	return cmd.Run() == nil
}

func changed(path string) bool {
	cmd := exec.Command("git", "diff", "--quiet", path)
	return cmd.Run() != nil
}

func modTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return timeZero, err
	}
	return info.ModTime(), nil
}

func commitTime(path string) (time.Time, error) {
	ts, err := commitTimestamp(path)
	if err != nil {
		return timeZero, err
	}
	return time.Unix(ts, 0), nil
}

func commitTimestamp(path string) (int64, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%ct", path)
	tsBytes, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	tsStr := strings.TrimSpace(string(tsBytes))
	if len(tsStr) == 0 {
		return 0, err
	}
	return strconv.ParseInt(tsStr, 10, 64)
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

func compileGitIgnore() *ignore.GitIgnore {
	gi, err := ignore.CompileIgnoreFile(".gitignore")
	if err != nil {
		panic(err)
	}
	return gi
}
