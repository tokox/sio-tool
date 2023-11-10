package sio_samples

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/Arapak/sio-tool/util"
)

const ErrorParsingSamples = `parsing samples failed`

func createRegGroupFromArray(regs []string) (reg string) {
	reg = fmt.Sprintf(`(%v`, regs[0])
	for _, r := range regs {
		reg += fmt.Sprintf("|%v", r)
	}
	reg += ")"
	return
}

func isWhitespaceOnly(input []byte) bool {
	for _, b := range input {
		if !unicode.IsSpace(rune(b)) {
			return false
		}
	}
	return true
}

func containsWhitespaceOnly(input [][]byte) bool {
	for _, x := range input {
		if isWhitespaceOnly(x) {
			return true
		}
	}
	return false
}

func findSamplesSection(body []byte) []byte {
	sectionStartReg := regexp.MustCompile(`Przykład(y)?\s+`)
	sectionEndReg := regexp.MustCompile(`Wyjaśnienie przykładu|Wyjaśnienie do przykładu|Komentarz do przykładu|Ocenianie|Autor(zy)? zadania:|Testy „ocen”|Testy przykładowe`)

	startIndex := sectionStartReg.FindIndex(body)
	if startIndex != nil {
		body = body[startIndex[1]:]
	}
	endIndex := sectionEndReg.FindIndex(body)
	if endIndex != nil {
		body = body[:endIndex[0]]
	}
	return body
}

func maxSpace(s string) int {
	maxSpc := 0
	index := len(s)
	currentSpace := 0
	for i := range s {
		if s[i] == ' ' {
			currentSpace++
			if currentSpace > maxSpc {
				maxSpc = currentSpace
				index = i
			}
		} else {
			currentSpace = 0
		}
	}
	return index
}

func trimSpace(s string) string {
	s = strings.TrimSpace(s)
	whitespaces := regexp.MustCompile(`\s+`)
	return whitespaces.ReplaceAllString(s, " ")
}

func parseSinolSample(sample []byte) (input string, output string, err error) {
	lines := strings.Split(string(sample), "\n")
	index := maxSpace(lines[0])
	for i := range lines {
		if len(lines[i]) >= index {
			input += trimSpace(lines[i][:index]) + "\n"
			output += trimSpace(lines[i][index:]) + "\n"
		} else {
			input += trimSpace(lines[i]) + "\n"
		}
	}
	input = strings.TrimSpace(input)
	output = strings.TrimSpace(output)
	return
}

const sioSampleInputStart = `Wejście\s*`
const sioSampleOutputStart = `Wyjście\s*`

const sampleEnd = `\n\n|\z`

func parseSinolStatement(pdf []byte) (input, output [][]byte, err error) {
	body, err := util.PdfToTextLayout(pdf)
	if err != nil {
		return
	}
	body = findSamplesSection(body)
	merged := findSubstring(body, sioSampleInputStart+sioSampleOutputStart, "("+sampleEnd+")", 1)
	for _, sample := range merged {
		inputString, outputString, err := parseSinolSample(sample)
		if err != nil {
			return input, output, err
		}
		input = append(input, []byte(inputString))
		output = append(output, []byte(outputString))
	}
	return
}

var oiSampleInputStarts = []string{
	`[Dd]la danych wejściowych:\s*`,
	`[AIai] dla danych wejściowych:\s*`,
	`[Aa] dla danych:\s*`,
	`[Nn]atomiast dla danych wejściowych:\s*`,
	`[Nn]atomiast dla danych:\s*`,
	`[Zz] kolei dla danych wejściowych:\s*`,
	`[Pp]rzykładowe wejście\s*`,
}

var oiSampleOutputStarts = []string{
	`[Pp]oprawnym wynikiem jest:\s*`,
	`[Jj]ednym z poprawnych wyników jest:\s*`,
	`[Mm]ożliwym poprawnym wynikiem jest:\s*`,
	`[Mm]ożliwym wynikiem jest:\s*`,
	`[Mm]ożliwą odpowiedzią jest:\s*`,
	`[Pp]oprawną odpowiedzią jest:\s*`,
	`[Pp]rzykładowe wyjście\s*`,
}

func FindSamples(body []byte, pdf []byte) (input, output [][]byte, err error) {
	body = findSamplesSection(body)

	inputEnd := createRegGroupFromArray(append(oiSampleInputStarts, append(oiSampleOutputStarts, sampleEnd)...))
	oiSampleInputStartReg := createRegGroupFromArray(oiSampleInputStarts)
	oiSampleOutputStartReg := createRegGroupFromArray(oiSampleOutputStarts)

	input = findSubstring(body, oiSampleInputStartReg, inputEnd, 2)
	output = findSubstring(body, oiSampleOutputStartReg, inputEnd, 2)

	if containsWhitespaceOnly(input) || containsWhitespaceOnly(output) {
		input = findSubstring(body, oiSampleInputStartReg+inputEnd+oiSampleOutputStartReg, inputEnd, 4)
		output = findSubstring(body, oiSampleInputStartReg+inputEnd+oiSampleOutputStartReg+`([\s\S]*?)`+inputEnd, inputEnd, 6)
	}

	if containsWhitespaceOnly(input) || containsWhitespaceOnly(output) {
		return nil, nil, errors.New(ErrorParsingSamples)
	}

	if len(input) == 0 && len(output) == 0 {
		return parseSinolStatement(pdf)

	}
	if len(input) != len(output) {
		return nil, nil, errors.New(ErrorParsingSamples)
	}
	regs := make([]*regexp.Regexp, len(oiSampleInputStarts)+len(oiSampleOutputStarts))
	for i, oiSampleInputStart := range oiSampleInputStarts {
		regs[i] = regexp.MustCompile(oiSampleInputStart)
	}
	for i, oiSampleOutputStart := range oiSampleOutputStarts {
		regs[len(oiSampleInputStarts)+i] = regexp.MustCompile(oiSampleOutputStart)
	}
	for _, inp := range input {
		for _, reg := range regs {
			if reg.Match(inp) {
				return nil, nil, errors.New(ErrorParsingSamples)
			}
		}
	}
	for _, out := range output {
		for _, reg := range regs {
			if reg.Match(out) {
				return nil, nil, errors.New(ErrorParsingSamples)
			}
		}
	}
	return
}

func findSubstring(body []byte, before, after string, group int) (result [][]byte) {
	re := regexp.MustCompile(before + `([\s\S]*?)` + after)
	matches := re.FindAllSubmatch(body, -1)
	for _, match := range matches {
		if len(match) > group {
			result = append(result, match[group])
		}
	}
	return
}
