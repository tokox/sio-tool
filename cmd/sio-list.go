package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/Arapak/sio-tool/util"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

func SioList() (err error) {
	info := Args.SioInfo
	if info.Contest == "" {
		return SioListContests()
	}
	cln := getSioClient()
	err = cln.Ping()
	if err != nil {
		return
	}
	problems, perf, err := cln.Statis(info)
	if err != nil {
		if err = loginAgainSio(cln, err); err == nil {
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
	table.SetHeader([]string{"round", "name", "alias", "points"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, prob := range problems {
		table.Append([]string{
			util.LimitNumOfChars(prob.Round, 20),
			util.LimitNumOfChars(prob.Name, 25),
			prob.Alias,
			prob.ParsePoint(),
		})
	}
	table.Render()

	scanner := bufio.NewScanner(io.Reader(&buf))
	for i := -2; scanner.Scan(); i++ {
		line := scanner.Text()
		_, _ = ansi.Println(line)
	}
	return
}

func SioListContests() (err error) {
	cln := getSioClient()
	err = cln.Ping()
	if err != nil {
		return
	}
	contests, perf, err := cln.ListContests()
	if err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			contests, perf, err = cln.ListContests()
		}
	}
	if err != nil {
		return
	}
	fmt.Printf("Statis: (%v)\n", perf.Parse())

	var buf bytes.Buffer
	output := io.Writer(&buf)
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"contests"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.SetAutoMergeCells(true)
	for _, prob := range contests {
		if prob.Subheader {
			table.Append([]string{util.GreenString(prob.Name)})
		} else {
			table.Append([]string{fmt.Sprintf("%v (%v)", util.LimitNumOfChars(prob.Name, 20), prob.Alias)})
		}
	}
	table.Render()

	scanner := bufio.NewScanner(io.Reader(&buf))
	for i := -2; scanner.Scan(); i++ {
		line := scanner.Text()
		_, _ = ansi.Println(line)
	}
	return
}
