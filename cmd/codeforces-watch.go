package cmd

import (
	"github.com/Arapak/sio-tool/codeforces_client"
)

func CodeforcesWatch() (err error) {
	cln := codeforces_client.Instance
	err = cln.Ping()
	if err != nil {
		return
	}
	info := Args.CodeforcesInfo
	n := 10
	if Args.All {
		n = -1
	}
	if _, err = cln.WatchSubmission(info, n, false); err != nil {
		if err = loginAgainCodeforces(cln, err); err == nil {
			_, err = cln.WatchSubmission(info, n, false)
		}
	}
	return
}
