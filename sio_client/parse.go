package sio_client

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Arapak/sio-tool/database_client"
	"github.com/Arapak/sio-tool/sio_samples"
	"github.com/Arapak/sio-tool/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

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
		input, output, err = sio_samples.FindSamples([]byte(statement))
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
	statement, err := util.PdfToText(body)
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
		input, output, err = sio_samples.FindSamples(statement)
	}
	return
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
		input[i] = util.AddNewLine(input[i])
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
		output[i] = util.AddNewLine(output[i])
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
	} else if err != nil && err.Error() == sio_samples.ErrorParsingSamples {
		warns = color.RedString("Error parsing samples")
		err = nil
	}

	if mu != nil {
		mu.Lock()
	}
	if err != nil {
		color.Red("Failed (%v). Error: %v", problemAlias, err.Error())
	} else {
		_, _ = ansi.Printf("%v %v\n", color.GreenString("Parsed %v (%v) with %v samples.", name, problemAlias, samples), warns)
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
			prompt := &survey.Select{
				Message: "Which round do you want to parse?",
				Options: append(rounds, "ALL"),
			}
			if err = survey.AskOne(prompt, &info.Round); err != nil {
				return
			}

			if info.Round != "ALL" {
				var filteredProblems []StatisInfo
				for _, problem := range problems {
					if problem.Round == info.Round {
						filteredProblems = append(filteredProblems, problem)
					}
				}
				problems = filteredProblems
			} else {
				info.Round = ""
			}
		}
	}

	if len(problems) == 0 {
		color.Red("no problems to parse")
		return
	}

	if len(problems) >= 10 {
		parseAll := true
		prompt := &survey.Confirm{Message: fmt.Sprintf("Are you sure you want to parse %v problems?", len(problems)), Default: true}
		if err = survey.AskOne(prompt, &parseAll); err != nil {
			return
		}
		if !parseAll {
			return
		}
	}

	for _, problem := range problems {
		pathInfo := Info{RootPath: info.RootPath, ProblemAlias: problem.Alias, Round: problem.Round, Contest: info.Contest}
		paths = append(paths, pathInfo.Path())
	}
	contestPath := info.Path()
	_, _ = ansi.Printf(color.CyanString("The problem(s) will be saved to %v\n"), color.GreenString(contestPath))

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
