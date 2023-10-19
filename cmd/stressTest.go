package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/judge"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
)

func StressTest() (err error) {
	cfg := config.Instance
	if len(cfg.Template) == 0 {
		return errors.New("you have to add at least one code template by `st config`")
	}
	if len(cfg.DefaultNaming) == 0 || cfg.DefaultNaming["solve"] == "" || cfg.DefaultNaming["brute"] == "" || cfg.DefaultNaming["gen"] == "" || cfg.DefaultNaming["test_in"] == "" {
		return errors.New("you have to add default naming by `st config`")
	}

	task := Args.Specifier[0]

	solveFilePattern := cfg.DefaultNaming["solve"]
	if Args.Solve != "" {
		solveFilePattern = Args.Solve
	} else {
		solveFilePattern = strings.ReplaceAll(solveFilePattern, "$%task%$", task)
	}
	solveFilename, index, err := getOneCode(solveFilePattern, cfg.Template, map[string]struct{}{})
	if err != nil {
		return
	}
	solvePath, solveFull := filepath.Split(solveFilename)
	ext := filepath.Ext(solveFilename)
	solveFile := solveFull[:len(solveFull)-len(ext)]

	bruteFilePattern := cfg.DefaultNaming["brute"]
	if Args.Brute != "" {
		bruteFilePattern = Args.Brute
	} else {
		bruteFilePattern = strings.ReplaceAll(bruteFilePattern, "$%task%$", task)
	}
	bruteFilename, _, err := getOneCode(bruteFilePattern, cfg.Template, map[string]struct{}{})
	if err != nil {
		return
	}
	brutePath, bruteFull := filepath.Split(bruteFilename)
	ext = filepath.Ext(bruteFilename)
	bruteFile := bruteFull[:len(bruteFull)-len(ext)]

	testsGenFilePattern := cfg.DefaultNaming["gen"]
	if Args.Generator != "" {
		testsGenFilePattern = Args.Generator
	} else {
		testsGenFilePattern = strings.ReplaceAll(testsGenFilePattern, "$%task%$", task)
	}
	testsGenFilename, _, err := getOneCode(testsGenFilePattern, cfg.Template, map[string]struct{}{})
	if err != nil {
		return
	}
	testsGenPath, testsGenFull := filepath.Split(testsGenFilename)
	ext = filepath.Ext(testsGenFilename)
	testsGenFile := testsGenFull[:len(testsGenFull)-len(ext)]

	template := cfg.Template[index]
	rand := util.RandString(8)

	filter := func(cmd, path, full, file string) string {
		cmd = strings.ReplaceAll(cmd, "$%rand%$", rand)
		cmd = strings.ReplaceAll(cmd, "$%path%$", path)
		cmd = strings.ReplaceAll(cmd, "$%full%$", full)
		cmd = strings.ReplaceAll(cmd, "$%file%$", file)
		cmd = strings.ReplaceAll(cmd, "$%task%$", task)
		return cmd
	}

	run := func(script, path, full, file string) error {
		if s := filter(script, path, full, file); len(s) > 0 {
			fmt.Println(s)
			cmds := util.SplitCmd(s)
			cmd := exec.Command(cmds[0], cmds[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
		return nil
	}

	if err = run(template.BeforeScript, solvePath, solveFull, solveFile); err != nil {
		return
	}
	if err = run(template.BeforeScript, brutePath, bruteFull, bruteFile); err != nil {
		return
	}
	if err = run(template.BeforeScript, testsGenPath, testsGenFull, testsGenFile); err != nil {
		return
	}

	solveScript := filter(template.Script, solvePath, solveFull, solveFile)
	bruteScript := filter(template.Script, brutePath, bruteFull, bruteFile)
	testsGenScript := filter(template.Script, testsGenPath, testsGenFull, testsGenFile)

	if len(solveScript) == 0 || len(bruteScript) == 0 || len(testsGenScript) == 0 {
		return errors.New("invalid script command. Please check config file")
	}

	testInFormat := strings.ReplaceAll(cfg.DefaultNaming["test_in"], "$%task%$", task)

	numberOfWorkers := 10

	wg := sync.WaitGroup{}
	wg.Add(numberOfWorkers)
	mu := sync.Mutex{}

	workerError := false
	currentTestNumber := 1

	oiejqOptions := &judge.OiejqOptions{MemorylimitInMegaBytes: Args.MemoryLimit, TimeLimitInSeconds: Args.TimeLimit}

	for i := 1; i <= numberOfWorkers; i++ {
		go func(workerID int) {
			defer func() {
				mu.Lock()
				workerError = true
				mu.Unlock()
				wg.Done()
			}()
			for {
				mu.Lock()
				if workerError {
					mu.Unlock()
					return
				}
				testNumber := currentTestNumber
				currentTestNumber++
				mu.Unlock()
				testID := strconv.Itoa(testNumber)
				var genProcessInfo *judge.ProcessInfo
				if Args.Oiejq {
					genProcessInfo, err = judge.RunProcessWithOiejq(testID, testsGenScript, strings.NewReader(testID), oiejqOptions)
				} else {
					genProcessInfo, err = judge.RunProcess(testID, testsGenScript, strings.NewReader(testID), nil)
				}
				if err != nil {
					mu.Lock()
					color.Red(err.Error())
					mu.Unlock()
					return
				}
				var bruteProcessInfo *judge.ProcessInfo
				if Args.Oiejq {
					bruteProcessInfo, err = judge.RunProcessWithOiejq(testID, testsGenScript, strings.NewReader(testID), oiejqOptions)
				} else {
					bruteProcessInfo, err = judge.RunProcess(testID, testsGenScript, strings.NewReader(testID), nil)
				}
				if err != nil {
					mu.Lock()
					color.Red(err.Error())
					mu.Unlock()
					return
				}
				var solveProcessInfo *judge.ProcessInfo
				if Args.Oiejq {
					solveProcessInfo, err = judge.RunProcessWithOiejq(testID, testsGenScript, strings.NewReader(testID), oiejqOptions)
				} else {
					solveProcessInfo, err = judge.RunProcess(testID, testsGenScript, strings.NewReader(testID), nil)
				}
				if err != nil {
					mu.Lock()
					color.Red(err.Error())
					mu.Unlock()
					return
				}
				verdict := judge.GenerateVerdict(testID, judge.Plain(bruteProcessInfo.Output), *solveProcessInfo)
				if !verdict.Correct {
					mu.Lock()
					if workerError {
						mu.Unlock()
						return
					}
					workerError = true
					fmt.Print(verdict.Message)
					err = os.WriteFile(strings.ReplaceAll(testInFormat, "$%test%$", testID), genProcessInfo.Output, 0644)
					color.Red(err.Error())
					mu.Unlock()
					return
				}
				mu.Lock()
				fmt.Print(verdict.Message)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	color.Blue("----FINISHED----")
	return
}
