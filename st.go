package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Arapak/sio-tool/cmd"
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/mitchellh/go-homedir"

	docopt "github.com/docopt/docopt-go"
)

const version = "$CI_VERSION"
const buildTime = "$CI_BUILD_TIME"
const configPath = "~/.st/config"
const codeforcesSessionPath = "~/.st/codeforces_session"
const szkopulSessionPath = "~/.st/szkopul_session"

func main() {
	usage := `SIO Tool $%version%$ (st). https://github.com/Arapak/sio-tool
You should run "st config" to configure your handle, password and code
templates at first.
If you want to compete, the best command is "st race"
Usage:
  st config
  st submit [-f <file>] [<specifier>...]
  st list [<specifier>...]
  st parse [<specifier>...]
  st gen [<alias>]
  st test [<file>]
  st watch [all] [<specifier>...]
  st open [<specifier>...]
  st stand [<specifier>...]
  st sid [<specifier>...]
  st race [<specifier>...]
  st pull [ac] [<specifier>...]
  st stress-test <specifier> [-s <solve>] [-b <brute>] [-g <generator>]
  st db add [--source <source>] [-n <name>] [-p <path>] [-l <link>] [-c <contest>] [--shortname <shortname>] [--stage <stage>]
  st db find [--source <source>] [-n <name>] [-p <path>] [-l <link>] [-c <contest>] [--shortname <shortname>] [--stage <stage>]
  st db goto [--source <source>] [-n <name>] [-p <path>] [-l <link>] [-c <contest>] [--shortname <shortname>] [--stage <stage>]
  st upgrade
Options:
  -h --help            Show this screen.
  --version            Show version.
  -s <solve>, --solve <solve>, <solve>
  					   Path to solve file
  -b <brute>, --brute <brute>, <brute>
  					   Path to brute force solution file
  -g <generator>, --generator <generator>, <generator>
  					   Path to tests generator file
  -f <file>, --file <file>, <file>
                       Path to file. E.g. "a.cpp", "./temp/a.cpp"
  --source <source>, <source>
					   For example the site from which the tasks originates (codeforces, szkopul)
  -n <name>, --name <name>, <name>
					   Problem name
  -p <path>, --path <path>, <path>
					   Path to folder where is a solution to a problem
  -l <link>, --link <link>, <link>
					   Link to problem site
  -c <contest>, --contest <contest>, <contest>
					   Problem's contest id
  --shortname <shortname>, <shortname>
					   Problem shortname
  --stage <stage>, <stage>
					   Problem's contest stage id
  <specifier>          Any useful text. E.g.
                       "https://codeforces.com/contest/100",
                       "https://codeforces.com/contest/180/problem/A",
                       "https://codeforces.com/group/Cw4JRyRGXR/contest/269760",
                       "https://szkopul.edu.pl/problemset/problem/kQ5ExYNkFhx3K2FvVuXAAbn4/site/?key=statement",
                       "1111A", "1111", "a", "Cw4JRyRGXR"
                       You can combine multiple specifiers to specify what you
                       want.
  <alias>              Template's alias. E.g. "cpp"
  ac                   The status of the submission is Accepted.
Examples:
  st config            Configure the sio-tool.
  st submit            st will detect what you want to submit automatically.
  st submit -f a.cpp
  st submit https://codeforces.com/contest/100/A
  st submit -f a.cpp 100A 
  st submit -f a.cpp 100 a
  st submit contest 100 a
  st submit gym 100001 a
  st list              List all problems' stats of a contest.
  st list 1119
  st parse 100         Fetch all problems' samples of contest 100 into
                       "{st}/{contest}/100/<problem-id>".
  st parse gym 100001a
                       Fetch samples of problem "a" of gym 100001 into
                       "{st}/{gym}/100001/a".
  st parse gym 100001
                       Fetch all problems' samples of gym 100001 into
                       "{st}/{gym}/100001".
  st parse             Fetch samples of current problem into current path.
  st gen               Generate a code from default template.
  st gen cpp           Generate a code from the template whose alias is "cpp"
                       into current path.
  st test              Run the commands of a template in current path. Then
                       test all samples. If you want to add a new testcase,
                       create two files "inK.txt" and "ansK.txt" where K is
                       a string with 0~9.
  st watch             Watch the first 10 submissions of current contest.
  st watch all         Watch all submissions of current contest.
  st open 1136a        Use default web browser to open the page of contest
                       1136, problem a.
  st open gym 100136   Use default web browser to open the page of gym
                       100136.
  st stand             Use default web browser to open the standing page.
  st sid 52531875      Use default web browser to open the submission
                       52531875's page.
  st sid               Open the last submission's page.
  st race 1136         If the contest 1136 has not started yet, it will
                       countdown. When the countdown ends, it will open all
                       problems' pages and parse samples.
  st pull 100          Pull all problems' latest codes of contest 100 into
                       "./100/<problem-id>".
  st pull 100 a        Pull the latest code of problem "a" of contest 100 into
                       "./100/<problem-id>".
  st pull ac 100 a     Pull the "Accepted" or "Pretests passed" code of problem
                       "a" of contest 100.
  st pull              Pull the latest codes of current problem into current
                       path.
  st stress-test abc   Stresstest a program with your solve, brute force solution and tests generator.
  st upgrade           Upgrade the "st" to the latest version from GitHub.
File:
  st will save some data in some files:
  "~/.st/config"        Configuration file, including templates, etc.
  "~/.st/session"       Session file, including cookies, handle, password, etc.
  "~" is the home directory of current user in your system.
Template:
  You can insert some placeholders into your template code. When generate a code
  from the template, st will replace all placeholders by following rules:
  $%U%$   Handle (e.g. Arapak)
  $%Y%$   Year   (e.g. 2019)
  $%M%$   Month  (e.g. 04)
  $%D%$   Day    (e.g. 09)
  $%h%$   Hour   (e.g. 08)
  $%m%$   Minute (e.g. 05)
  $%s%$   Second (e.g. 00)
Script in template:
  Template will run 3 scripts in sequence when you run "st test":
    - before_script   (execute once)
    - script          (execute the number of samples times)
    - after_script    (execute once)
  You could set "before_script" or "after_script" to empty string, meaning
  not executing.
  You have to run your program in "script" with standard input/output (no
  need to redirect).
  You can insert some placeholders in your scripts. When execute a script,
  st will replace all placeholders by following rules:
  $%path%$   Path to source file (Excluding $%full%$, e.g. "/home/arapak/")
  $%full%$   Full name of source file (e.g. "a.cpp")
  $%file%$   Name of source file (Excluding suffix, e.g. "a")
  $%rand%$   Random string with 8 character (including "a-z" "0-9")
  $%task%$   Task name supplied when using stress-test command`
	color.Output = ansi.NewAnsiStdout()

	usage = strings.Replace(usage, `$%version%$`, version, 1)
	opts, _ := docopt.ParseArgs(usage, os.Args[1:], fmt.Sprintf("SIO Tool (st) %v\nLast built: %v\n", version, buildTime))
	opts[`{version}`] = version

	cfgPath, _ := homedir.Expand(configPath)
	codeforcesClnPath, _ := homedir.Expand(codeforcesSessionPath)
	szkopulClnPath, _ := homedir.Expand(szkopulSessionPath)
	config.Init(cfgPath)
	codeforces_client.Init(codeforcesClnPath, config.Instance.CodeforcesHost, config.Instance.Proxy)
	szkopul_client.Init(szkopulClnPath, config.Instance.SzkopulHost, config.Instance.Proxy)

	err := cmd.Eval(opts)
	if err != nil {
		fmt.Println(util.RedString(err.Error()))
		os.Exit(1)
	}
	color.Unset()
}
