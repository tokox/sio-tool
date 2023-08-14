# SIO Tool

[![Github release](https://img.shields.io/github/release/Arapak/sio-tool.svg)](https://github.com/Arapak/sio-tool/releases)
[![platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue.svg)](https://github.com/Arapak/sio-tool/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arapak/sio-tool)](https://goreportcard.com/report/github.com/Arapak/sio-tool)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-green.svg)](https://github.com/golang)
[![license](https://img.shields.io/badge/license-MIT-%23373737.svg)](https://raw.githubusercontent.com/Arapak/sio-tool/main/LICENSE)

SIO Tool is a command-line interface tool for [Codeforces](https://codeforces.com), [Szkopul (OI archive)](https://szkopul.edu.pl/task_archive/oi/), [SIO2 (staszic)](https://sio2.staszic.waw.pl) and [SIO2 (mimuw)](https://sio2.mimuw.edu.pl).

It's fast, small, cross-platform, and powerful.

[Installation](#installation) | [Configuration](#configuration) | [Usage](#usage) | [FAQ](#faq)

## Features

- Supports Codeforces (Contests, Gym, Groups, and acmsguru), Sio and Szkopul (OI Archive).
- Supports all programming languages in Codeforces (Sio and Szkopul support only C++).
- Submit codes.
- Watch submissions' status dynamically.
- Fetch problems' samples.
- Compile and test locally.
- Generate codes from the specified template (including timestamp, author, etc.).
- List problems' stats for one contest.
- Use default web browser to open problems' pages, standings' pages, etc.
- Set up a network proxy. Set up a mirror host.
- Colorful CLI

Pull requests are always welcome.

## Installation

Please refer to the [INSTALL.md](/INSTALL.md) file

## Configuration

Please refer to the [CONFIG.md](/CONFIG.md) file

## Usage

Let's explain the structure of the folders created by sio-tool: First, there are roots for every website (by default, they are in corresponding folders in `~/st`), then it differs based on the website (we will explain them below).
### Codeforces

Folders structure:
- First layer is the section in Codeforces (contests, gym, etc.).
- Then there are the actual contests (the folder names are the contests' IDs).
- Next there are problem aliases (for example: a, b, c, etc.).
- And here is your solution file and the sample test cases.

Let's simulate a Codeforces competition.

You will have to start in your codeforces root path (you can configure it, the default is `~/st/codeforces`).

`st race 1136` or `st race https://codeforces.com/contest/1136`

To start competing in the contest 1136!

If the contest has not started yet, `st` will count down. If the contest has started or the countdown ends, `st` will use the default browser to open the dashboard's page and problems' page and fetch all samples to the local.

`cd ./contest/1136/a` (This may be different from this; please notice the message on your screen.)

Enter the directory for problem A; the directory should contain all samples of the problem.

`st gen`

Generate code with the default template. The filename of the code is problem id by default.

`vim a.cpp`

Use Vim to write the code (It depends on yourself).

`st test`

Compile and test all samples.

`st submit`

Submit the code.

`st list`

List problems' stats for the contest.

`st stand`

Open the Standings page of the contest.

### Szkopul

Folders structure:
- Archives (currently only supporting OI)
- Stage
- Task
- Your code and samples


Let's say you want to solve the problem "Meteory" from XXIX OI. You should start in the szkopul's root folder (by default `~/st/szkopul`), then write


`st parse XXIX met`

Then sio-tool will create a directory `XXIX` in which there will be a directory for stage 3, and after that, `met` folder for the task "Meteory".
Then change the directory to the `met` folder.
`cd XXIX/3/met`

In this folder, there will be 2 files, `in1.txt` and `out1.txt`. These are the samples for the "Meteory" problem.

Now you can generate a file for your solution using the default template (which you have to configure first, refer to the configuration section).

`st gen`

For Szkopul, your file extension has to be one of: `cpp`, `cc`, `c`, `py` or `pas` (some problems may even accept only `cpp` or `cc`), so only files with those extensions will be allowed.
If you use C++, this will create a `met.cpp` file containing your template.

You now proceed to solve the problem, and when you are ready, you want to test it on the samples.

`st test`

This compiles and runs your program using the scripts you specified in the template.

Your solution passes the samples, and you want to submit it.

`st submit`

You get 100 points, are happy, and now want to check what other problems from this stage you didn't already solve.

`cd ..`

`st list`

You see, you didn't solve the problem "Rzeki", so you want to open its statement page.

`st open rze`

### Sio

Folders structure:
- Contest
- Round
- Task
- Your code and samples

Let's say you want to solve the problem "Permutacje" from Wiekuisty Ontak 2023 on Sio Mimuw, you should start in Sio Mimuw's root folder (by default `~/st/sio-mimuw`), then write

`st parse wiekuisty_ontak2023 per`

Then sio-tool will create a directory `wiekuisty_ontak2023` in there will be the directory `day_2` and after that, `per` folder for the task "Permutacje"

Then change the directory to the `per` folder.

`cd wiekuisty_ontak2023/day_2/per`

In this folder will be 2 files, `in1.txt` and `out1.txt`. These are the samples for the "Permutacje" problem.

Now you can generate a file for your solution using the default template (which you have to configure first; refer to the configuration section).

`st gen`

For Sio, your file extension has to be one of: `cpp`, `cc`, `c`, `py` or `pas` (some contests may even accept only `cpp` or `cc`), so only files with those extensions will be allowed.
If you use C++, this will create a `per.cpp` file containing your template.

You now proceed to solve the problem, and when you are ready, you want to test it on the samples.

`st test`

This compiles and runs your program using the scripts you specified in the template.

Your solution passes the samples, and you want to submit it.

`st submit`

You get 100 points, are happy, and now want to check if your position in the standings changed.

`st stand`

### Database

You vaguely remember a problem but don't know from where; you just remember it was something about chess. Now you can search all the problems you solved using the sio-tool's db command.

`st db find -n "chess"`

Now you have a list of all the problems you solved whose names contain "chess".
(You can do similar things with contest IDs, stages, aliases, etc.)

You solved a problem without using Sio-Tool but want to add it to the database, no problem.
`st db add`
You will be asked several questions, and your problem will be added.

You are working on a problem, but you close your terminal every time.

You can now very easily go back to the correct location just by knowing the name of your task (or a part of it).

`st db goto -n "chess"`

(This command only works after configuring your shell, checkout configuration.)

### All options

```plain
You should run "st config" to configure your handle, password, and code.
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
  $%rand%$   Random string with 8 characters (including "a-z" "0-9")
```

## Template Example

The placeholders inside the template will be replaced with the corresponding content when you run `st gen`.

```
$%U%$   Handle (e.g. Arapak)
$%Y%$   Year   (e.g. 2023)
$%M%$   Month  (e.g. 08)
$%D%$   Day    (e.g. 12)
$%h%$   Hour   (e.g. 20)
$%m%$   Minute (e.g. 05)
$%s%$   Second (e.g. 00)
```

```cpp
/* Generated by the powerful Sio Tool
 * You can download the binary file here: https://github.com/Arapak/sio-tool
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

## FAQ

### I double-clicked the program, but it doesn't work.

The SIO Tool is a command-line tool. You should run it in the terminal.

### I cannot use `st` command.

You should put the `st` program to a path (e.g. for Linux `/usr/bin/`)
Or just google "how to add a path to the system environment variable PATH".

### How to add a new testcase

Create two extra testcase files, `inK.txt` and `outK.txt` (K is a string with 0~9).
Or a different possible name would be `$%file%$K.in` and `$%file%$K.out` where `$%file%$` is the file name of your program.

### Enable tab completion in the terminal.

Use this [Infinidat/infi.docopt_completion](https://github.com/Infinidat/infi.docopt_completion).

Note: If there is a new version released (especially a new command added), you should run `docopt-completion st` again.

### Can't parse tasks from Szkopul or Sio

You get the error: `Error: exec: "pdftotext": executable file not found in $PATH`?

This program depends on a package named (on most OSes) poppler-utils. Below is an example install on a Debian-based OS (ex. Ubuntu, Linux Mint):
```bash
sudo apt install poppler-utils
```
