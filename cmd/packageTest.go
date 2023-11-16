package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/judge"
	"github.com/Arapak/sio-tool/util"
	"github.com/k0kubun/go-ansi"

	"github.com/fatih/color"
)

const ErrorPackageNotFound = "no package found"

const ErrorTestsNotFound = "no tests found"

func getOnePackage(path string) (packagePath string, err error) {
	paths, err := os.ReadDir(path)
	if err != nil {
		return
	}
	var packages []string
	var packagesMessage []string
	for _, path := range paths {
		if path.IsDir() {
			packages = append(packages, path.Name())
			info, err := path.Info()
			if err != nil {
				return "", err
			}
			packagesMessage = append(packagesMessage, fmt.Sprintf("%v (added %v)", path.Name(), info.ModTime().Format("2006-01-02 15:04")))
		}
	}
	if len(packages) == 0 {
		return "", errors.New(ErrorPackageNotFound)
	} else if len(packages) == 1 {
		return packages[0], nil
	} else {
		packageIndex := 0
		prompt := &survey.Select{
			Message: "Multiple packages can be selected.",
			Options: packagesMessage,
		}
		if err = survey.AskOne(prompt, &packageIndex); err != nil {
			return
		}
		return packages[packageIndex], nil
	}
}

type testPattern struct {
	in  string
	out string
}

var testPatterns = [...]testPattern{
	{`^in(\w+)\.txt$`, `^out(\w+)\.txt$`},
	{`^(\w+)\.in$`, `^(\w+)\.out$`},
	{`^in/(\w+)\.in$`, `^out/(\w+)\.out$`},
	{`^in/in(\w+)$`, `^out/out(\w+)$`},
	{`^in/(\w+)$`, `^out/(\w+)$`},
}

func checkMatching(s string, pattern string) (string, bool) {
	reg := regexp.MustCompile(pattern)
	tmp := reg.FindStringSubmatch(s)
	if len(tmp) < 2 {
		return "", false
	}
	return string(tmp[1]), true
}

func getTestsByPattern(paths []string, pattern testPattern) (in []string, out []string) {
	inFiles := make(map[string]int)
	for i, path := range paths {
		if val, match := checkMatching(path, pattern.in); match {
			inFiles[val] = i
		}
	}
	for _, path := range paths {
		if val, match := checkMatching(path, pattern.out); match {
			if i, ok := inFiles[val]; ok {
				in = append(in, paths[i])
				out = append(out, path)
			}
		}
	}
	return
}

func getAllTests(path string) (in []string, out []string, err error) {
	var paths []string
	err = filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			paths = append(paths, strings.TrimPrefix(filePath, filepath.Clean(path)+"/"))
			return nil
		})
	if err != nil {
		return
	}
	for _, pattern := range testPatterns {
		in, out = getTestsByPattern(paths, pattern)
		if len(in) != 0 {
			return
		}
	}
	return nil, nil, errors.New(ErrorTestsNotFound)
}

func printVerdict(verdict judge.Verdict, testID string) {
	if verdict.Status == judge.OK || verdict.Status == judge.WA {
		fmt.Print(verdict.Message)
	} else if verdict.Status == judge.RE {
		color.Red("runtime error #%v: %v", testID, verdict.Err.Error())
	} else if verdict.Status == judge.TLE {
		color.Red("time limit exceeded #%v", testID)
	} else if verdict.Status == judge.MLE {
		color.Red("memory limit exceeded #%v", testID)
	} else if verdict.Status == judge.OLE {
		color.Red("output limit exceeded #%v", testID)
	} else if verdict.Status == judge.INT {
		color.Red("internal error #%v: %v", testID, verdict.Err.Error())
	}
}

func printReport(m map[judge.VerdictStatus]int, testsRan int, maxTime, maxMemory float64) {
	_, _ = ansi.Printf("TESTS RAN: %v", util.BlueString(fmt.Sprint(testsRan)))
	_, _ = ansi.Printf(" MAX TIME: %0.3fs", maxTime)
	_, _ = ansi.Printf(" MAX MEMORY: %v", judge.ParseMemory(maxMemory))
	for status, num := range m {
		if status == judge.OK {
			_, _ = ansi.Printf(" OK: %v", util.GreenString(fmt.Sprint(num)))
		} else {
			_, _ = ansi.Printf(" %v: %v", status, util.RedString(fmt.Sprint(num)))
		}
	}
}

func PackageTest() (err error) {
	cfg := config.Instance
	if len(cfg.Template) == 0 {
		return errors.New("you have to add at least one code template by `st config`")
	}

	filename, index, err := getOneCode(Args.File, cfg.Template, map[string]struct{}{})
	if err != nil {
		return
	}
	template := cfg.Template[index]
	path, full := filepath.Split(filename)
	ext := filepath.Ext(filename)
	file := full[:len(full)-len(ext)]
	rand := util.RandString(8)

	packagesPath, err := ArgsPackagePath()
	if err != nil {
		return
	}
	packagePath, err := getOnePackage(packagesPath)
	if err != nil {
		return
	}
	packagePath = filepath.Join(packagesPath, packagePath)
	in, out, err := getAllTests(packagePath)
	if err != nil {
		return
	}

	filter := func(cmd string) string {
		cmd = strings.ReplaceAll(cmd, "$%rand%$", rand)
		cmd = strings.ReplaceAll(cmd, "$%path%$", path)
		cmd = strings.ReplaceAll(cmd, "$%full%$", full)
		cmd = strings.ReplaceAll(cmd, "$%file%$", file)
		return cmd
	}

	run := func(script string) error {
		if s := filter(script); len(s) > 0 {
			fmt.Println(s)
			cmds := util.SplitCmd(s)
			cmd := exec.Command(cmds[0], cmds[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
		return nil
	}

	if err = run(template.BeforeScript); err != nil {
		return
	}

	numberOfWorkers := 10

	wg := sync.WaitGroup{}
	wg.Add(numberOfWorkers)
	mu := sync.Mutex{}

	currentTestNumber := 0

	runScript := filter(template.Script)

	var oiejqOptions *judge.OiejqOptions
	if Args.Oiejq {
		err = judge.InstallSio2Jail()
		if err != nil {
			return
		}
		oiejqOptions = &judge.OiejqOptions{MemorylimitInMegaBytes: Args.MemoryLimit, TimeLimitInSeconds: Args.TimeLimit}
	}

	m := make(map[judge.VerdictStatus]int)
	testsRan := 0
	maxTime := 0.0
	maxMemory := 0.0

	for i := 1; i <= numberOfWorkers; i++ {
		go func(workerID int) {
			defer func() {
				wg.Done()
			}()

			for {
				mu.Lock()
				testNumber := currentTestNumber
				currentTestNumber++
				if testNumber >= len(in) {
					mu.Unlock()
					return
				}
				mu.Unlock()

				verdict := judge.Judge(filepath.Join(packagePath, in[testNumber]), filepath.Join(packagePath, out[testNumber]), in[testNumber], runScript, oiejqOptions)

				mu.Lock()
				ansi.EraseInLine(2)
				ansi.CursorHorizontalAbsolute(0)
				if Args.Verbose {
					printVerdict(verdict, in[testNumber])
				}
				m[verdict.Status]++
				testsRan++
				maxTime = math.Max(maxTime, verdict.TimeInSeconds)
				maxMemory = math.Max(maxMemory, verdict.MemoryInMegabytes)
				printReport(m, testsRan, maxTime, maxMemory)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	color.Blue("\n----FINISHED----")
	return
}
