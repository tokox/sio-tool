package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"sio-tool/codeforces_client"
	"sio-tool/config"

	"github.com/docopt/docopt-go"
)

func parseArgsCodeforces(opts docopt.Opts) error {
	cfg := config.Instance
	cln := codeforces_client.Instance
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	if Args.Handle == "" {
		Args.Handle = cln.Handle
	}
	info := codeforces_client.Info{}
	for _, arg := range Args.Specifier {
		parsed := parseArgCodeforces(arg)
		if value, ok := parsed["problemType"]; ok {
			if info.ProblemType != "" && info.ProblemType != value {
				return fmt.Errorf("problem type conflicts: %v %v", info.ProblemType, value)
			}
			info.ProblemType = value
		}
		if value, ok := parsed["contestID"]; ok {
			if info.ContestID != "" && info.ContestID != value {
				return fmt.Errorf("contest ID conflicts: %v %v", info.ContestID, value)
			}
			info.ContestID = value
		}
		if value, ok := parsed["groupID"]; ok {
			if info.GroupID != "" && info.GroupID != value {
				return fmt.Errorf("group ID conflicts: %v %v", info.GroupID, value)
			}
			info.GroupID = value
		}
		if value, ok := parsed["problemID"]; ok {
			if info.ProblemID != "" && info.ProblemID != value {
				return fmt.Errorf("problem ID conflicts: %v %v", info.ProblemID, value)
			}
			info.ProblemID = value
		}
		if value, ok := parsed["submissionID"]; ok {
			if info.SubmissionID != "" && info.SubmissionID != value {
				return fmt.Errorf("submission ID conflicts: %v %v", info.SubmissionID, value)
			}
			info.SubmissionID = value
		}
	}
	if info.ProblemType == "" {
		parsed := parsePathCodeforces(path)
		if value, ok := parsed["problemType"]; ok {
			info.ProblemType = value
		}
		if value, ok := parsed["contestID"]; ok && info.ContestID == "" {
			info.ContestID = value
		}
		if value, ok := parsed["groupID"]; ok && info.GroupID == "" {
			info.GroupID = value
		}
		if value, ok := parsed["problemID"]; ok && info.ProblemID == "" {
			info.ProblemID = value
		}
	}
	if info.ProblemType == "" || info.ProblemType == "contest" {
		if len(info.ContestID) < 6 {
			info.ProblemType = "contest"
		} else {
			info.ProblemType = "gym"
		}
	}
	if info.ProblemType == "acmsguru" {
		if info.ContestID != "99999" && info.ContestID != "" {
			info.ProblemID = info.ContestID
		}
		info.ContestID = "99999"
	}
	root := cfg.FolderName["codeforces-root"]
	info.RootPath = filepath.Join(root, cfg.FolderName[info.ProblemType])
	Args.CodeforcesInfo = info
	return nil
}

const CodeforcesProblemRegStr = `\w+`

const CodeforcesStrictProblemRegStr = `[a-zA-Z]+\d*`

const CodeforcesContestRegStr = `\d+`

const CodeforcesGroupRegStr = `\w{10}`

const CodeforcesSubmissionRegStr = `\d+`

var CodeforcesArgRegStr = [...]string{
	`^[cC][oO][nN][tT][eE][sS][tT][sS]?$`,
	`^[gG][yY][mM][sS]?$`,
	`^[gG][rR][oO][uU][pP][sS]?$`,
	`^[aA][cC][mM][sS][gG][uU][rR][uU]$`,
	fmt.Sprintf(`/contest/(?P<contestID>%v)(/problem/(?P<problemID>%v))?`, CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf(`/gym/(?P<contestID>%v)(/problem/(?P<problemID>%v))?`, CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf(`/problemset/problem/(?P<contestID>%v)/(?P<problemID>%v)`, CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf(`/group/(?P<groupID>%v)(/contest/(?P<contestID>%v)(/problem/(?P<problemID>%v))?)?`, CodeforcesGroupRegStr, CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf(`/problemsets/acmsguru/problem/(?P<contestID>%v)/(?P<problemID>%v)`, CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf(`/problemsets/acmsguru/submission/(?P<contestID>%v)/(?P<submissionID>%v)`, CodeforcesContestRegStr, CodeforcesSubmissionRegStr),
	fmt.Sprintf(`/submission/(?P<submissionID>%v)`, CodeforcesSubmissionRegStr),
	fmt.Sprintf(`^(?P<contestID>%v)(?P<problemID>%v)$`, CodeforcesContestRegStr, CodeforcesStrictProblemRegStr),
	fmt.Sprintf(`^(?P<contestID>%v)$`, CodeforcesContestRegStr),
	fmt.Sprintf(`^(?P<problemID>%v)$`, CodeforcesStrictProblemRegStr),
	fmt.Sprintf(`^(?P<groupID>%v)$`, CodeforcesGroupRegStr),
}

var CodeforcesArgTypePathRegStr = [...]string{
	fmt.Sprintf("%v/%v/((?P<contestID>%v)/((?P<problemID>%v)/)?)?", "%v", "%v", CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf("%v/%v/((?P<contestID>%v)/((?P<problemID>%v)/)?)?", "%v", "%v", CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf("%v/%v/((?P<groupID>%v)/((?P<contestID>%v)/((?P<problemID>%v)/)?)?)?", "%v", "%v", CodeforcesGroupRegStr, CodeforcesContestRegStr, CodeforcesProblemRegStr),
	fmt.Sprintf("%v/%v/((?P<problemID>%v)/)?", "%v", "%v", CodeforcesProblemRegStr),
}

var CodeforcesArgType = [...]string{
	"contest",
	"gym",
	"group",
	"acmsguru",
	"contest",
	"gym",
	"contest",
	"group",
	"acmsguru",
	"acmsguru",
	"",
	"",
	"",
	"",
	"",
}

func parseArgCodeforces(arg string) map[string]string {
	output := make(map[string]string)
	for k, regStr := range CodeforcesArgRegStr {
		reg := regexp.MustCompile(regStr)
		names := reg.SubexpNames()
		for i, val := range reg.FindStringSubmatch(arg) {
			if names[i] != "" && val != "" {
				output[names[i]] = val
			}
			if CodeforcesArgType[k] != "" {
				output["problemType"] = CodeforcesArgType[k]
				if k < 4 {
					return output
				}
			}
		}
	}
	return output
}

func parsePathCodeforces(path string) map[string]string {
	path = filepath.ToSlash(path) + "/"
	output := make(map[string]string)
	cfg := config.Instance
	for k, problemType := range codeforces_client.ProblemTypes {
		reg := regexp.MustCompile(fmt.Sprintf(CodeforcesArgTypePathRegStr[k], cfg.FolderName["codeforces-root"], cfg.FolderName[problemType]))
		names := reg.SubexpNames()
		for i, val := range reg.FindStringSubmatch(path) {
			if names[i] != "" && val != "" {
				output[names[i]] = val
			}
			output["problemType"] = problemType
		}
	}
	if (output["problemType"] != "" && output["problemType"] != "group") ||
		output["groupID"] == output["contestID"] ||
		output["groupID"] == fmt.Sprintf("%v%v", output["contestID"], output["problemID"]) {
		output["groupID"] = ""
	}
	if output["groupID"] != "" && output["problemType"] == "" {
		output["problemType"] = "group"
	}
	return output
}
