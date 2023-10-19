package judge

import (
	"os"
	"strings"
)

func Judge(inPath, ansPath, sampleID, command string, oiejqOptions *OiejqOptions) Verdict {
	input, err := os.Open(inPath)
	if err != nil {
		return Verdict{false, "", err}
	}
	defer input.Close()

	var processInfo *ProcessInfo
	if oiejqOptions != nil {
		processInfo, err = RunProcessWithOiejq(sampleID, command, input, oiejqOptions)
	} else {
		processInfo, err = RunProcess(sampleID, command, input, nil)
	}
	if err != nil {
		return Verdict{false, "", err}
	}

	b, err := os.ReadFile(ansPath)
	if err != nil {
		return Verdict{false, "", err}
	}
	return GenerateVerdict(sampleID, Plain(b), *processInfo)
}

func ExtractTaskName(file string) (task string) {
	task, _, _ = strings.Cut(file, "-")
	return
}
