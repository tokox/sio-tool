package sio_client

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Arapak/sio-tool/database_client"
)

type Info struct {
	Contest      string `json:"contest_id"`
	ProblemID    string `json:"problem_id"`
	ProblemAlias string `json:"problem_alias"`
	Round        string `json:"round"`
	SubmissionID string `json:"submission_id"`
	RootPath     string
}

const ErrorNeedProblemAlias = "you have to specify the Problem alias"

const ErrorNeedContest = "you have to specify the contest"

const ErrorNeedSubmissionID = "you have to specify the Submission ID"

func (info *Info) Hint() string {
	text := ""
	if info.Contest != "" {
		text = "contest " + info.Contest
		if info.Round != "" {
			text = text + ", round " + strings.ToLower(info.Round)
		}
		if info.ProblemAlias != "" {
			text = text + ", problem " + strings.ToLower(info.ProblemAlias)
		}
	}
	if info.SubmissionID != "" {
		if text != "" {
			text += ", "
		}
		text = text + "submission " + info.SubmissionID
	}
	return text
}

var polishCharMap = map[rune]rune{
	'ą': 'a',
	'ć': 'c',
	'ę': 'e',
	'ł': 'l',
	'ń': 'n',
	'ó': 'o',
	'ś': 's',
	'ź': 'z',
	'ż': 'z',
}

func replacePolishChars(str string) string {
	runes := []rune(str)
	for i := range runes {
		if val, ok := polishCharMap[runes[i]]; ok == true {
			runes[i] = val
		}
	}
	return string(runes)
}

var nonAlphanumericRegExp = regexp.MustCompile(`\W+`)
var underscoreRegex = regexp.MustCompile("_+")

func clearString(str string) string {
	str = strings.ToLower(str)
	str = replacePolishChars(str)
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, ".", "_")
	str = strings.ReplaceAll(str, "-", "_")

	str = nonAlphanumericRegExp.ReplaceAllString(str, "")
	return underscoreRegex.ReplaceAllString(str, "_")
}

func (info *Info) Path() string {
	path := info.RootPath
	if info.Contest != "" {
		path = filepath.Join(path, info.Contest)
		if info.Round != "" {
			path = filepath.Join(path, clearString(info.Round))
			if info.ProblemAlias != "" {
				path = filepath.Join(path, strings.ToLower(info.ProblemAlias))
			}
		}
	}
	return path
}

func ProblemURL(host, contest, problemAlias string) string {
	return fmt.Sprintf(host+"/c/%v/p/%v", contest, problemAlias)
}

func (info *Info) ProblemURL(host string) (string, error) {
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	if info.ProblemAlias == "" {
		return "", errors.New(ErrorNeedProblemAlias)
	}
	return ProblemURL(host, info.Contest, info.ProblemAlias), nil
}

func (info *Info) ContestURL(host string) (string, error) {
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	return fmt.Sprintf(host+"/c/%v/p", info.Contest), nil
}

func (info *Info) MySubmissionURL(host string) (string, error) {
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	return fmt.Sprintf(host+"/c/%v/submissions/", info.Contest), nil
}

func (info *Info) SubmissionURL(host string) (string, error) {
	if info.SubmissionID == "" {
		return "", errors.New(ErrorNeedSubmissionID)
	}
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	return fmt.Sprintf(host+"/c/%v/s/%v", info.Contest, info.SubmissionID), nil
}

func (info *Info) SubmitURL(host string) (string, error) {
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	return fmt.Sprintf(host+"/c/%v/submit/", info.Contest), nil
}

func (info *Info) StandingsURL(host string) (string, error) {
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	return fmt.Sprintf(host+"/c/%v/r/", info.Contest), nil
}

func (info *Info) StatusURL(host string) (string, error) {
	if info.Contest == "" {
		return "", errors.New(ErrorNeedContest)
	}
	return fmt.Sprintf(host+"/c/%v/status", info.Contest), nil
}

func (info *Info) OpenURL(host string) (string, error) {
	if info.Contest != "" {
		if info.ProblemAlias != "" {
			return info.ProblemURL(host)
		}
		return info.ContestURL(host)
	}
	return host, nil
}

func (info *Info) ToTask() database_client.Task {
	return database_client.Task{ShortName: info.ProblemAlias, Source: "sio", ContestID: info.Contest}
}
