package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/database_client"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"
	"github.com/fatih/color"

	_ "modernc.org/sqlite"
)

const ErrorNoProblemFound = "no problem found matching criteria"
const ErrorMultipleProblemsFound = "more than one problem found matching criteria"

func getLinkToProblemFromDatabase(task database_client.Task) (link string, err error) {
	cfg := config.Instance
	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		fmt.Printf("failed to open database connection: %v\n", err)
		return
	}
	defer db.Close()
	tasks, err := database_client.FindTasks(db, task)
	if err != nil {
		return
	}
	if len(tasks) == 0 {
		err = errors.New(ErrorNoProblemFound)
	} else if len(tasks) > 1 {
		err = errors.New(ErrorMultipleProblemsFound)
	} else {
		link = tasks[0].Link
	}
	return
}

func getLinkToProblemFromStatis() (link string, err error) {
	cln := szkopul_client.Instance
	var problems []szkopul_client.StatisInfo
	var perf util.Performance
	color.Green("Fetching...")
	problems, perf, err = cln.Statis(Args.SzkopulInfo)
	if err != nil {
		return
	}
	fmt.Printf("Statis: (%v)\n", perf.Parse())
	if len(problems) == 0 {
		err = errors.New(ErrorNoProblemFound)
	} else if len(problems) > 1 {
		err = errors.New(ErrorMultipleProblemsFound)
	} else {
		link = szkopul_client.ProblemURL(config.Instance.SzkopulHost, problems[0].ID)
	}
	return
}

func SzkopulOpen() (err error) {
	var URL string
	if Args.SzkopulInfo.ProblemID == "" && Args.SzkopulInfo.ProblemAlias != "" {
		if util.Confirm("You didn't specify the problemID, but have given a problem alias, do you want to search for the problem with given criteria (Y/n):") {
			color.Green("Searching in database...")
			URL, err = getLinkToProblemFromDatabase(Args.SzkopulInfo.ToTask())
			if err == nil {
				return openURL(URL)
			}
			if err.Error() == ErrorNoProblemFound || err.Error() == ErrorMultipleProblemsFound {
				color.Red(err.Error())
				URL, err = getLinkToProblemFromStatis()
				if err == nil {
					return openURL(URL)
				}
			}
			return
		}
	}
	URL, err = Args.SzkopulInfo.OpenURL(config.Instance.SzkopulHost)
	if err != nil {
		return
	}
	return openURL(URL)
}

func SzkopulSid() (err error) {
	info := Args.SzkopulInfo
	if info.SubmissionID == "" && szkopul_client.Instance.LastSubmission != nil {
		info = *szkopul_client.Instance.LastSubmission
	}
	URL, err := info.SubmissionURL(config.Instance.SzkopulHost)
	if err != nil {
		return
	}
	return openURL(URL)
}
