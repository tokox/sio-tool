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

func getPrefix(str string, length int) string {
	i := 0
	for j := range str {
		if i == length {
			return str[:j]
		}
		i++
	}
	return str
}

func (t *Task) ParseName() string {
	return getPrefix(t.Name, maxNameLength)
}

func (t *Task) ParseAlias() string {
	return getPrefix(t.ShortName, maxAliasLength)
}

func (t *Task) ParseStage() string {
	return getPrefix(t.ContestStageID, maxStageLength)
}

func (t *Task) Display() {
	ansi.Printf("       #: %v\n", t.ParseID())
	ansi.Printf("    name: %v\n", t.Name)
	ansi.Printf("  source: %v\n", t.Source)
	ansi.Printf("    path: %v\n", t.Path)
	ansi.Printf("   alias: %v\n", t.ShortName)
	ansi.Printf("    link: %v\n", t.Link)
	ansi.Printf(" contest: %v\n", t.ContestID)
	ansi.Printf("   stage: %v\n", t.ContestStageID)
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
			t.ParseName(),
			t.Source,
			t.ParseAlias(),
			t.ContestID,
			t.ParseStage(),
		})
	}
	table.Render()
	scanner := bufio.NewScanner(io.Reader(&buf))
	for scanner.Scan() {
		line := scanner.Text()
		ansi.Println(line)
	}
}
