package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"

	"github.com/fatih/color"

	"github.com/PuerkitoBio/goquery"
)

func less(a, b string) bool {
	if a == "$CI_VERSION" {
		return true
	}
	if b == "$CI_VERSION" {
		return false
	}
	reg := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	x := reg.FindSubmatch([]byte(a))
	y := reg.FindSubmatch([]byte(b))
	num := func(s []byte) int {
		n, _ := strconv.Atoi(string(s))
		return n
	}
	for i := 1; i <= 3; i++ {
		if num(x[i]) < num(y[i]) {
			return true
		} else if num(x[i]) > num(y[i]) {
			return false
		}
	}
	return false
}

func getLatest() (version, ptime, url string, err error) {
	goos := ""
	switch runtime.GOOS {
	case "darwin":
		goos = "darwin"
	case "linux":
		goos = "linux"
	case "windows":
		goos = "windows"
	default:
		err = fmt.Errorf("not support %v", runtime.GOOS)
		return
	}

	arch := ""
	switch runtime.GOARCH {
	case "386":
		arch = "x32"
	case "amd64":
		arch = "x64"
	default:
		err = fmt.Errorf("not support %v", runtime.GOARCH)
		return
	}

	resp, err := http.Get("https://github.com/tokox/sio-tool/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	location := resp.Request.URL

	version = location.Path[strings.LastIndex(location.Path, "/")+1:]

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}


	tm, _ := time.Parse("2006-01-02T15:04:05Z", doc.Find("relative-time").First().AttrOr("datetime", ""))
	ptime = tm.In(time.Local).Format("2006-01-02 15:04")
	url = fmt.Sprintf("https://github.com/tokox/sio-tool/releases/download/%v/st_%v_%v.zip", version, goos, arch)

	return
}

type WriteCounter struct {
	Count uint
	last  uint
}

func (w *WriteCounter) Print() {
	fmt.Printf("\rProgress: %v KB  Speed: %v KB/s           ",
		w.Count/1024, (w.Count-w.last)/1024)
	w.last = w.Count
}

func (w *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	w.Count += uint(n)
	return n, nil
}

func upgrade(url, exePath string) (err error) {
	updateDir := filepath.Dir(exePath)

	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filepath.Base(exePath)))
	color.Cyan("Move the old one to %v", oldPath)
	if err = os.Rename(exePath, oldPath); err != nil {
		return
	}
	defer func() {
		if err != nil {
			color.Cyan("Move the old one back")
			if e := os.Rename(oldPath, exePath); e != nil {
				color.Red(e.Error())
			}
		} else {
			color.Cyan("Remove the old one")
			if e := os.Remove(oldPath); e != nil {
				color.Red(e.Error() + "\nYou could remove it manually")
			}
		}
	}()

	color.Cyan("Download %v", url)
	counter := &WriteCounter{Count: 0, last: 0}
	counter.Print()

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			counter.Print()
		}
	}()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.TeeReader(resp.Body, counter))
	ticker.Stop()
	counter.Print()
	fmt.Println()
	if err != nil {
		return
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return
	}

	rc, err := reader.File[0].Open()
	if err != nil {
		return
	}
	newData, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		return
	}

	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filepath.Base(exePath)))
	color.Cyan("Save the new one to %v", newPath)
	if err = os.WriteFile(newPath, newData, 0755); err != nil {
		return
	}

	if err = os.Rename(newPath, exePath); err != nil {
		color.Cyan("Delete the new one %v", newPath)
		if e := os.Remove(newPath); e != nil {
			color.Red(e.Error())
		}
	}

	return
}

func Upgrade() (err error) {
	color.Cyan("Checking version")
	latest, ptime, url, err := getLatest()
	if err != nil {
		return
	}
	version := Args.Version
	if !less(version, latest) {
		color.Green("Current version %v is the latest", version)
		return
	}

	color.Red("Current version is %v", version)
	color.Green("The latest version is %v, published at %v", latest, ptime)

	doUpgrade := true
	prompt := &survey.Confirm{Message: "Do you want to upgrade?", Default: true}
	if err = survey.AskOne(prompt, &doUpgrade); err != nil {
		return
	}
	if !doUpgrade {
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		return
	}

	if exePath, err = filepath.EvalSymlinks(exePath); err != nil {
		return
	}

	if err = upgrade(url, exePath); err != nil {
		return
	}

	color.Green("Successfully updated to version %v", latest)
	return
}
