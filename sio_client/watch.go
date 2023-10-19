package sio_client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Arapak/sio-tool/sio_submissions"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"

	"github.com/PuerkitoBio/goquery"
)

func getSubmissionID(body string) (string, error) {
	reg := regexp.MustCompile(`<td id="submission(\d+?)-score">`)
	tmp := reg.FindStringSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("cannot find submission id")
	}
	return tmp[1], nil
}

func findSubmission(body []byte) (submissions []*goquery.Selection, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	doc.Find("table.submission tbody").First().Find("tr").Each(func(_ int, s *goquery.Selection) {
		submissions = append(submissions, s)
	})
	return
}

func getProblemNames(name string) (string, string) {
	reg := regexp.MustCompile(`([\s\S]+?) \((\S+?)\)`)
	tmp := reg.FindSubmatch([]byte(name))
	if len(tmp) < 3 {
		return name, ""
	}
	return string(tmp[1]), string(tmp[2])
}

func parseSubmission(s *goquery.Selection) (ret sio_submissions.Submission, err error) {
	body, err := s.Html()
	if err != nil {
		return
	}
	id, err := getSubmissionID(body)
	if err != nil {
		return
	}

	get := func(sel string) string {
		return strings.TrimSpace(s.Find(sel).Text())
	}
	when := strings.TrimSpace(s.Find("a").First().Text())
	combinedName := get(fmt.Sprintf("td#submission%v-problem-instance", id))
	name, shortName := getProblemNames(combinedName)
	points := sio_submissions.ToInt(get(fmt.Sprintf("td#submission%v-score", id)))
	kind := get(fmt.Sprintf("td#submission%v-kind", id))
	status := get(fmt.Sprintf("td#submission%v-status", id))
	statusLowercase := strings.ToLower(status)
	end := true
	if strings.Contains(statusLowercase, "oczekuje") || strings.Contains(statusLowercase, "pending") {
		status = fmt.Sprintf("${c-waiting}%v", status)
		end = false
	} else if strings.Contains(statusLowercase, "ok") {
		status = fmt.Sprintf("${c-accepted}%v", status)
		if points == sio_submissions.Inf && (kind == "" || kind == "Normalne" || kind == "Normal" || kind == "Zignorowane" || kind == "Ignored") {
			end = false
		}
	} else if strings.Contains(statusLowercase, "błąd") || strings.Contains(statusLowercase, "failed") {
		status = fmt.Sprintf("${c-failed}%v", status)
		if points == sio_submissions.Inf && !strings.Contains(statusLowercase, "kompilacji") && !strings.Contains(statusLowercase, "compilation") {
			end = false
		}
	} else {
		status = fmt.Sprintf("${c-rejected}%v", status)
	}
	return sio_submissions.Submission{
		Id:        sio_submissions.ToInt(id),
		Name:      name,
		ShortName: shortName,
		Status:    status,
		Points:    points,
		When:      when,
		End:       end,
	}, nil
}

func (c *SioClient) getSubmissions(URL string, n int) (submissions []sio_submissions.Submission, err error) {
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	if _, err = findUsername(body); err != nil {
		return
	}

	submissionsBody, err := findSubmission(body)
	if err != nil {
		return
	}

	for _, submissionBody := range submissionsBody {
		if submission, err := parseSubmission(submissionBody); err == nil {
			submissions = append(submissions, submission)
		}
		if len(submissions) == n {
			break
		}
	}

	if len(submissions) < 1 {
		return nil, errors.New("cannot find any submission")
	}

	return
}

func (c *SioClient) RevealSubmission(info Info) (err error) {
	submissionURL, err := info.SubmissionURL(c.host, false)
	if err != nil {
		return
	}
	body, err := util.GetBody(c.client, submissionURL)
	if err != nil {
		return
	}
	if !bytes.Contains(body, []byte("<h4>Score revealing</h4>")) && !bytes.Contains(body, []byte("<h4>Ujawnianie wyniku</h4>")) {
		return
	}
	if bytes.Contains(body, []byte("<p>Unfortunately, this submission has not been scored yet, so you can&#39;t see your score. Please come back later.</p>")) || bytes.Contains(body, []byte("<p>Niestety to zgłoszenie nie zostało jeszcze ocenione, więc nie możesz zobaczyć swojego wyniku. Spróbuj później.</p>")) {
		return
	}
	csrf, err := findCsrf(body)
	if err != nil {
		return
	}
	revealURL, err := info.SubmissionURL(c.host, true)
	if err != nil {
		return
	}
	postBody := &bytes.Buffer{}
	writer := multipart.NewWriter(postBody)
	part, err := writer.CreateFormField("csrfmiddlewaretoken")
	if err != nil {
		return
	}
	_, err = io.Copy(part, strings.NewReader(csrf))
	if err != nil {
		return
	}
	writer.Close()
	req, err := http.NewRequest("POST", revealURL, postBody)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Referer", submissionURL)
	_, err = c.client.Do(req)
	return
}

func (c *SioClient) WatchSubmission(info Info, n int, line bool) (submissions []sio_submissions.Submission, err error) {
	URL, err := info.MySubmissionURL(c.host)
	if err != nil {
		return
	}

	maxWidth := 0
	first := true
	for {
		st := time.Now()
		if info.SubmissionID != "" {
			err = c.RevealSubmission(info)
			if err != nil {
				return
			}
		}
		if c.instanceClient == Staszic {
			submissions, err = c.getSubmissions(URL, n)
		} else {
			submissions, err = szkopul_client.GetSubmissions(c.client, URL, n)
		}
		if err != nil {
			return
		}
		sio_submissions.Display(submissions, first, &maxWidth, line)
		first = false
		endCount := 0
		for _, submission := range submissions {
			if submission.End {
				endCount++
			}
		}
		if endCount == len(submissions) {
			return
		}

		if n == 1 && len(submissions) == 1 {
			info.SubmissionID = submissions[0].ParseID()
		}

		sub := time.Since(st)
		if sub < time.Second {
			time.Sleep(time.Second - sub)
		}
	}
}
