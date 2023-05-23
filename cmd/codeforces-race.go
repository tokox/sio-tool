package cmd

import (
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
)

// Race command
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
	openURL(URL)
	openURL(URL + "/problems")
	return CodeforcesParse()
}
