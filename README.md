# SIO Tool

[![Github release](https://img.shields.io/github/release/Arapak/sio-tool.svg)](https://github.com/Arapak/sio-tool/releases)
[![platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue.svg)](https://github.com/Arapak/sio-tool/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arapak/sio-tool)](https://goreportcard.com/report/github.com/Arapak/sio-tool)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.12-green.svg)](https://github.com/golang)
[![license](https://img.shields.io/badge/license-MIT-%23373737.svg)](https://raw.githubusercontent.com/Arapak/sio-tool/main/LICENSE)

SIO Tool is a command-line interface tool for [Codeforces](https://codeforces.com), [Szkopul (OI archive)](https://szkopul.edu.pl/task_archive/oi/), [SIO2 (staszic)](https://sio2.staszic.waw.pl) and [SIO2 (mimuw)](https://sio2.mimuw.edu.pl).

It's fast, small, cross-platform and powerful.

[Installation](##installation) | [Usage](##usage) | [FAQ](##faq)

## Features

-   Supports Contests, Gym, Groups and acmsguru.
-   Supports all programming languages in Codeforces.
-   Submit codes.
-   Watch submissions' status dynamically.
-   Fetch problems' samples.
-   Compile and test locally.
-   Generate codes from the specified template (including timestamp, author, etc.)
-   List problems' stats of one contest.
-   Use default web browser to open problems' pages, standings' page, etc.
-   Setup a network proxy. Setup a mirror host.
-   Colorful CLI.

Pull requests are always welcome.

## Installation

You can download the pre-compiled binary file in [here](https://github.com/Arapak/sio-tool/releases).

Then enjoy the sio-tool~

Or you can compile it from the source **(go >= 1.12)**:

```plain
$ go get github.com/Arapak/sio-tool
$ cd $GOPATH/src/github.com/Arapak/sio-tool
$ go build -ldflags "-s -w" st.go
```

If you don't know what's the `$GOPATH`, please see here <https://github.com/golang/go/wiki/GOPATH>.

## Configuration
To configure this program to your preferences you will have to run the command `st config`

You will then be able to: login, add and remove templates, set hosts, set folders naming, change your database path.

Alternatively you can edit the `~/.st/config` file

## Usage

Let's simulate a codeforces competition.

You will have to start in your codeforces root path (you can configure it, the default is `~/st/codeforces`).

`st race 1136` or `st race https://codeforces.com/contest/1136`

To start competing the contest 1136!

If the contest has not started yet, `st` will count down. If the contest have started or the countdown ends, `st` will use the default browser to open dashboard's page and problems' page, and fetch all samples to the local.

`cd ./st/contest/1136/a` (May be different from this, please notice the message on your screen)

Enter the directory of problem A, the directory should contain all samples of the problem.

`st gen`

Generate a code with the default template. The filename of the code is problem id by default.

`vim a.cpp`

Use Vim to write the code (It depends on yourself).

`st test`

Compile and test all samples.

`st submit`

Submit the code.

`st list`

List problems' stats of the contest.

`st stand`

Open the standings' page of the contest.

```plain
You should run "st config" to configure your handle, password and code
templates at first.

If you want to compete, the best command is "st race".

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
                       "{st}/{contest}/100/".
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
                       create two files "inK.txt" and "outK.txt" where K is
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
  st db add            Add a new task to the database with problems you solved (problems parsed by sio-tool are automatically added)
  st db find -n "square"
					   Find all problems in the database which contain the string "square" (ignoring capitalization)
  st db goto -n "square" -c 100
					   Returns the path of the task with name that contains "square" and has contest id 100 (if you configure your shell correctly it can automatically cd in to the path (example of .bashrc below))
  st upgrade           Upgrade the "st" to the latest version from GitHub.


File:
  st will save some data in some files:

  "~/.st/config"        Configuration file, including templates, etc.
  "~/.st/codeforces_session"    Codeforces session file, including cookies, handle, password, etc.
  "~/.st/szkopul_session"       Szkopul session file, including username and password
  "~/.st/sio_session"           Sio session file, including username and password

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
```

## Template Example

The placeholders inside the template will be replaced with the corresponding content when you run `st gen`.

```
$%U%$   Handle (e.g. xalanq)
$%Y%$   Year   (e.g. 2019)
$%M%$   Month  (e.g. 04)
$%D%$   Day    (e.g. 09)
$%h%$   Hour   (e.g. 08)
$%m%$   Minute (e.g. 05)
$%s%$   Second (e.g. 00)
```

```cpp
/* Generated by powerful Sio Tool
 * You can download the binary file in here https://github.com/Arapak/sio-tool (Linux, Windows, MacOS)
 * Author: $%U%$
 * Time: $%Y%$-$%M%$-$%D%$ $%h%$:$%m%$:$%s%$
**/

#include <bits/stdc++.h>
using namespace std;

typedef long long ll;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(0);

    return 0;
}
```

## Configure your shell

To have some additional features, you can also configure your shell by adding the script below to your `.bashrc` file (or `.zshrc`)

```bash
st() {
  if [ "$1" = stress-test ]; then
		# the command below runs stress tests in a separate terminal. If you want this behaviour in your terminal, you have to find out what command to use. This example is for alacritty 
    command alacritty --hold -e st "$@"
  elif [ "$1" = db ] && [ "$2" = goto ]; then
		res=$(command st "$@")
		code=$?
		if [ "$code" = 0 ]; then
			cd $res
		else
			echo -e "$res"
		fi
	else
  	command st "$@"
  fi
}
```

## FAQ

### I double click the program but it doesn't work

SIO Tool is a command-line tool. You should run it in terminal.

### I cannot use `st` command

You should put the `st` program to a path (e.g. for Linux `/usr/bin/`, or even better create a `~/.local/bin` and add the whole folder to PATH in your `.bashrc` (or `.zshrc`) using the script below, it is necessary for the upgrade script to work) which has been added to system environment variable PATH.

```bash
export PATH=$PATH:$HOME/.local/bin
```

Or just google "how to add a path to system environment variable PATH".

### How to add a new testcase

Create two extra testcase files `inK.txt` and `outK.txt` (K is a string with 0~9).
Or a different possible naming would be `$%file%$K.in` and `$%file%$K.out` where `$%file%$` is the file name of your program

### Enable tab completion in terminal

Use this [Infinidat/infi.docopt_completion](https://github.com/Infinidat/infi.docopt_completion).

Note: If there is a new version released (especially a new command added), you should run `docopt-completion st` again.

### Can't parse tasks from szkopul

You get the error: `Error: exec: "pdftotext": executable file not found in $PATH`?

This program depends on package named(on most os's) poppler-utils. Below, example install on a debian-based os (ex. ubuntu, linux mint):
```bash
sudo apt install poppler-utils
```
