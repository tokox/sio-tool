package judge

import (
	"fmt"

	"github.com/fatih/color"
)

type VerdictStatus string

const (
	OK   VerdictStatus = "OK"
	WA   VerdictStatus = "WA"
	TLE  VerdictStatus = "TLE"
	MLE  VerdictStatus = "MLE"
	OLE  VerdictStatus = "OLE"
	RE   VerdictStatus = "RE"
	PERF VerdictStatus = "perf_event_paranoid"
	INT  VerdictStatus = "INT"
)

type Verdict struct {
	Status            VerdictStatus
	TimeInSeconds     float64
	MemoryInMegabytes float64
	Message           string
	Err               error
}

func ParseMemory(memory float64) string {
	if memory >= 1 {
		return fmt.Sprintf("%.3fMB", memory)
	} else if memory*1024.0 >= 1 {
		return fmt.Sprintf("%.3fKB", memory*1024.0)
	}
	return fmt.Sprintf("%.0fB", memory*1024.0*1024.0)
}

func GenerateVerdict(testID, answer string, processInfo ProcessInfo) Verdict {
	state := ""
	diff := ""
	var status VerdictStatus
	output := Plain(processInfo.Output)
	if output == answer {
		status = OK
		state = color.New(color.FgGreen).Sprintf("Passed #%v", testID)
	} else {
		status = WA
		state = color.New(color.FgRed).Sprintf("Failed #%v", testID)
		diff += color.New(color.FgCyan).Sprintf("-----Output-----\n")
		diff += output + "\n"
		diff += color.New(color.FgCyan).Sprintf("-----Answer-----\n")
		diff += answer + "\n"
	}
	return Verdict{status, processInfo.TimeInSeconds, processInfo.MemoryInMegabytes, fmt.Sprintf("%v ... %.3fs %v\n%v", state, processInfo.TimeInSeconds, ParseMemory(processInfo.MemoryInMegabytes), diff), nil}
}
