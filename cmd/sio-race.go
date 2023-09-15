package cmd

func SioRace() (err error) {
	cln := getSioClient()
	err = cln.Ping()
	if err != nil {
		return
	}
	info := Args.SioInfo
	if Args.SioInfo.Round, err = cln.RaceContest(info); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			Args.SioInfo.Round, err = cln.RaceContest(info)
		}
	}
	if err != nil {
		return
	}
	URL, err := info.ContestURL(getSioHost())
	if err != nil {
		return
	}
	err = openURL(URL)
	if err != nil {
		return err
	}
	return SioParse()
}
