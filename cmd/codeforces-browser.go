package cmd

import (
	"sio-tool/codeforces_client"
	"sio-tool/config"

	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
)

func openURL(url string) error {
	color.Green("Open %v", url)
	return open.Run(url)
}

// Open command
func CodeforcesOpen() (err error) {
	URL, err := Args.CodeforcesInfo.OpenURL(config.Instance.CodeforcesHost)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Stand command
func CodeforcesStand() (err error) {
	URL, err := Args.CodeforcesInfo.StandingsURL(config.Instance.CodeforcesHost)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Sid command
func CodeforcesSid() (err error) {
	info := Args.CodeforcesInfo
	if info.SubmissionID == "" && codeforces_client.Instance.LastSubmission != nil {
		info = *codeforces_client.Instance.LastSubmission
	}
	URL, err := info.SubmissionURL(config.Instance.CodeforcesHost)
	if err != nil {
		return
	}
	return openURL(URL)
}
