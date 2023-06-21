package cmd

import (
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
)

func SioRace() (err error) {
	cfg := config.Instance
	cln := sio_client.Instance
	info := Args.SioInfo
	if Args.SioInfo.Round, err = cln.RaceContest(info); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			Args.SioInfo.Round, err = cln.RaceContest(info)
		}
	}
	if err != nil {
		return
	}
	URL, err := info.ContestURL(cfg.SioHost)
	if err != nil {
		return
	}
	openURL(URL)
	return SioParse()
}
