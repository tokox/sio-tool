package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docopt/docopt-go"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
)

// Eval opts
func Eval(opts docopt.Opts) error {
	Args = &ParsedArgs{}
	opts.Bind(Args)
	if err := parseArgs(opts); err != nil {
		return err
	}
	if Args.Config {
		return Config()
	} else if Args.Gen {
		return Gen()
	} else if Args.Test {
		return Test()
	} else if Args.StressTest {
		return StressTest()
	} else if Args.Upgrade {
		return Upgrade()
	} else if Args.Database {
		if Args.Add {
			return DatabaseAdd()
		} else if Args.Find {
			return DatabaseFind()
		} else if Args.Goto {
			return DatabaseGoto()
		}
	} else {
		if Args.Codeforces {
			if Args.Submit {
				return CodeforcesSubmit()
			} else if Args.List {
				return CodeforcesList()
			} else if Args.Parse {
				return CodeforcesParse()
			} else if Args.Watch {
				return CodeforcesWatch()
			} else if Args.Open {
				return CodeforcesOpen()
			} else if Args.Stand {
				return CodeforcesStand()
			} else if Args.Sid {
				return CodeforcesSid()
			} else if Args.Race {
				return CodeforcesRace()
			} else if Args.Pull {
				return CodeforcesPull()
			}
		} else if Args.Szkopul {
			if Args.Submit {
				return SzkopulSubmit()
			} else if Args.Watch {
				return SzkopulWatch()
			} else if Args.Parse {
				return SzkopulParse()
			} else if Args.Sid {
				return SzkopulSid()
			} else if Args.Open {
				return SzkopulOpen()
			} else if Args.List {
				return SzkopulList()
			}
		}
	}
	color.Red("This function is not available here. Maybe you are in the wrong folder?")
	return nil
}

func getSampleID() (samples []string) {
	path, err := os.Getwd()
	if err != nil {
		return
	}
	paths, err := os.ReadDir(path)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`in(\d+).txt`)
	for _, path := range paths {
		name := path.Name()
		tmp := reg.FindSubmatch([]byte(name))
		if tmp != nil {
			idx := string(tmp[1])
			ans := fmt.Sprintf("out%v.txt", idx)
			if _, err := os.Stat(ans); err == nil {
				samples = append(samples, idx)
			}
		}
	}
	return
}

func getSampleByName(filename string) (samples []string) {
	path, err := os.Getwd()
	if err != nil {
		return
	}
	paths, err := os.ReadDir(path)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(fmt.Sprintf("%s(\\d+).in", filename))
	for _, path := range paths {
		name := path.Name()
		tmp := reg.FindSubmatch([]byte(name))
		if tmp != nil {
			idx := string(tmp[1])
			ans := fmt.Sprintf("%s%v.out", filename, idx)
			if _, err := os.Stat(ans); err == nil {
				samples = append(samples, idx)
			}
		}
	}
	return
}

// CodeList Name matches some template suffix, index are template array indexes
type CodeList struct {
	Name  string
	Index []int
}

func getCode(filename string, templates []config.CodeTemplate) (codes []CodeList, err error) {
	mp := make(map[string][]int)
	for i, temp := range templates {
		suffixMap := map[string]bool{}
		for _, suffix := range temp.Suffix {
			if _, ok := suffixMap[suffix]; !ok {
				suffixMap[suffix] = true
				sf := "." + suffix
				mp[sf] = append(mp[sf], i)
			}
		}
	}

	if filename != "" {
		ext := filepath.Ext(filename)
		if idx, ok := mp[ext]; ok {
			return []CodeList{{filename, idx}}, nil
		}
		return nil, fmt.Errorf("%v can not match any template. You could add a new template by `st config`", filename)
	}

	path, err := os.Getwd()
	if err != nil {
		return
	}
	paths, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, path := range paths {
		name := path.Name()
		ext := filepath.Ext(name)
		if idx, ok := mp[ext]; ok {
			codes = append(codes, CodeList{name, idx})
		}
	}

	return codes, nil
}

func getOneCode(filename string, templates []config.CodeTemplate) (name string, index int, err error) {
	codes, err := getCode(filename, templates)
	if err != nil {
		return
	}
	if len(codes) < 1 {
		return "", 0, errors.New("cannot find any code,\nmaybe you should add a new template by `st config`")
	}
	if len(codes) > 1 {
		color.Cyan("There are multiple files can be selected.")
		for i, code := range codes {
			fmt.Printf("%3v: %v\n", i, code.Name)
		}
		i := util.ChooseIndex(len(codes))
		codes[0] = codes[i]
	}
	if len(codes[0].Index) > 1 {
		color.Cyan("There are multiple languages match the file.")
		for i, idx := range codes[0].Index {
			fmt.Printf("%3v: %v\n", i, codeforces_client.Langs[templates[idx].Lang])
		}
		i := util.ChooseIndex(len(codes[0].Index))
		codes[0].Index[0] = codes[0].Index[i]
	}
	return codes[0].Name, codes[0].Index[0], nil
}
