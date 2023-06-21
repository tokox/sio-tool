package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
)

func parseArgsSio() error {
	cfg := config.Instance
	cln := sio_client.Instance
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	if Args.Handle == "" {
		Args.Handle = cln.Username
	}
	info := sio_client.Info{}
	for _, arg := range Args.Specifier {
		parsed := parseArgSio(arg)
		if value, ok := parsed["contestID"]; ok {
			if info.Contest != "" && info.Contest != value {
				return fmt.Errorf("contest ID conflicts: %v %v", info.Contest, value)
			}
			info.Contest = value
		}
		if value, ok := parsed["problemAlias"]; ok {
			if info.ProblemAlias != "" && info.ProblemAlias != value {
				return fmt.Errorf("problem alias conflicts: %v %v", info.ProblemAlias, value)
			}
			info.ProblemAlias = value
		}
		if value, ok := parsed["problemID"]; ok {
			if info.ProblemID != "" && info.ProblemID != value {
				return fmt.Errorf("problemID conflicts: %v %v", info.ProblemID, value)
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
	parsedPath := parsePathSio(path)
	if info.Contest == "" {
		if value, ok := parsedPath["contestID"]; ok {
			info.Contest = value
		}
		if info.Round == "" {
			if value, ok := parsedPath["round"]; ok {
				info.Round = value
			}
			if info.ProblemAlias == "" {
				if value, ok := parsedPath["problemAlias"]; ok {
					info.ProblemAlias = value
				}
			}
		}
	}
	info.RootPath = filepath.Join(cfg.FolderName["sio-root"])
	Args.SioInfo = info
	return nil
}

const SioProblemRegStr = `[a-z]+\d*`
const SioProblemIdRegStr = `\d+`
const SioContestRegStr = `[\w-]+?`
const SioRoundRegStr = `\w+?`

var SioArgRegStr = [...]string{
	`^[sS][iI][oO]?$`,
	fmt.Sprintf(`/c/(?P<contestID>%v)/(p/(?P<problemAlias>%v)?)?`, SioContestRegStr, SioProblemRegStr),
	fmt.Sprintf(`^(?P<problemID>%v)$`, SioProblemIdRegStr),
	fmt.Sprintf(`^(?P<problemAlias>%v)$`, SioProblemRegStr),
	fmt.Sprintf(`^(?P<contestID>%v)$`, SioContestRegStr),
}

func parseArgSio(arg string) map[string]string {
	output := make(map[string]string)
	for _, regStr := range SioArgRegStr {
		reg := regexp.MustCompile(regStr)
		names := reg.SubexpNames()
		found := false
		for i, val := range reg.FindStringSubmatch(arg) {
			if names[i] != "" && val != "" {
				output[names[i]] = val
				found = true
			}
		}
		if found {
			break
		}
	}
	return output
}

var SioPathRegStr = fmt.Sprintf("%v/((?P<contestID>%v)/((?P<round>%v)/((?P<problemAlias>%v)/)?)?)?", "%v", SioContestRegStr, SioRoundRegStr, SioProblemRegStr)

func parsePathSio(path string) map[string]string {
	path = filepath.ToSlash(path) + "/"
	output := make(map[string]string)
	cfg := config.Instance
	reg := regexp.MustCompile(fmt.Sprintf(SioPathRegStr, cfg.FolderName["sio-root"]))
	names := reg.SubexpNames()
	for i, val := range reg.FindStringSubmatch(path) {
		if names[i] != "" && val != "" {
			output[names[i]] = val
		}
	}
	return output
}
