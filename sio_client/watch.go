package sio_client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Arapak/sio-tool/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

type Submission struct {
	name      string
	shortName string
	id        uint64
	status    string
	points    uint64
	when      string
	end       bool
}

func (s *Submission) ParseStatus() string {
	status := s.status
	for k, v := range colorMap {
		tmp := strings.ReplaceAll(status, k, "")
		if tmp != status {
			status = color.New(v).Sprint(tmp)
		}
	}
	return status
}

func (s *Submission) ParseID() string {
	return fmt.Sprintf("%v", s.id)
}

const inf = 1000000009

var intReg = regexp.MustCompile(`\d+`)

func toInt(sel string) uint64 {
	if tmp := intReg.FindString(sel); tmp != "" {
		t, _ := strconv.Atoi(tmp)
		return uint64(t)
	}
	return inf
}

func (s *Submission) ParsePoints() string {
	if s.points == inf {
		return ""
	} else if s.points == 0 {
		return color.New(colorMap["${c-failed}"]).Sprint(s.points)
	} else if s.points < 100 {
		return color.New(colorMap["${c-partial}"]).Sprint(s.points)
	} else if s.points == 100 {
		return color.New(colorMap["${c-accepted}"]).Sprint(s.points)
	}
	return fmt.Sprintf("%v", s.points)
}

func refreshLine(n int, maxWidth int) {
	for i := 0; i < n; i++ {
		_, _ = ansi.Printf("%v\n", strings.Repeat(" ", maxWidth))
	}
	ansi.CursorUp(n)
}

func updateLine(line string, maxWidth *int) string {
	*maxWidth = len(line)
	return line
}

func (s *Submission) display(first bool, maxWidth *int) {
	if !first {
		ansi.CursorUp(6)
	}
	_, _ = ansi.Printf("      #: %v\n", s.ParseID())
	_, _ = ansi.Printf("   when: %v\n", s.when)
	_, _ = ansi.Printf("   prob: %v\n", s.name)
	_, _ = ansi.Printf("  alias: %v\n", s.shortName)
	refreshLine(1, *maxWidth)
	_, _ = ansi.Printf(updateLine(fmt.Sprintf(" status: %v\n", s.ParseStatus()), maxWidth))
	_, _ = ansi.Printf(" points: %v\n", s.ParsePoints())
}

func display(submissions []Submission, first bool, maxWidth *int, line bool) {
	if line {
		submissions[0].display(first, maxWidth)
		return
	}
	var buf bytes.Buffer
	output := io.Writer(&buf)
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"#", "when", "problem", "alias", "status", "points"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, sub := range submissions {
		table.Append([]string{
			sub.ParseID(),
			sub.when,
			sub.name,
			sub.shortName,
			sub.ParseStatus(),
			sub.ParsePoints(),
		})
	}
	table.Render()

	if !first {
		ansi.CursorUp(len(submissions) + 2)
	}
	refreshLine(len(submissions)+2, *maxWidth)

	scanner := bufio.NewScanner(io.Reader(&buf))
	for scanner.Scan() {
		line := scanner.Text()
		*maxWidth = len(line)
		_, _ = ansi.Println(line)
	}
}

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

func parseSubmission(s *goquery.Selection) (ret Submission, err error) {
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
	points := toInt(get(fmt.Sprintf("td#submission%v-score", id)))
	kind := get(fmt.Sprintf("td#submission%v-kind", id))
	status := get(fmt.Sprintf("td#submission%v-status", id))
	statusLowercase := strings.ToLower(status)
	end := true
	if strings.Contains(statusLowercase, "oczekuje") || strings.Contains(statusLowercase, "pending") {
		status = fmt.Sprintf("${c-waiting}%v", status)
		end = false
	} else if strings.Contains(statusLowercase, "ok") {
		status = fmt.Sprintf("${c-accepted}%v", status)
		if points == inf && (kind == "" || kind == "Normalne" || kind == "Normal" || kind == "Zignorowane" || kind == "Ignored") {
			end = false
		}
	} else if strings.Contains(statusLowercase, "błąd") || strings.Contains(statusLowercase, "failed") {
		status = fmt.Sprintf("${c-failed}%v", status)
		if points == inf && !strings.Contains(statusLowercase, "kompilacji") && !strings.Contains(statusLowercase, "compilation") {
			end = false
		}
	} else {
		status = fmt.Sprintf("${c-rejected}%v", status)
	}
	return Submission{
		id:        toInt(id),
		name:      name,
		shortName: shortName,
		status:    status,
		points:    points,
		when:      when,
		end:       end,
	}, nil
}

func (c *SioClient) getSubmissions(URL string, n int) (submissions []Submission, err error) {
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

func (c *SioClient) WatchSubmission(info Info, n int, line bool) (submissions []Submission, err error) {
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
		submissions, err = c.getSubmissions(URL, n)
		if err != nil {
			return
		}
		display(submissions, first, &maxWidth, line)
		first = false
		endCount := 0
		for _, submission := range submissions {
			if submission.end {
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

var colorMap = map[string]color.Attribute{
	"${c-waiting}":  color.FgBlue,
	"${c-failed}":   color.FgRed,
	"${c-accepted}": color.FgGreen,
	"${c-partial}":  color.FgCyan,
	"${c-rejected}": color.FgBlue,
}
