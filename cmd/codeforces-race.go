package cmd

import (
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
)

func CodeforcesRace() (err error) {
	cfg := config.Instance
	cln := codeforces_client.Instance
	info := Args.CodeforcesInfo
	if err = cln.RaceContest(info); err != nil {
		if err = loginAgainCodeforces(cln, err); err == nil {
			err = cln.RaceContest(info)
		}
	}
	if err != nil {
		return
	}
	URL, err := info.ProblemSetURL(cfg.CodeforcesHost)
	if err != nil {
		return
	}
	err = openURL(URL)
	if err != nil {
		return err
	}
	err = openURL(URL + "/problems")
	if err != nil {
		return err
	}
	return CodeforcesParse()
}
