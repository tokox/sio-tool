package szkopul_client

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Arapak/sio-tool/database_client"
)

var Archives = [...]string{
	"OI",
}

// Info information
type Info struct {
	Archive      string `json:"archive"`
	ContestID    string `json:"contest_id"`
	StageID      string `json:"stage_id"`
	ProblemAlias string `json:"problem_alias"`
	ProblemID    string `json:"problem_id"`
	SubmissionID string `json:"submission_id"`
	RootPath     string
}

// ErrorNeedProblemID error
const ErrorNeedProblemID = "you have to specify the Problem ID"

// ErrorNeedProblemID error
const ErrorNeedProblemAlias = "you have to specify the Problem alias"

// ErrorNeedContestID error
const ErrorNeedContestID = "you have to specify the Contest ID"

// ErrorNeedSubmissionID error
const ErrorNeedSubmissionID = "you have to specify the Submission ID"

// ErrorNeedSubmissionID error
const ErrorNeedArchive = "you have to specify the archive"

// Hint hint text
func (info *Info) Hint() string {
	text := ""
	if info.ContestID != "" {
		text = "contest " + info.ContestID
	}
	if info.StageID != "" {
		text = text + ", stage " + info.StageID
	}
	if info.ProblemAlias != "" {
		text = text + ", problem " + strings.ToLower(info.ProblemAlias)
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
			if info.ProblemAlias != "" {
				path = filepath.Join(path, strings.ToLower(info.ProblemAlias))
			}
		}
	}
	return path
}

func ProblemURL(host, problemID string) string {
	return fmt.Sprintf(host+"/problemset/problem/%v/site", problemID)
}

// ProblemURL parse problem url
func (info *Info) ProblemURL(host string) (string, error) {
	if info.ProblemID == "" {
		return "", errors.New(ErrorNeedProblemID)
	}
	return ProblemURL(host, info.ProblemID), nil
}

func (info *Info) ProblemSetURL(host string) (string, error) {
	if info.Archive == "" {
		return "", errors.New(ErrorNeedArchive)
	}
	return fmt.Sprintf(host+"/task_archive/%v/", strings.ToLower(info.Archive)), nil
}

// MySubmissionURL parse submission url
func (info *Info) MySubmissionURL(host string) string {
	if info.ProblemID == "" {
		return host + "/submissions/"
	} else {
		return fmt.Sprintf(host+"/problemset/problem/%v/site/?key=submissions", info.ProblemID)
	}
}

// SubmissionURL parse submission url
func (info *Info) SubmissionURL(host string) (string, error) {
	if info.SubmissionID == "" {
		return "", errors.New(ErrorNeedSubmissionID)
	}
	return fmt.Sprintf(host+"/s/%v", info.SubmissionID), nil
}

// APISubmitURL submit url
func (info *Info) APISubmitURL(host string) (string, error) {
	if info.ProblemID == "" {
		return "", errors.New(ErrorNeedProblemID)
	}
	return fmt.Sprintf(host+"/api/problemset/submit/%v", info.ProblemID), nil
}

// SubmitURL submit url
func (info *Info) SubmitURL(host string) (string, error) {
	if info.ProblemID == "" {
		return "", errors.New(ErrorNeedProblemID)
	}
	return fmt.Sprintf(host+"/problemset/problem/%v/site/?key=submit", info.ProblemID), nil
}

// OpenURL open url
func (info *Info) OpenURL(host string) (string, error) {
	if info.ProblemID != "" {
		return info.ProblemURL(host)
	}
	return host + "/task_archive/oi", nil
}

func (info *Info) ToTask() database_client.Task {
	return database_client.Task{ShortName: info.ProblemAlias, Source: info.Archive, ContestID: info.ContestID, ContestStageID: info.StageID}
}
