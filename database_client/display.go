package database_client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/Arapak/sio-tool/util"

	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
)

func (t *Task) ParseID() string {
	return fmt.Sprintf("%v", t.ID)
}

const maxNameLength = 22
const maxAliasLength = 8
const maxStageLength = 5

func (t *Task) Display() {
	_, _ = ansi.Printf("       #: %v\n", t.ParseID())
	_, _ = ansi.Printf("    name: %v\n", t.Name)
	_, _ = ansi.Printf("  source: %v\n", t.Source)
	_, _ = ansi.Printf("    path: %v\n", t.Path)
	_, _ = ansi.Printf("   alias: %v\n", t.ShortName)
	_, _ = ansi.Printf("    link: %v\n", t.Link)
	_, _ = ansi.Printf(" contest: %v\n", t.ContestID)
	_, _ = ansi.Printf("   stage: %v\n", t.ContestStageID)
}

func Display(tasks []Task) {
	var buf bytes.Buffer
	output := io.Writer(&buf)
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"#", "name", "source", "alias", "contest", "stage"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, t := range tasks {
		table.Append([]string{
			util.GreenString(t.ParseID()),
			util.LimitNumOfChars(t.Name, maxNameLength),
			t.Source,
			util.LimitNumOfChars(t.ShortName, maxAliasLength),
			t.ContestID,
			util.LimitNumOfChars(t.ContestStageID, maxStageLength),
		})
	}
	table.Render()
	scanner := bufio.NewScanner(io.Reader(&buf))
	for scanner.Scan() {
		line := scanner.Text()
		_, _ = ansi.Println(line)
	}
}
