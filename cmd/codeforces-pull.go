package cmd

import (
	"os"

	"github.com/Arapak/sio-tool/codeforces_client"
)

func CodeforcesPull() (err error) {
	cln := codeforces_client.Instance
	info := Args.CodeforcesInfo
	ac := Args.Accepted
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	if err = cln.Pull(info, rootPath, ac); err != nil {
		if err = loginAgainCodeforces(cln, err); err == nil {
			err = cln.Pull(info, rootPath, ac)
		}
	}
	return
}
