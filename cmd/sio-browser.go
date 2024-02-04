package cmd

import (
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"

	_ "modernc.org/sqlite"
)

func SioOpen() (err error) {
	var URL string
	URL, err = Args.SioInfo.OpenURL(getSioHost())
	if err != nil {
		return
	}
	return openURL(URL)
}

func SioSid() (err error) {
	info := Args.SioInfo
	cln := getSioClient()
	if info.SubmissionID == "" && cln.LastSubmission != nil {
		info = *cln.LastSubmission
	}
	URL, err := info.SubmissionURL(getSioHost(), false)
	if err != nil {
		return
	}
	return openURL(URL)
}

func SioStand() (err error) {
	URL, err := Args.SioInfo.StandingsURL(getSioClient(), getSioHost())
	if err != nil {
		return
	}
	return openURL(URL)
}

func getSioClient() *sio_client.SioClient {
	if Args.SioStaszic {
		return sio_client.StaszicInstance
	} else if Args.SioMimuw {
		return sio_client.MimuwInstance
	} else if Args.SioTalent {
		return sio_client.TalentInstance
	}
	return nil
}

func getSioHost() string {
	if Args.SioStaszic {
		return config.Instance.SioStaszicHost
	} else if Args.SioMimuw {
		return config.Instance.SioMimuwHost
	} else if Args.SioTalent {
		return config.Instance.SioTalentHost
	}
	return ""
}
