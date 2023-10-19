package sio_submissions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

type Submission struct {
	Name      string
	ShortName string
	Id        uint64
	Status    string
	Points    uint64
	When      string
	End       bool
}

func (s *Submission) ParseStatus() string {
	status := s.Status
	for k, v := range colorMap {
		tmp := strings.ReplaceAll(status, k, "")
		if tmp != status {
			status = color.New(v).Sprint(tmp)
		}
	}
	return status
}

func (s *Submission) ParseID() string {
	return fmt.Sprintf("%v", s.Id)
}

const Inf = 1000000009

var intReg = regexp.MustCompile(`\d+`)

func ToInt(sel string) uint64 {
	if tmp := intReg.FindString(sel); tmp != "" {
		t, _ := strconv.Atoi(tmp)
		return uint64(t)
	}
	return Inf
}

func (s *Submission) ParsePoints() string {
	if s.Points == Inf {
		return ""
	} else if s.Points == 0 {
		return color.New(colorMap["${c-failed}"]).Sprint(s.Points)
	} else if s.Points < 100 {
		return color.New(colorMap["${c-partial}"]).Sprint(s.Points)
	} else if s.Points == 100 {
		return color.New(colorMap["${c-accepted}"]).Sprint(s.Points)
	}
	return fmt.Sprintf("%v", s.Points)
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
	_, _ = ansi.Printf("   when: %v\n", s.When)
	_, _ = ansi.Printf("   prob: %v\n", s.Name)
	_, _ = ansi.Printf("  alias: %v\n", s.ShortName)
	refreshLine(1, *maxWidth)
	_, _ = ansi.Printf(updateLine(fmt.Sprintf(" status: %v\n", s.ParseStatus()), maxWidth))
	_, _ = ansi.Printf(" points: %v\n", s.ParsePoints())
}

func Display(submissions []Submission, first bool, maxWidth *int, line bool) {
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
			sub.When,
			sub.Name,
			sub.ShortName,
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

var colorMap = map[string]color.Attribute{
	"${c-waiting}":  color.FgBlue,
	"${c-failed}":   color.FgRed,
	"${c-accepted}": color.FgGreen,
	"${c-partial}":  color.FgCyan,
	"${c-rejected}": color.FgBlue,
}
