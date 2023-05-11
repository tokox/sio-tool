package szkopul_client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

const SubmitIDRegStr = `\d+`

func (c *SzkopulClient) Submit(info Info, sourcePath string) (err error) {
	color.Cyan("Submit " + info.Hint())

	URL, err := info.APISubmitURL(c.host)
	if err != nil {
		return
	}
	refererURL, err := info.SubmitURL(c.host)
	if err != nil {
		return
	}

	fmt.Printf("Current user: %v\n", c.Username)

	csrf, err := c.GetCsrf(refererURL)
	if err != nil {
		return
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(sourceFile.Name()))
	io.Copy(part, sourceFile)
	part, err = writer.CreateFormField("csrfmiddlewaretoken")
	io.Copy(part, strings.NewReader(csrf))
	writer.Close()

	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Referer", refererURL)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	isSubmitID, err := regexp.MatchString(SubmitIDRegStr, string(responseBody))

	if isSubmitID {
		color.Green("Submitted")

		submissions, err := c.WatchSubmission(info, 1, true)
		if err != nil {
			return err
		}

		info.SubmissionID = submissions[0].ParseID()
		c.LastSubmission = &info
	} else {
		fmt.Print("an error occured: ")
		color.Red(string(responseBody))
	}
	return c.save()
}
