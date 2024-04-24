package cmd

import (
	"os"
	"path"
)

func SioUploadPackage() (err error) {
	cln := getSioClient()
	err = cln.Ping()
	if err != nil {
		return
	}
	info := Args.SioInfo
	file := Args.File

	var rootPath string
	if !path.IsAbs(file) {
		rootPath, err = os.Getwd()
		if err != nil {
			return
		}
		file = path.Join(rootPath, file)
	}

	if _, err = cln.UploadPackage(info, file); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			_, err = cln.UploadPackage(info, file)
		}
	}
	return
}
