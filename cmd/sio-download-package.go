package cmd

import (
	"os"
)

func SioDownloadPackages() (err error) {
	cln := getSioClient()
	err = cln.Ping()
	if err != nil {
		return
	}
	info := Args.SioInfo
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	if _, err = cln.DownloadAllPackages(info, rootPath); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			_, err = cln.DownloadAllPackages(info, rootPath)
		}
	}
	return
}
