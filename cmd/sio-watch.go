package cmd

import "github.com/Arapak/sio-tool/sio_client"

// Watch command
func SioWatch() (err error) {
	cln := sio_client.Instance
	info := Args.SioInfo
	n := 10
	if Args.All {
		n = -1
	}
	if _, err = cln.WatchSubmission(info, n, false); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			_, err = cln.WatchSubmission(info, n, false)
		}
	}
	return
}
