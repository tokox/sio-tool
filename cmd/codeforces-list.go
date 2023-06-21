package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Arapak/sio-tool/codeforces_client"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

func CodeforcesList() (err error) {
	cln := codeforces_client.Instance
	info := Args.CodeforcesInfo
	problems, perf, err := cln.Statis(info)
	if err != nil {
		if err = loginAgainCodeforces(cln, err); err == nil {
			problems, perf, err = cln.Statis(info)
		}
	}
	if err != nil {
		return
	}
	fmt.Printf("Statis: (%v)\n", perf.Parse())
	var buf bytes.Buffer
	output := io.Writer(&buf)
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"#", "problem", "passed", "limit", "IO"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, prob := range problems {
		table.Append([]string{
			prob.ID,
			prob.Name,
			prob.Passed,
			prob.Limit,
			prob.IO,
		})
	}
	table.Render()

	scanner := bufio.NewScanner(io.Reader(&buf))
	for i := -2; scanner.Scan(); i++ {
		line := scanner.Text()
		if i >= 0 {
			if strings.Contains(problems[i].State, "accepted") {
				line = color.New(color.BgGreen).Sprint(line)
			} else if strings.Contains(problems[i].State, "rejected") {
				line = color.New(color.BgRed).Sprint(line)
			}
		}
		_, err := ansi.Println(line)
		if err != nil {
			return err
		}
	}
	return
}
