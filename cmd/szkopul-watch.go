package cmd

import (
	"sio-tool/szkopul_client"
)

// Watch command
func SzkopulWatch() (err error) {
	cln := szkopul_client.Instance
	info := Args.SzkopulInfo
	n := 10
	if Args.All {
		n = -1
	}
	if _, err = cln.WatchSubmission(info, n, false); err != nil {
		if err = loginAgainSzkopul(cln, err); err == nil {
			_, err = cln.WatchSubmission(info, n, false)
		}
	}
	return
}
