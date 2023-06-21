package cmd

import (
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"

	_ "modernc.org/sqlite"
)

func SioOpen() (err error) {
	var URL string
	URL, err = Args.SioInfo.OpenURL(config.Instance.SioHost)
	if err != nil {
		return
	}
	return openURL(URL)
}

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

func SioStand() (err error) {
	URL, err := Args.SioInfo.StandingsURL(config.Instance.SioHost)
	if err != nil {
		return
	}
	return openURL(URL)
}
