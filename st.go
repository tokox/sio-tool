package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Arapak/sio-tool/cmd"
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/go-homedir"

	"github.com/docopt/docopt-go"
)

const version = "$CI_VERSION"
const buildTime = "$CI_BUILD_TIME"
const configPath = "~/.st/config"
const codeforcesSessionPath = "~/.st/codeforces_session"
const szkopulSessionPath = "~/.st/szkopul_session"
const sioStaszicSessionPath = "~/.st/sio_staszic_session"
const sioMimuwSessionPath = "~/.st/sio_mimuw_session"
const sioTalentSessionPath = "~/.st/sio_talent_session"

func main() {
	usage := `SIO Tool $%version%$ (st). https://github.com/Arapak/sio-tool
You should run "st config" to configure your handle, password, and code.
templates at first.

If you want to compete, the best command is "st race".

Usage:
  st config
  st submit [-f <file>] [<specifier>...]
  st list [<specifier>...]
  st parse [<specifier>...]
  st gen [<alias>]
  st test [--oiejq] [--memory_limit <memory_limit>] [--time_limit <time_limit>] [<file>]
  st package_test [--oiejq] [--verbose] [--memory_limit <memory_limit>] [--time_limit <time_limit>] [<file>]
  st add_package <file>
  st watch [all] [<specifier>...]
  st open [<specifier>...]
  st stand [<specifier>...]
  st sid [<specifier>...]
  st race [<specifier>...]
  st pull [ac] [<specifier>...]
  st stress-test [--oiejq] [--memory_limit <memory_limit>] [--time_limit <time_limit>] <specifier> [-s <solve>] [-b <brute>] [-g <generator>]
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
  					   Path to the brute force solution file
  -g <generator>, --generator <generator>, <generator>
  					   Path to the test generator file
  -f <file>, --file <file>, <file>
                       Path to the file. E.g. "a.cpp", "./temp/a.cpp"
  --source <source>, <source> 
					   For example, the site from which the tasks originate (codeforces, szkopul)
  -n <name>, --name <name>, <name>
					   Problem name
  -p <path>, --path <path>, <path>
					   Path to a folder where there is a solution to a problem
  -l <link>, --link <link>, <link>
					   Link to the problem site
  -c <contest>, --contest <contest>, <contest>
					   Problem's contest ID
  --shortname <shortname>, <shortname>
					   Problem shortname
  --stage <stage>, <stage>
					   Problem's contest stage ID
  <specifier>          Any useful text E.g.
                       "https://codeforces.com/contest/100",
                       "https://codeforces.com/contest/180/problem/A",
                       "https://codeforces.com/group/Cw4JRyRGXR/contest/269760",
                       "https://szkopul.edu.pl/problemset/problem/kQ5ExYNkFhx3K2FvVuXAAbn4/site/?key=statement",
                       "1111A", "1111", "a", "Cw4JRyRGXR"
                       You can combine multiple specifiers to specify what you
                       want.
  <alias>              Template's alias, e.g., "cpp"
  ac                   The status of the submission is Accepted.
  -o, --oiejq          Use oiejq for running tests
  -v, --verbose        Print verdict of every test
  -m <memory_limit>, --memory_limit <memory_limit>, <memory_limit>
             Set oiejq's memory limit in MiB (default is 1024 (1 GiB))
  -t <time_limit>, --time_limit <time_limit>, <time_limit>  
             Set oiejq's time limit in seconds (default is 10s)

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
  st parse 100         Fetch all problems' samples from contest 100 into
                       "{st}/{contest}/100/".
  st parse gym 100001a
                       Fetch samples of problem "a" of gym 100001 into
                       "{st}/{gym}/100001/a".
  st parse gym 100001
                       Fetch all problems' samples of gym 100001 into
                       "{st}/{gym}/100001".
  st parse             Fetch samples of the current problem onto the current path.
  st gen               Generate code from the default template.
  st gen cpp           Generate a code from the template whose alias is "cpp"
                       into the current path.
  st test              Run the commands of a template in the current path. Then
                       test all samples. If you want to add a new test case,
                       Create two files, "inK.txt" and "outK.txt" where K is
                       a string with 0~9.
  st add_package ~/tests
                       Add package (set of tests) for a task you are currently in 
  st test_package      Test your solution on a package added before
  st watch             Watch the first 10 submissions for the current contest.
  st watch all         Watch all submissions for the current contest.
  st open 1136a        Use your default web browser to open the page for the contest.
                       1136, problem a.
  st open gym 100136   Use the default web browser to open the page of gym.
                       100136.
  st stand             Use the default web browser to open the standing page.
  st sid 52531875      Use the default web browser to open the submission.
                       52531875's page.
  st sid               Open the last submission's page.
  st race 1136         If the contest 1136 has not started yet, it will
                       countdown. When the countdown ends, it will open all
                       problems' pages and parse samples.
  st pull 100          Pull all problems' latest codes from contest 100 into
                       "./100/<problem-id>".
  st pull 100 a        Pull the latest code of problem "a" of contest 100 into
                       "./100/<problem-id>".
  st pull ac 100 a     Pull the "Accepted" or "Pretests passed" code of the problem.
                       "a" of contest 100.
  st pull              Pull the latest codes for the current problem into the current
                       path.
  st stress-test abc   Stresstest a program with your solve, brute force solution, and test generator.
  st db add            Add a new task to the database with problems you solved (problems parsed by sio-tool are automatically added).
  st db find -n "square"
					   Find all problems in the database that contain the string "square" (ignoring capitalization).
  st db goto -n "square" -c 100
					   Returns the path of the task with a name that contains "square" and has contest id 100 (if you configure your shell correctly, it can automatically cd into the path (example of .bashrc in CONFIG.md))
  st upgrade           Upgrade the "st" to the latest version from GitHub.


Files:
  st will save some data in some files:

  "~/.st/config"        Configuration file, including templates, etc.
  "~/.st/codeforces_session"    Codeforces session file, including cookies, handle, password, etc.
  "~/.st/szkopul_session"       Szkopul session file, including username and password
  "~/.st/sio_session"           Sio session file, including username and password

  "~" is the home directory of the current user on your system.

  Don't share the session files with anyone, your password is encrypted, but if someone knows you are using this program, he can easily decrypt it.

Template:
  You can insert some placeholders into your template code. When generating a code
  from the template, st will replace all placeholders by the following rules:

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
  You could set "before_script" or "after_script" to an empty string, meaning
  not executing.
  You have to run your program in "script" with standard input/output (no
  need to redirect).

  You can insert some placeholders in your scripts. When executing a script,
  st will replace all placeholders by the following rules:

  $%path%$   Path to source file (Excluding $%full%$, e.g. "/home/arapak/")
  $%full%$   Full name of source file (e.g. "a.cpp")
  $%file%$   Name of source file (Excluding suffix, e.g. "a")
  $%rand%$   Random string with 8 characters (including "a-z" "0-9")`

	color.Output = ansi.NewAnsiStdout()

	usage = strings.Replace(usage, `$%version%$`, version, 1)
	opts, _ := docopt.ParseArgs(usage, os.Args[1:], fmt.Sprintf("SIO Tool (st) %v\nLast built: %v\n", version, buildTime))
	opts[`{version}`] = version

	cfgPath, _ := homedir.Expand(configPath)
	codeforcesClnPath, _ := homedir.Expand(codeforcesSessionPath)
	szkopulClnPath, _ := homedir.Expand(szkopulSessionPath)
	sioStaszicClnPath, _ := homedir.Expand(sioStaszicSessionPath)
	sioMimuwClnPath, _ := homedir.Expand(sioMimuwSessionPath)
	sioTalentClnPath, _ := homedir.Expand(sioTalentSessionPath)
	config.Init(cfgPath)
	codeforces_client.Init(codeforcesClnPath, config.Instance.CodeforcesHost, config.Instance.Proxy)
	szkopul_client.Init(szkopulClnPath, config.Instance.SzkopulHost, config.Instance.Proxy)
	sio_client.Init(sioStaszicClnPath, config.Instance.SioStaszicHost, config.Instance.Proxy, sio_client.Staszic)
	sio_client.Init(sioMimuwClnPath, config.Instance.SioMimuwHost, config.Instance.Proxy, sio_client.Mimuw)
	sio_client.Init(sioTalentClnPath, config.Instance.SioTalentHost, config.Instance.Proxy, sio_client.Talent)

	err := cmd.Eval(opts)
	if err != nil {
		fmt.Println(util.RedString(err.Error()))
		os.Exit(1)
	}
	color.Unset()
}
