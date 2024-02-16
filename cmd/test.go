package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/judge"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
)

func Test() (err error) {
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
	task := judge.ExtractTaskName(file)

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

	var oiejqOptions *judge.OiejqOptions
	if Args.Oiejq {
		err = judge.InstallSio2Jail()
		if err != nil {
			return
		}
		oiejqOptions = &judge.OiejqOptions{MemorylimitInMegaBytes: Args.MemoryLimit, TimeLimitInSeconds: Args.TimeLimit}
	}

	if s := filter(template.Script); len(s) > 0 {
		for _, i := range samples {
			var verdict judge.Verdict

			if samplesWithName {
				verdict = judge.Judge(fmt.Sprintf("%s%v.in", task, i), fmt.Sprintf("%s%v.out", task, i), i, s, oiejqOptions)
			} else {
				verdict = judge.Judge(fmt.Sprintf("in%v.txt", i), fmt.Sprintf("out%v.txt", i), i, s, oiejqOptions)
			}

			if verdict.Err != nil {
				color.Red(verdict.Err.Error())
			} else {
				fmt.Print(verdict.Message)
			}
		}
	} else {
		return errors.New("invalid script command, please check config file")
	}
	return run(template.AfterScript)
}
