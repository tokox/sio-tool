package sio_client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Arapak/sio-tool/util"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const ErrorNeedProblemIdentification = "you have to specify the problem alias or the problem instance id"
const ErrorNeedCantFindProblemID = "couldn't find problem instance id (maybe a problem with this alias doesn't exist in this contest)"

const SubmissionsPageRegExp = `<title>My submissions[\S\s]+</title>`
const LoginPageRegExp = `<h1>Log in</h1>`

func getErrorsFromBody(body []byte) (err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	s := doc.Find(".help-block").First().Text()
	if s != "" {
		return errors.New(s)
	}
	return
}

func findProblemID(body []byte, info *Info) (err error) {
	if info.ProblemAlias == "" {
		return errors.New(ErrorNeedProblemIdentification)
	}
	reg := regexp.MustCompile(fmt.Sprintf(`<option value="(?P<problemID>\d+?)">[\S ]+?\(%v\)</option>`, info.ProblemAlias))
	names := reg.SubexpNames()
	for i, val := range reg.FindSubmatch(body) {
		if names[i] == "problemID" {
			info.ProblemID = string(val)
		}
	}
	if info.ProblemID == "" {
		return errors.New(ErrorNeedCantFindProblemID)
	}
	return
}

func (c *SioClient) Submit(info Info, sourcePath string) (err error) {
	URL, err := info.SubmitURL(c.host)
	if err != nil {
		return
	}

	submitPageBody, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	_, err = findUsername(submitPageBody)
	if err != nil {
		return
	}

	csrf, err := findCsrf(submitPageBody)
	if err != nil {
		return
	}

	if info.ProblemID == "" {
		err = findProblemID(submitPageBody, &info)
		if err != nil {
			return
		}
	}

	color.Cyan("Submit " + info.Hint())
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
	part, err = writer.CreateFormField("csrfmiddlewaretoken")
	io.Copy(part, strings.NewReader(csrf))
	part, err = writer.CreateFormField("problem_instance_id")
	io.Copy(part, strings.NewReader(info.ProblemID))
	part, err = writer.CreateFormField("user")
	io.Copy(part, strings.NewReader(c.Username))
	part, err = writer.CreateFormField("kind")
	io.Copy(part, strings.NewReader("NORMAL"))
	writer.Close()

	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Referer", URL)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	isLoginPage, err := regexp.Match(LoginPageRegExp, responseBody)
	if err != nil {
		return
	}
	if isLoginPage {
		return errors.New(ErrorNotLogged)
	}

	isSubmissionsPage, err := regexp.Match(SubmissionsPageRegExp, responseBody)
	if err != nil {
		return
	}

	if isSubmissionsPage {
		color.Green("Submitted")

		submissions, err := c.WatchSubmission(info, 1, true)
		if err != nil {
			return err
		}

		info.SubmissionID = submissions[0].ParseID()
		c.LastSubmission = &info
	} else {
		fmt.Print("an error occurred: ")
		err = getErrorsFromBody(responseBody)
		if err != nil {
			color.Red(err.Error())
		}
	}
	return c.save()
}
