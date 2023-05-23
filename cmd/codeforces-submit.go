package cmd

import (
	"os"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"

	"github.com/fatih/color"
)

// Submit command
func CodeforcesSubmit() (err error) {
	cln := codeforces_client.Instance
	cfg := config.Instance
	info := Args.CodeforcesInfo
	filename, index, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	source := string(bytes)

	lang := cfg.Template[index].Lang
	if err = cln.Submit(info, lang, source); err != nil {
		if err = loginAgainCodeforces(cln, err); err == nil {
			err = cln.Submit(info, lang, source)
		}
	}
	return
}

func loginAgainCodeforces(cln *codeforces_client.CodeforcesClient, err error) error {
	if err != nil && err.Error() == codeforces_client.ErrorNotLogged {
		color.Red("Not logged. Try to login\n")
		err = cln.Login()
	}
	return err
}
