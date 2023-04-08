package szkopul_client

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var Archives = [...]string{
	"OI",
}

// Info information
type Info struct {
	Archive          string `json:"archive"`
	ContestID        string `json:"contest_id"`
	StageID          string `json:"stage_id"`
	ProblemID        string `json:"problem_id"`
	ProblemSecretKey string `json:"problem_secret_key"`
	SubmissionID     string `json:"submission_id"`
	RootPath         string
}

// ErrorNeedProblemID error
const ErrorNeedProblemID = "you have to specify the Problem ID"

// ErrorNeedProblemSecretKey error
const ErrorNeedProblemSecretKey = "you have to specify the Problem Secret Key"

// ErrorNeedContestID error
const ErrorNeedContestID = "you have to specify the Contest ID"

// ErrorNeedSubmissionID error
const ErrorNeedSubmissionID = "you have to specify the Submission ID"

// Hint hint text
func (info *Info) Hint() string {
	text := ""
	if info.ContestID != "" {
		text = "contest " + info.ContestID
		if info.StageID != "" {
			text = text + ", stage " + info.StageID
			if info.ProblemID != "" {
				text = text + ", problem " + strings.ToLower(info.ProblemID)
			}
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

// Path path
func (info *Info) Path() string {
	path := info.RootPath
	if info.ContestID != "" {
		path = filepath.Join(path, info.ContestID)
		if info.StageID != "" {
			path = filepath.Join(path, info.StageID)
			if info.ProblemID != "" {
				path = filepath.Join(path, strings.ToLower(info.ProblemID))
			}
		}
	}
	return path
}

// ProblemURL parse problem url
func (info *Info) ProblemURL(host string) (string, error) {
	if info.ProblemSecretKey == "" {
		return "", errors.New(ErrorNeedProblemSecretKey)
	}
	return fmt.Sprintf(host+"/problemset/problem/%v/site", info.ProblemSecretKey), nil
}

// MySubmissionURL parse submission url
func (info *Info) MySubmissionURL(host string) string {
	if info.ProblemSecretKey == "" {
		return host + "/submissions/"
	} else {
		return fmt.Sprintf(host+"/problemset/problem/%v/site/?key=submissions", info.ProblemSecretKey)
	}
}

// SubmissionURL parse submission url
func (info *Info) SubmissionURL(host string) (string, error) {
	if info.SubmissionID == "" {
		return "", errors.New(ErrorNeedSubmissionID)
	}
	return fmt.Sprintf(host+"/s/%v", info.SubmissionID), nil
}

// SubmitURL submit url
func (info *Info) SubmitURL(host string) (string, error) {
	if info.ProblemSecretKey == "" {
		return "", errors.New(ErrorNeedProblemSecretKey)
	}
	return fmt.Sprintf(host+"/api/problemset/submit/%v", info.ProblemSecretKey), nil
}

// OpenURL open url
func (info *Info) OpenURL(host string) (string, error) {
	if info.ProblemSecretKey != "" {
		return info.ProblemURL(host)
	}
	return host + "/task_archive/oi", nil
}
