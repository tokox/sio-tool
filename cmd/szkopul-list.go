package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/Arapak/sio-tool/szkopul_client"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

func SzkopulList() (err error) {
	cln := szkopul_client.Instance
	info := Args.SzkopulInfo
	problems, perf, err := cln.Statis(info)
	if err != nil {
		if err = loginAgainSzkopul(cln, err); err == nil {
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
	table.SetHeader([]string{"contest", "stage", "name", "alias", "points"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, prob := range problems {
		table.Append([]string{
			prob.Contest,
			prob.Stage,
			prob.Name,
			prob.Alias,
			prob.ParsePoint(),
		})
	}
	table.Render()

	scanner := bufio.NewScanner(io.Reader(&buf))
	for i := -2; scanner.Scan(); i++ {
		line := scanner.Text()
		ansi.Println(line)
	}
	return
}
