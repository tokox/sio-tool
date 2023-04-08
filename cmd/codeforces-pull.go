package cmd

import (
	"os"

	"sio-tool/codeforces_client"
)

// Pull command
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
