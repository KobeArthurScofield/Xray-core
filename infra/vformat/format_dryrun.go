package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of vformat:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if !filepath.IsAbs(*directory) {
		pwd, wdErr := os.Getwd()
		if wdErr != nil {
			fmt.Println("Can not get current working directory.")
			os.Exit(1)
		}
		*directory = filepath.Join(pwd, *directory)
	}

	pwd := *directory
	GOBIN := GetGOBIN()
	binPath := os.Getenv("PATH")
	pathSlice := []string{pwd, GOBIN, binPath}
	binPath = strings.Join(pathSlice, string(os.PathListSeparator))
	os.Setenv("PATH", binPath)

	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	gofmt := "gofumpt" + suffix
	goimports := "gci" + suffix

	if gofmtPath, err := exec.LookPath(gofmt); err != nil {
		fmt.Println("Can not find", gofmt, "in system path or current working directory.")
		os.Exit(1)
	} else {
		gofmt = gofmtPath
	}

	if goimportsPath, err := exec.LookPath(goimports); err != nil {
		fmt.Println("Can not find", goimports, "in system path or current working directory.")
		os.Exit(1)
	} else {
		goimports = goimportsPath
	}

	rawFilesSlice := make([]string, 0, 1000)
	walkErr := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".go") &&
			!strings.HasSuffix(filename, ".pb.go") &&
			!strings.Contains(dir, filepath.Join("testing", "mocks")) &&
			!strings.Contains(path, filepath.Join("main", "distro", "all", "all.go")) {
			rawFilesSlice = append(rawFilesSlice, path)
		}

		return nil
	})
	if walkErr != nil {
		fmt.Println(walkErr)
		os.Exit(1)
	}

	gofmtListArgs := []string{
		"-l", "-e",
	}

	gofmtShowArgs := []string{
		"-d", "-e",
	}

	goimportsListArgs := []string{
		"list",
	}

	goimportsShowArgs := []string{
		"diff",
	}

	fmt.Println("Checking files thar are not properly formatted...")
	RunMany(gofmt, gofmtListArgs, rawFilesSlice)
	RunMany(goimports, goimportsListArgs, rawFilesSlice)
	RunMany(gofmt, gofmtShowArgs, rawFilesSlice)
	RunMany(goimports, goimportsShowArgs, rawFilesSlice)
	fmt.Println("Do NOT forget to commit file changes.")
}
