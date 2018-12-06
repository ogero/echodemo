// +build mage

package main

import (
	"errors"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var (
	outputFile  = "echodemo"
	userID      = "echodemo"
	userGroup   = "echodemo"
	serviceName = "echodemo"
	deployPath  = "/home/echodemo/app"
	packageName = "bitbucket.org/ogero/echodemo"
)

func init() {
	if runtime.GOOS == "windows" {
		outputFile += ".exe"
	}
}

var Default = Build

// Dev starts realize
func Dev() error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "realize", "start")
	return err
}

// Clean deletes generated files from dist folder
func Clean() error {
	fmt.Println("Cleaning...")
	if err := sh.Rm("www/assets/css/default.bundle.css"); err != nil {
		return err
	}
	if err := sh.Rm("dist/" + outputFile); err != nil {
		return err
	}
	return nil
}

// Generate performs sass pre-processing and fileb0x embedding
func Generate() error {
	fmt.Println("Generating...")
	if err := sh.Run("sassc", "--style", "compressed", "www/assets/css/default.scss", "www/assets/css/default.css"); err != nil {
		return err
	}
	if err := copy("www/assets/materialize-src/js/bin/materialize.min.js", "www/assets/js/materialize.min.js", true); err != nil {
		return err
	}
	if err := sh.Run("fileb0x", "embed.json"); err != nil {
		return err
	}
	return nil
}

// Build performs a clean, generates required files and builds the binary placing it at dist folder
func Build() error {
	mg.Deps(Clean)
	mg.Deps(Generate)

	fmt.Println("Obtaining git info...")
	commitHash, _ := sh.Output("git", "log", "-1", "--format=\"%H\"")
	gitTag, _ := sh.Output("git", "describe", "--tags", "--abbrev=0")

	fmt.Println("Building...")
	args := []string{
		"build",
		"-ldflags",
		"-s -w -X \"" + packageName + "/dist.Timestamp=" + time.Now().Format(time.RFC1123) + "\"" +
			" -X \"" + packageName + "/dist.CommitHash=" + strings.Trim(commitHash, "\"") + "\"" +
			" -X \"" + packageName + "/dist.GitTag=" + strings.Trim(gitTag, "\"") + "\"",
		"-o",
		path.Join("dist/", outputFile),
		packageName,
	}
	_, err := sh.Output("go", args...)
	return err
}

// Deploy performs a build, stops the service, reinstall the binary and restart the service
func Deploy() error {
	if runtime.GOOS == "windows" {
		return errors.New("Can't deploy on windows")
	}

	fmt.Println("Pulling last sources...")
	if err := sh.Run("git", "reset", "--hard"); err != nil {
		return err
	}
	if err := sh.Run("git", "pull"); err != nil {
		return err
	}
	if err := sh.Run("git", "fetch", "--tags"); err != nil {
		return err
	}

	fmt.Println("Dep ensure...")
	if err := sh.Run("dep", "ensure"); err != nil {
		return err
	}

	mg.Deps(Build)
	fmt.Println("Deploying...")

	fmt.Println("Stopping service...")
	if err := sh.Run("service", serviceName, "stop"); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Copying binary...")
	if err := sh.Run("cp", path.Join("dist/", outputFile), path.Join(deployPath, outputFile)); err != nil {
		fmt.Println(err)
	}
	if err := sh.Run("chown", userID+":"+userGroup, path.Join(deployPath, outputFile)); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Starting service...")
	if err := sh.Run("service", serviceName, "start"); err != nil {
		return err
	}

	return nil
}

// copy the src file to dst
func copy(src, dst string, skipIfExists ...bool) error {
	if len(skipIfExists) == 1 && skipIfExists[0] {
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
