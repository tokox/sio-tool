package szkopul_client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func (c *SzkopulClient) Submit(info Info, sourcePath string) (err error) {
	color.Cyan("Submit " + info.Hint())

	URL, err := info.SubmitURL(c.host)
	if err != nil {
		return
	}

	fmt.Printf("Current user: %v\n", c.Username)

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(sourceFile.Name()))
	io.Copy(part, sourceFile)
	writer.Close()

	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	token, err := c.DecryptToken()
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %v", token))

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	color.Green("Submitted")

	info.SubmissionID = string(responseBody)
	c.LastSubmission = &info
	return c.save()
}
