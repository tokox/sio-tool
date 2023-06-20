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

var nonAlphanumericRegExp = regexp.MustCompile(`\W+`)

func clearString(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, " ", "_", -1)
	return nonAlphanumericRegExp.ReplaceAllString(str, "")
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
