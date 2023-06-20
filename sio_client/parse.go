package sio_client

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Arapak/sio-tool/database_client"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/Arapak/sio-tool/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

const ErrorParsingSamples = `parsing samples failed`
const ErrorServiceUnavailable = `service unavailable`
const ErrorFetchingProblemSiteFailed = `fetching problem site failed`
const ErrorParsingProblemsFailed = `parsing some problems failed`
const ErrorUnrecognizedStatementFormat = `unrecognized statement format`

const StandardIOReg = `(\nKomunikacja\n|\nOpis interfejsu\s+)`

func parseSiteStatement(body []byte) (name string, standardIO bool, input [][]byte, output [][]byte, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	name = doc.Find("h1").First().Text()
	reg := regexp.MustCompile(StandardIOReg)
	standardIO = !reg.Match(body)
	statement := doc.Find("section.main-content").First().Text()
	if standardIO {
		input, output, err = findSamples([]byte(statement))
	}
	return
}

func getNameFromOiPdf(statement []byte) (problemName string) {
	reg := regexp.MustCompile(`^Zadanie: (?P<alias>\S+?)\n(?P<problemName>[\S ]+?)\n`)
	names := reg.SubexpNames()
	for i, val := range reg.FindSubmatch(statement) {
		if names[i] == "problemName" {
			problemName = string(val)
		}
	}
	return
}

func getNameFromSinolPdf(statement []byte) (problemName string) {
	reg := regexp.MustCompile(`\A([\S ]*?\n+)?([\S ]*?\n+)?Dostępna pamięć: \d+MB\n+(?P<problemName>[\S ]+?)\n`)
	names := reg.SubexpNames()
	for i, val := range reg.FindSubmatch(statement) {
		if names[i] == "problemName" {
			problemName = string(val)
		}
	}
	return
}

func getNameFromPdf(statement []byte) (problemName string) {
	problemName = getNameFromSinolPdf(statement)
	if problemName != "" {
		return
	}
	return getNameFromOiPdf(statement)
}

func parsePdf(body []byte) (name string, standardIO bool, input [][]byte, output [][]byte, err error) {
	statement, err := pdfToText(body)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(StandardIOReg)
	standardIO = !reg.Match(statement)
	name = getNameFromPdf(statement)
	if name == "" {
		err = errors.New("parsing problem failed")
		return
	}
	if standardIO {
		input, output, err = findSamples(statement)
	}
	return
}

func pdfToText(body []byte) ([]byte, error) {
	cmd := exec.Command("pdftotext", "-", "-")
	cmd.Stdin = bytes.NewReader(body)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

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

func noWhitespaceOnly(input [][]byte) bool {
	for _, x := range input {
		if isWhitespaceOnly(x) {
			return false
		}
	}
	return true
}

func findSamples(body []byte) (input, output [][]byte, err error) {
	sectionStart1Reg := regexp.MustCompile(`Przykład(y)?\s+`)
	sectionEnd1Reg := regexp.MustCompile(`Wyjaśnienie przykładu|Wyjaśnienie do przykładu|Komentarz do przykładu|Ocenianie|Autor(zy)? zadania:|Testy „ocen”`)

	startIndex := sectionStart1Reg.FindIndex(body)
	if startIndex != nil {
		body = body[startIndex[1]:]
	}
	endIndex := sectionEnd1Reg.FindIndex(body)
	if endIndex != nil {
		body = body[:endIndex[0]]
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
		`[Pp]oprawną odpowiedzią jest:\s*`,
		`[Pp]rzykładowe wyjście\s*`,
	}
	const sioSampleInputStart = `Wejście\s*`
	const sioSampleOutputStart = `Wyjście\s*`

	const sampleEnd = `\n\n|\z`

	inputEnd := createRegGroupFromArray(append(oiSampleInputStarts, append(oiSampleOutputStarts, sampleEnd)...))
	oiSampleInputStartReg := createRegGroupFromArray(oiSampleInputStarts)
	oiSampleOutputStartReg := createRegGroupFromArray(oiSampleOutputStarts)

	input = findSubstring(body, oiSampleInputStartReg, inputEnd, 2)
	output = findSubstring(body, oiSampleOutputStartReg, inputEnd, 2)

	if !noWhitespaceOnly(input) || !noWhitespaceOnly(output) {
		input = findSubstring(body, oiSampleInputStartReg+inputEnd+oiSampleOutputStartReg, inputEnd, 4)
		output = findSubstring(body, oiSampleInputStartReg+inputEnd+oiSampleOutputStartReg+`([\s\S]*?)`+inputEnd, inputEnd, 6)
	}

	if !noWhitespaceOnly(input) || !noWhitespaceOnly(output) {
		return nil, nil, errors.New(ErrorParsingSamples)
	}

	if len(input) == 0 && len(output) == 0 {
		input = findSubstring(body, sioSampleInputStart, "("+sampleEnd+")", 1)
		output = findSubstring(body, sioSampleOutputStart, "("+sampleEnd+")", 1)
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

func addNewLine(body []byte) []byte {
	if !bytes.HasSuffix(body, []byte("\n")) {
		return append(body, byte('\n'))
	}
	return body
}

func (c *SioClient) ParseProblem(host, contestID, problemAlias, path string, mu *sync.Mutex) (name string, samples int, standardIO bool, perf util.Performance, err error) {
	perf.StartFetching()

	resp, err := c.client.Get(ProblemURL(host, contestID, problemAlias))
	if err != nil {
		return
	}

	perf.StopFetching()
	perf.StartParsing()

	var input, output [][]byte
	var body []byte
	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if resp.Header.Get("Content-Type") == "application/pdf" {
			name, standardIO, input, output, err = parsePdf(body)
		} else if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			name, standardIO, input, output, err = parseSiteStatement(body)
		} else {
			err = errors.New(ErrorUnrecognizedStatementFormat)
		}
	} else {
		err = errors.New(ErrorFetchingProblemSiteFailed)
	}
	if err != nil {
		return
	}

	perf.StopParsing()

	samples = len(input)
	for i := 0; i < samples; i++ {
		fileIn := filepath.Join(path, fmt.Sprintf("in%v.txt", i+1))
		fileOut := filepath.Join(path, fmt.Sprintf("out%v.txt", i+1))
		input[i] = addNewLine(input[i])
		e := os.WriteFile(fileIn, input[i], 0644)
		if e != nil {
			if mu != nil {
				mu.Lock()
			}
			color.Red(e.Error())
			if mu != nil {
				mu.Unlock()
			}
		}
		output[i] = addNewLine(output[i])
		e = os.WriteFile(fileOut, output[i], 0644)
		if e != nil {
			if mu != nil {
				mu.Lock()
			}
			color.Red(e.Error())
			if mu != nil {
				mu.Unlock()
			}
		}
	}
	return
}

func (c *SioClient) parse(contestID, problemAlias, path string, mu *sync.Mutex) (perf util.Performance, err error) {
	name, samples, standardIO, perf, err := c.ParseProblem(c.host, contestID, problemAlias, path, mu)

	warns := ""
	if !standardIO {
		warns = color.YellowString("Non standard input output format.")
	} else if err != nil && err.Error() == ErrorParsingSamples {
		warns = color.RedString("Error parsing samples")
		err = nil
	}

	if mu != nil {
		mu.Lock()
	}
	if err != nil {
		color.Red("Failed (%v). Error: %v", problemAlias, err.Error())
	} else {
		ansi.Printf("%v %v\n", color.GreenString("Parsed %v (%v) with %v samples.", name, problemAlias, samples), warns)
	}
	if mu != nil {
		mu.Unlock()
	}
	return
}

func (c *SioClient) Parse(info Info, db *sql.DB) (problems []StatisInfo, paths []string, err error) {
	start := time.Now()

	color.Cyan("Parse " + info.Hint())
	problems, statisPerf, err := c.Statis(info)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("Statis: (%v)\n", statisPerf.Parse())

	if info.Round == "" && info.ProblemAlias == "" {
		var rounds []string
		for _, problem := range problems {
			if len(rounds) == 0 || problem.Round != rounds[len(rounds)-1] {
				rounds = append(rounds, problem.Round)
			}
		}
		if len(rounds) > 1 {
			color.Cyan("Which round do you want to parse?")
			for i, round := range rounds {
				fmt.Printf("(%v) %v\n", i, round)
			}
			fmt.Printf("(%v) all\n", len(rounds))
			index := util.ChooseIndex(len(rounds) + 1)
			if index < len(rounds) {
				info.Round = rounds[index]
				var filteredProblems []StatisInfo
				for _, problem := range problems {
					if problem.Round == info.Round {
						filteredProblems = append(filteredProblems, problem)
					}
				}
				problems = filteredProblems
			}
		}
	}

	if len(problems) >= 10 && !util.Confirm(fmt.Sprintf("Are you sure you want to parse %v problems? (Y/n): ", len(problems))) {
		return
	}

	for _, problem := range problems {
		pathInfo := Info{RootPath: info.RootPath, ProblemAlias: problem.Alias, Round: problem.Round, Contest: info.Contest}
		paths = append(paths, pathInfo.Path())
	}
	contestPath := info.Path()
	ansi.Printf(color.CyanString("The problem(s) will be saved to %v\n"), color.GreenString(contestPath))

	var avgPerformance util.Performance

	parsed := 0

	const numberOfWorkers = 50
	index := 0

	wg := sync.WaitGroup{}
	wg.Add(numberOfWorkers)
	mu := sync.Mutex{}
	for i := 1; i <= numberOfWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				mu.Lock()
				if index >= len(problems) {
					mu.Unlock()
					break
				}
				problemIndex := index
				index++
				fmt.Printf("Parsing %v (%v)\n", problems[problemIndex].Name, problems[problemIndex].Alias)
				path := paths[problemIndex]
				problemAlias := problems[problemIndex].Alias
				contestID := info.Contest
				mu.Unlock()

				err = os.MkdirAll(path, os.ModePerm)
				if err != nil {
					continue
				}

				var perf util.Performance
				perf, err = c.parse(contestID, problemAlias, path, &mu)
				mu.Lock()
				avgPerformance.Fetching += perf.Fetching
				avgPerformance.Parsing += perf.Parsing
				if err == nil {
					task := database_client.Task{
						Name:      problems[problemIndex].Name,
						Source:    "sio",
						Path:      path,
						ShortName: problems[problemIndex].Alias,
						Link:      ProblemURL(c.host, info.Contest, problemAlias),
						ContestID: info.Contest,
					}
					err := database_client.AddTask(db, task)
					if err != nil {
						color.Red(err.Error())
					}
					parsed++
				}
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	color.Green("Successfully parsed %v problems.\n", parsed)
	avgPerformance.Fetching = util.AverageTime(avgPerformance.Fetching, len(problems))
	avgPerformance.Parsing = util.AverageTime(avgPerformance.Parsing, len(problems))
	fmt.Printf("Average: (%v)\n", avgPerformance.Parse())
	fmt.Printf("Total: %s\n", time.Since(start).Round(time.Millisecond))
	if parsed != len(problems) {
		err = errors.New(ErrorParsingProblemsFailed)
	}
	return
}
