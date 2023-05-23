package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/process"
)

func splitCmd(s string) (res []string) {
	// https://github.com/vrischmann/shlex/blob/master/shlex.go
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}

func plain(raw []byte) string {
	buf := bufio.NewScanner(bytes.NewReader(raw))
	var b bytes.Buffer
	newline := []byte{'\n'}
	for buf.Scan() {
		b.Write(bytes.TrimSpace(buf.Bytes()))
		b.Write(newline)
	}
	return b.String()
}

func parseMemory(memory uint64) string {
	if memory > 1024*1024 {
		return fmt.Sprintf("%.3fMB", float64(memory)/1024.0/1024.0)
	} else if memory > 1024 {
		return fmt.Sprintf("%.3fKB", float64(memory)/1024.0)
	}
	return fmt.Sprintf("%vB", memory)
}

type ProcessInfo struct {
	time   float64
	memory string
	output []byte
}

func runProcess(processID, command string, input io.Reader) (*ProcessInfo, error) {
	var o bytes.Buffer
	output := io.Writer(&o)

	cmds := splitCmd(command)

	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("runtime error #%v ... %v", processID, err.Error())
	}

	pid := int32(cmd.Process.Pid)
	maxMemory := uint64(0)
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	running := true
	for running {
		select {
		case err := <-ch:
			if err != nil {
				return nil, fmt.Errorf("runtime error #%v ... %v", processID, err.Error())
			}
			running = false
		default:
			p, err := process.NewProcess(pid)
			if err == nil {
				m, err := p.MemoryInfo()
				if err == nil && m.RSS > maxMemory {
					maxMemory = m.RSS
				}
			}
		}
	}
	return &ProcessInfo{cmd.ProcessState.UserTime().Seconds(), parseMemory(maxMemory), o.Bytes()}, nil
}

type Verdict struct {
	correct bool
	message string
	err     error
}

func generateVerdict(testID, answer string, processInfo ProcessInfo) Verdict {
	state := ""
	diff := ""
	var correct bool
	output := plain(processInfo.output)
	if output == answer {
		correct = true
		state = color.New(color.FgGreen).Sprintf("Passed #%v", testID)
	} else {
		correct = false
		state = color.New(color.FgRed).Sprintf("Failed #%v", testID)
		diff += output + "\n"
		diff += answer + "\n"
	}
	return Verdict{correct, fmt.Sprintf("%v ... %.3fs %v\n%v", state, processInfo.time, processInfo.memory, diff), nil}
}

func judge(sampleID, inPathFormat, ansPathFormat, command string) Verdict {
	inPath := fmt.Sprintf(inPathFormat, sampleID)
	ansPath := fmt.Sprintf(ansPathFormat, sampleID)
	input, err := os.Open(inPath)
	if err != nil {
		return Verdict{false, "", err}
	}
	defer input.Close()

	processInfo, err := runProcess(sampleID, command, input)
	if err != nil {
		return Verdict{false, "", err}
	}

	b, err := os.ReadFile(ansPath)
	if err != nil {
		return Verdict{false, "", err}
	}
	return generateVerdict(sampleID, plain(b), *processInfo)
}

func ExtractTaskName(file string) (task string) {
	task, _, _ = strings.Cut(file, "-")
	return
}

// Test command
func Test() (err error) {
	cfg := config.Instance
	if len(cfg.Template) == 0 {
		return errors.New("you have to add at least one code template by `st config`")
	}

	filename, index, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}
	template := cfg.Template[index]
	path, full := filepath.Split(filename)
	ext := filepath.Ext(filename)
	file := full[:len(full)-len(ext)]
	rand := util.RandString(8)
	task := ExtractTaskName(file)

	samples := getSampleByName(task)
	samplesWithName := true
	if len(samples) == 0 {
		samplesWithName = false
		samples = getSampleID()
		if len(samples) == 0 {
			return errors.New("cannot find any sample file")
		}
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
			cmds := splitCmd(s)
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
	if s := filter(template.Script); len(s) > 0 {
		for _, i := range samples {
			var verdict Verdict
			if samplesWithName {
				verdict = judge(i, fmt.Sprintf("%s%%v.in", file), fmt.Sprintf("%s%%v.out", file), s)
			} else {
				verdict = judge(i, "in%v.txt", "out%v.txt", s)
			}

			if verdict.err != nil {
				color.Red(err.Error())
			} else {
				fmt.Print(verdict.message)
			}
		}
	} else {
		return errors.New("invalid script command, please check config file")
	}
	return run(template.AfterScript)
}
