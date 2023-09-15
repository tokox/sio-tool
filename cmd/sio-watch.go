package cmd

func SioWatch() (err error) {
	cln := getSioClient()
	err = cln.Ping()
	if err != nil {
		return
	}
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
