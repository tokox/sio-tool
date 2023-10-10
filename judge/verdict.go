package judge

import (
	"fmt"

	"github.com/fatih/color"
)

type Verdict struct {
	Correct bool
	Message string
	Err     error
}

func GenerateVerdict(testID, answer string, processInfo ProcessInfo) Verdict {
	state := ""
	diff := ""
	var correct bool
	output := Plain(processInfo.Output)
	if output == answer {
		correct = true
		state = color.New(color.FgGreen).Sprintf("Passed #%v", testID)
	} else {
		correct = false
		state = color.New(color.FgRed).Sprintf("Failed #%v", testID)
		diff += color.New(color.FgCyan).Sprintf("-----Output-----\n")
		diff += output + "\n"
		diff += color.New(color.FgCyan).Sprintf("-----Answer-----\n")
		diff += answer + "\n"
	}
	return Verdict{correct, fmt.Sprintf("%v ... %.3fs %v\n%v", state, processInfo.Time, processInfo.Memory, diff), nil}
}
