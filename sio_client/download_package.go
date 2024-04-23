package sio_client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/Arapak/sio-tool/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const ErrorNoFileAttached = "no file attached"

type PackageInfo struct {
	Name    string
	Alias   string
	Package string
}

func findPackages(body []byte) (packages []PackageInfo, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	doc.Find("table tbody").First().Find("tr").Each(func(_ int, s *goquery.Selection) {
		info := PackageInfo{}
		info.Name = strings.TrimSpace(s.Find(".field-name_link a").Last().Text())
		info.Alias = strings.TrimSpace(s.Find(".field-short_name_link a").First().Text())
		info.Package, _ = s.Find(".field-package a").First().Attr("href")
		packages = append(packages, info)
	})
	return
}

func (c *SioClient) DownloadPackage(p PackageInfo, rootPath string) (err error) {
	resp, err := c.client.Get(c.host + p.Package)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`attachment; filename="(\S+)"`)
	match := re.FindStringSubmatch(resp.Header.Get("Content-Disposition"))
	if match == nil {
		return errors.New(ErrorNoFileAttached)
	}
	filename := match[1]
	extension := path.Ext(filename)
	return os.WriteFile(path.Join(rootPath, p.Name+extension), body, 0644)
}

func (c *SioClient) FindAllPackages(info Info) (packages []PackageInfo, perf util.Performance, err error) {
	URL, err := info.ProblemInstanceURL(c.host)
	if err != nil {
		return
	}

	var body []byte
	var current_packages []PackageInfo
	perf.StartFetching()
	pageNum := 0
	for {
		body, err = util.GetBody(c.client, (fmt.Sprintf("%v/?p=%v", URL, pageNum)))
		if err != nil {
			return
		}
		if bytes.Contains(body, []byte("<p>404 &mdash; Page not found</p>")) {
			err = errors.New(ErrorContestNotFound)
			return
		}

		pageNum++

		if _, err = findUsername(body); err != nil {
			return
		}

		current_packages, err = findPackages(body)
		if err != nil {
			return
		}
		if len(current_packages) == 0 || (len(packages) != 0 && current_packages[0] == packages[0]) {
			break
		}
		packages = append(packages, current_packages...)
	}
	perf.StopFetching()
	return
}

func (c *SioClient) DownloadAllPackages(info Info, rootPath string) (perf util.Performance, err error) {
	packages, perf, err := c.FindAllPackages(info)
	if err != nil {
		return
	}
	numberOfWorkers := 10

	wg := sync.WaitGroup{}
	wg.Add(numberOfWorkers)
	mu := sync.Mutex{}

	workerError := false
	packageNumber := 0

	for i := 1; i <= numberOfWorkers; i++ {
		go func(workerID int) {
			defer func() {
				mu.Lock()
				workerError = true
				mu.Unlock()
				wg.Done()
			}()
			for {
				mu.Lock()
				if workerError || packageNumber >= len(packages) {
					mu.Unlock()
					return
				}
				currentPackage := packages[packageNumber]
				packageNumber++
				mu.Unlock()
				err = c.DownloadPackage(currentPackage, rootPath)
				if err != nil {
					mu.Lock()
					color.Red(err.Error())
					mu.Unlock()
					return
				}
				mu.Lock()
				color.Green("Downloaded package for task: %v", currentPackage.Name)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	color.Blue("----FINISHED----")

	return
}
