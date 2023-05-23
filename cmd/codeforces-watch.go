package cmd

import (
	"sio-tool/codeforces_client"
)

// Watch command
func CodeforcesWatch() (err error) {
	cln := codeforces_client.Instance
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