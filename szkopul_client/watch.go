package szkopul_client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sio-tool/util"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

// Submission submit state
type Submission struct {
	name       string
	short_name string
	id         uint64
	status     string
	points     uint64
	when       string
	end        bool
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

// ParseID formatter
func (s *Submission) ParseID() string {
	return fmt.Sprintf("%v", s.id)
}

const inf = 1000000009

// ParsePoints formatter
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
		ansi.Printf("%v\n", strings.Repeat(" ", maxWidth))
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
	ansi.Printf("      #: %v\n", s.ParseID())
	ansi.Printf("   when: %v\n", s.when)
	ansi.Printf("   prob: %v\n", s.name)
	ansi.Printf("  alias: %v\n", s.short_name)
	refreshLine(1, *maxWidth)
	ansi.Printf(updateLine(fmt.Sprintf(" status: %v\n", s.ParseStatus()), maxWidth))
	ansi.Printf(" points: %v\n", s.ParsePoints())
}

func display(submissions []Submission, problemID string, first bool, maxWidth *int, line bool) {
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
			sub.short_name,
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
		ansi.Println(line)
	}
}

func getSubmissionID(body []byte) (string, error) {
	reg := regexp.MustCompile(`<tr id="report(\d+?)row">`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("cannot find submission id")
	}
	return string(tmp[1]), nil
}

func findSubmission(body []byte, n int) ([][]byte, error) {
	reg := regexp.MustCompile(`<tr id="report\d+row">[\s\S]+?</tr>`)
	tmp := reg.FindAll(body, n)
	if tmp == nil {
		return nil, errors.New("cannot find any submission")
	}
	return tmp, nil
}

func getProblemNames(name string) (string, string) {
	reg := regexp.MustCompile(`([\s\S]+?) \((\S+?)\)`)
	tmp := reg.FindSubmatch([]byte(name))
	if len(tmp) < 3 {
		return name, ""
	}
	return string(tmp[1]), string(tmp[2])
}

func parseSubmission(body []byte) (ret Submission, err error) {
	reg := regexp.MustCompile(`\d+`)
	toInt := func(sel string) uint64 {
		if tmp := reg.FindString(sel); tmp != "" {
			t, _ := strconv.Atoi(tmp)
			return uint64(t)
		}
		return inf
	}
	id, err := getSubmissionID(body)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf("<table>%v</table>", string(body))))
	if err != nil {
		return
	}
	get := func(sel string) string {
		return strings.TrimSpace(doc.Find(sel).Text())
	}
	when := strings.TrimSpace(doc.Find("a").First().Text())
	combined_name := get(fmt.Sprintf("td#submission%v-problem-instance", id))
	name, short_name := getProblemNames(combined_name)
	points := toInt(get(fmt.Sprintf("td#submission%v-score", id)))
	status := get(fmt.Sprintf("td#submission%v-status", id))
	end := true
	if status == "Oczekuje" || status == "PENDING" {
		end = false
	}
	if strings.Contains(status, "OK") {
		status = fmt.Sprintf("${c-accepted}%v", status)
		if points == inf {
			end = false
		}
	} else if strings.Contains(status, "błąd") {
		status = fmt.Sprintf("${c-failed}%v", status)
		if points == inf {
			end = false
		}
	} else {
		status = fmt.Sprintf("${c-rejected}%v", status)
	}
	return Submission{
		id:         toInt(id),
		name:       name,
		short_name: short_name,
		status:     status,
		points:     points,
		when:       when,
		end:        end,
	}, nil
}

func getProblemName(body []byte) string {
	reg := regexp.MustCompile(`<div class="problem-title text-center content-row">\s+?<h1>([\s\S]+?)</h1>`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return ""
	}
	return string(tmp[1])
}

func (c *SzkopulClient) getSubmissions(URL string, n int) (submissions []Submission, err error) {
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	if _, err = findUsername(body); err != nil {
		return
	}

	submissionsBody, err := findSubmission(body, n)
	if err != nil {
		return
	}

	name, short_name := getProblemNames(getProblemName(body))

	for _, submissionBody := range submissionsBody {
		if submission, err := parseSubmission(submissionBody); err == nil {
			if submission.name == "" && submission.short_name == "" {
				submission.name = name
				submission.short_name = short_name
			}
			submissions = append(submissions, submission)
		}
	}

	if len(submissions) < 1 {
		return nil, errors.New("cannot find any submission")
	}

	return
}

// WatchSubmission n is the number of submissions
func (c *SzkopulClient) WatchSubmission(info Info, n int, line bool) (submissions []Submission, err error) {
	URL := info.MySubmissionURL(c.host)
	if err != nil {
		return
	}

	maxWidth := 0
	first := true
	for {
		st := time.Now()
		submissions, err = c.getSubmissions(URL, n)
		if err != nil {
			return
		}
		display(submissions, info.ProblemID, first, &maxWidth, line)
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
		sub := time.Since(st)
		if sub < time.Second {
			time.Sleep(time.Duration(time.Second - sub))
		}
	}
}

var colorMap = map[string]color.Attribute{
	"${c-waiting}":  color.FgWhite,
	"${c-failed}":   color.FgRed,
	"${c-accepted}": color.FgGreen,
	"${c-partial}":  color.FgCyan,
	"${c-rejected}": color.FgBlue,
}
