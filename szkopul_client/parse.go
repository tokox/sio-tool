package szkopul_client

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

const ErrorServiceUnavailable = `service unavailable`
const ErrorFetchingProblemSiteFailed = `fetching problem site failed`
const ErrorParsingProblemsFailed = `parsing some problems failed`

const SiteStatementProblemURL = `/problemset/problem/%v/site/?key=statement`
const PdfStatementProblemURL = `/problemset/problem/%v/statement`

const StandardIOReg = `(\nKomunikacja\n|\nOpis interfejsu\s+)`

func parseSiteStatement(body []byte) (name string, alias string, standardIO bool, input [][]byte, output [][]byte, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	alias, name = GetAliasAndName(doc.Find("h1").First().Text())
	reg := regexp.MustCompile(StandardIOReg)
	standardIO = !reg.Match(body)
	statement := doc.Find(".nav-content").First().Text()
	if standardIO {
		input, output, err = sio_samples.FindSamples([]byte(statement))
	}
	return
}

func getNamesFromPdf(statement []byte) (problemName string, alias string) {
	reg := regexp.MustCompile(`^Zadanie: (?P<alias>\S+?)\n(?P<problemName>[\S ]+?)\n`)
	names := reg.SubexpNames()
	for i, val := range reg.FindSubmatch(statement) {
		if names[i] == "problemName" {
			problemName = string(val)
		} else if names[i] == "alias" {
			alias = strings.ToLower(string(val))
		}
	}
	return
}

func parsePdf(body []byte) (name string, alias string, standardIO bool, input [][]byte, output [][]byte, err error) {
	statement, err := util.PdfToText(body)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(StandardIOReg)
	standardIO = !reg.Match(statement)
	name, alias = getNamesFromPdf(statement)
	if name == "" || alias == "" {
		err = errors.New("parsing problem failed")
		return
	}
	if standardIO {
		input, output, err = sio_samples.FindSamples(statement)
	}
	return
}

func (c *SzkopulClient) ParseProblem(host, problemID, path string, mu *sync.Mutex) (name string, alias string, samples int, standardIO bool, perf util.Performance, err error) {
	perf.StartFetching()

	resp, err := c.client.Get(fmt.Sprintf(host+PdfStatementProblemURL, problemID))
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
		name, alias, standardIO, input, output, err = parsePdf(body)
	} else if resp.StatusCode == 403 {
		resp, err = c.client.Get(fmt.Sprintf(host+SiteStatementProblemURL, problemID))
		if err != nil {
			return
		}
		if resp.StatusCode == 503 {
			err = errors.New(ErrorServiceUnavailable)
			return
		} else if resp.StatusCode != 200 {
			err = errors.New(ErrorFetchingProblemSiteFailed)
			return
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err == nil {
			name, alias, standardIO, input, output, err = parseSiteStatement(body)
		}
	} else if resp.StatusCode == 503 {
		err = errors.New(ErrorServiceUnavailable)
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

func (c *SzkopulClient) parse(problemID, path string, mu *sync.Mutex) (perf util.Performance, err error) {
	name, alias, samples, standardIO, perf, err := c.ParseProblem(c.host, problemID, path, mu)

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
		color.Red("Failed (%v). Error: %v", problemID, err.Error())
	} else {
		_, _ = ansi.Printf("%v %v\n", color.GreenString("Parsed %v (%v) with %v samples.", name, alias, samples), warns)
	}
	if mu != nil {
		mu.Unlock()
	}
	return
}

func (c *SzkopulClient) Parse(info Info, db *sql.DB) (problems []StatisInfo, paths []string, err error) {
	start := time.Now()

	color.Cyan("Parse " + info.Hint())
	problems, statisPerf, err := c.Statis(info)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("Statis: (%v)\n", statisPerf.Parse())

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
		pathInfo := Info{RootPath: info.RootPath, ProblemAlias: problem.Alias, ContestID: problem.Contest, StageID: problem.Stage}
		paths = append(paths, pathInfo.Path())
	}
	contestPath := info.Path()
	_, _ = ansi.Printf(color.CyanString("The problem(s) will be saved to %v\n"), color.GreenString(contestPath))

	var retry []int
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
				problemID := problems[problemIndex].ID
				mu.Unlock()

				err = os.MkdirAll(path, os.ModePerm)
				if err != nil {
					continue
				}

				var perf util.Performance
				perf, err = c.parse(problemID, path, &mu)
				if err != nil && err.Error() == ErrorServiceUnavailable {
					mu.Lock()
					retry = append(retry, problemIndex)
					mu.Unlock()
					continue
				}
				mu.Lock()
				avgPerformance.Fetching += perf.Fetching
				avgPerformance.Parsing += perf.Parsing
				if err == nil {
					task := database_client.Task{
						Name:           problems[problemIndex].Name,
						Source:         info.Archive,
						Path:           path,
						ShortName:      problems[problemIndex].Alias,
						Link:           ProblemURL(c.host, problemID),
						ContestID:      problems[problemIndex].Contest,
						ContestStageID: problems[problemIndex].Stage,
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
	index = 0
	for _, index := range retry {
		color.Blue("Retrying %v (%v)\n", problems[index].Name, problems[index].Alias)

		err = os.MkdirAll(paths[index], os.ModePerm)
		if err != nil {
			return
		}

		var perf util.Performance
		perf, err = c.parse(problems[index].ID, paths[index], nil)

		avgPerformance.Fetching += perf.Fetching
		avgPerformance.Parsing += perf.Parsing
		if err == nil {
			task := database_client.Task{
				Name:           problems[index].Name,
				Source:         info.Archive,
				Path:           paths[index],
				ShortName:      problems[index].Alias,
				Link:           ProblemURL(c.host, problems[index].ID),
				ContestID:      problems[index].Contest,
				ContestStageID: problems[index].Stage,
			}
			err := database_client.AddTask(db, task)
			if err != nil {
				color.Red(err.Error())
			}
			parsed++
		}
	}
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
