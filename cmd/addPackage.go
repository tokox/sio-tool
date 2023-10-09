package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/otiai10/copy"
)

const ErrorFileIsNotADirectory = "this file is not a directory"

func AddPackage() (err error) {
	fileInfo, err := os.Stat(Args.File)
	if err != nil {
		return
	}
	destination, err := ArgsPackagePath()
	if err != nil {
		return
	}
	destination = getPackageNumber(destination)

	if !fileInfo.IsDir() {
		return errors.New(ErrorFileIsNotADirectory)
	}
	color.Green("Coping package to destination: %v", destination)
	err = copy.Copy(Args.File, destination)
	if err != nil {
		return
	}
	color.Green("Successfully copied package")
	deleteFolder := true
	if err = survey.AskOne(&survey.Confirm{Message: `Do you want to delete the original folder?`, Default: true}, &deleteFolder); err != nil {
		return
	}
	if deleteFolder {
		return os.RemoveAll(Args.File)
	}
	return
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func getPackageNumber(path string) (packagePath string) {
	i := 0
	for {
		packagePath = filepath.Join(path, strconv.Itoa(i))
		if !fileExists(packagePath) {
			return
		}
		i++
	}
}
