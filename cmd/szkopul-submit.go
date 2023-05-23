package cmd

import (
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/szkopul_client"

	"github.com/fatih/color"
)

// Submit command
func SzkopulSubmit() (err error) {
	cln := szkopul_client.Instance
	cfg := config.Instance
	info := Args.SzkopulInfo
	filename, _, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}

	// lang := cfg.Template[index].Lang
	if err = cln.Submit(info, filename); err != nil {
		if err = loginAgainSzkopul(cln, err); err == nil {
			err = cln.Submit(info, filename)
		}
	}
	return
}

func loginAgainSzkopul(cln *szkopul_client.SzkopulClient, err error) error {
	if err != nil && err.Error() == szkopul_client.ErrorNotLogged {
		color.Red("Not logged. Try to login\n")
		err = cln.Login()
	}
	return err
}
