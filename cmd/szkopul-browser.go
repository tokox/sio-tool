package cmd

import (
	"sio-tool/config"
	"sio-tool/szkopul_client"
)

// Open command
func SzkopulOpen() (err error) {
	URL, err := Args.SzkopulInfo.OpenURL(config.Instance.SzkopulHost)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Sid command
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
