package cmd

import (
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"

	_ "modernc.org/sqlite"
)

// Open command
func SioOpen() (err error) {
	var URL string
	URL, err = Args.SioInfo.OpenURL(config.Instance.SioHost)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Sid command
func SioSid() (err error) {
	info := Args.SioInfo
	if info.SubmissionID == "" && sio_client.Instance.LastSubmission != nil {
		info = *sio_client.Instance.LastSubmission
	}
	URL, err := info.SubmissionURL(config.Instance.SioHost)
	if err != nil {
		return
	}
	return openURL(URL)
}
