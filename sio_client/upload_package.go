package sio_client

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arapak/sio-tool/util"
	"github.com/fatih/color"
)

const ErrorProblemNotFound = "problem not found"

func (c *SioClient) uploadPackageFile(url string, file string) (err error) {
	csrf, err := c.getCsrf(url)

	packageFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer packageFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("package_file", filepath.Base(packageFile.Name()))
	if err != nil {
		return
	}
	_, err = io.Copy(part, packageFile)
	if err != nil {
		return
	}
	part, err = writer.CreateFormField("csrfmiddlewaretoken")
	if err != nil {
		return
	}
	_, err = io.Copy(part, strings.NewReader(csrf))
	if err != nil {
		return
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Referer", url)

	_, err = c.client.Do(req)
	return
}

func (c *SioClient) UploadPackage(info Info, file string) (perf util.Performance, err error) {
	packages, perf, err := c.FindAllPackages(info)
	if err != nil {
		return
	}
	var url string
	for _, p := range packages {
		if p.Name == info.ProblemAlias {
			url, err = info.ReuploadPackageURL(c.host, p.ReuploadId)
			if err != nil {
				return
			}
			err = c.uploadPackageFile(url, file)
			if err == nil {
				color.Green("Uploaded package for %v", p.Name)
			}
			return
		}
	}
	err = errors.New(ErrorProblemNotFound)
	return
}
