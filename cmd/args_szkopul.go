package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"sio-tool/config"
	"sio-tool/szkopul_client"

	"github.com/docopt/docopt-go"
)

func parseArgsSzkopul(opts docopt.Opts) error {
	cfg := config.Instance
	cln := szkopul_client.Instance
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	if Args.Handle == "" {
		Args.Handle = cln.Username
	}
	info := szkopul_client.Info{}
	for _, arg := range Args.Specifier {
		parsed := parseArgSzkopul(arg)
		if value, ok := parsed["archive"]; ok {
			if info.Archive != "" && info.Archive != value {
				return fmt.Errorf("archive conflicts: %v %v", info.Archive, value)
			}
			info.Archive = value
		}
		if value, ok := parsed["contestID"]; ok {
			if info.ContestID != "" && info.ContestID != value {
				return fmt.Errorf("contest ID conflicts: %v %v", info.ContestID, value)
			}
			info.ContestID = value
		}
		if value, ok := parsed["stageID"]; ok {
			if info.StageID != "" && info.StageID != value {
				return fmt.Errorf("group ID conflicts: %v %v", info.StageID, value)
			}
			info.StageID = value
		}
		if value, ok := parsed["problemID"]; ok {
			if info.ProblemID != "" && info.ProblemID != value {
				return fmt.Errorf("problem ID conflicts: %v %v", info.ProblemID, value)
			}
			info.ProblemID = value
		}
		if value, ok := parsed["problemSecretKey"]; ok {
			if info.ProblemSecretKey != "" && info.ProblemSecretKey != value {
				return fmt.Errorf("problemSecretKey conflicts: %v %v", info.ProblemSecretKey, value)
			}
			info.ProblemSecretKey = value
		}
		if value, ok := parsed["submissionID"]; ok {
			if info.SubmissionID != "" && info.SubmissionID != value {
				return fmt.Errorf("submission ID conflicts: %v %v", info.SubmissionID, value)
			}
			info.SubmissionID = value
		}
	}
	if info.Archive == "" {
		parsed := parsePathSzkopul(path)
		if value, ok := parsed["archive"]; ok {
			info.Archive = value
		}
		if value, ok := parsed["contestID"]; ok && info.ContestID == "" {
			info.ContestID = value
		}
		if value, ok := parsed["stageID"]; ok && info.StageID == "" {
			info.StageID = value
		}
		if value, ok := parsed["problemID"]; ok && info.ProblemID == "" {
			info.ProblemID = value
		}
	}
	// util.DebugJSON(info)
	info.RootPath = filepath.Join(cfg.FolderName["szkopul-root"], cfg.FolderName[fmt.Sprintf("codeforces-%v", info.Archive)])
	Args.SzkopulInfo = info
	return nil
}

const SzkopulProblemRegStr = `\w+`
const SzkopulProblemSecretKeyRegStr = `[A-Za-z0-9]{24}`

const StrictSzkopulProblemRegStr = `[a-z]{3}\d*`
const OIContestRegStr = `[MCLXVI]+`
const OIStageRegStr = `[1-3]`

var SzkopulArgRegStr = [...]string{
	`^[oO][iI]?$`,
	fmt.Sprintf(`/problemset/problem/(?P<problemSecretKey>%v)(/site(/\?key=\w+)?)?`, SzkopulProblemSecretKeyRegStr),
	fmt.Sprintf(`^(?P<problemSecretKey>%v)$`, SzkopulProblemSecretKeyRegStr),
	fmt.Sprintf(`^(?P<problemID>%v)$`, StrictSzkopulProblemRegStr),
	fmt.Sprintf(`^(?P<contestID>%v)$`, OIContestRegStr),
	fmt.Sprintf(`^(?P<stageID>%v)$`, OIStageRegStr),
}

var SzkopulArgType = [...]string{
	"OI",
	"OI",
	"OI",
	"OI",
	"OI",
	"OI",
}

func parseArgSzkopul(arg string) map[string]string {
	output := make(map[string]string)
	for k, regStr := range SzkopulArgRegStr {
		reg := regexp.MustCompile(regStr)
		names := reg.SubexpNames()
		for i, val := range reg.FindStringSubmatch(arg) {
			if names[i] != "" && val != "" {
				output[names[i]] = val
			}
			if SzkopulArgType[k] != "" {
				output["archive"] = SzkopulArgType[k]
			}
		}
	}
	return output
}

var SzkopulPathRegStr = [...]string{
	fmt.Sprintf("%v/%v/((?P<contestID>%v)/((?P<stageID>%v)/((?P<problemID>%v)/)?)?)?", "%v", "%v", OIContestRegStr, OIStageRegStr, StrictSzkopulProblemRegStr),
}

func parsePathSzkopul(path string) map[string]string {
	path = filepath.ToSlash(path) + "/"
	output := make(map[string]string)
	cfg := config.Instance
	for k, archive := range szkopul_client.Archives {
		reg := regexp.MustCompile(fmt.Sprintf(SzkopulPathRegStr[k], cfg.FolderName["szkopul-root"], cfg.FolderName[fmt.Sprintf("codeforces-%v", archive)]))
		names := reg.SubexpNames()
		for i, val := range reg.FindStringSubmatch(path) {
			if names[i] != "" && val != "" {
				output[names[i]] = val
			}
			output["archive"] = archive
		}
	}
	return output
}
