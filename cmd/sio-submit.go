package cmd

import (
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"

	"github.com/fatih/color"
)

// Submit command
func SioSubmit() (err error) {
	cln := sio_client.Instance
	cfg := config.Instance
	info := Args.SioInfo
	filename, _, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}

	if err = cln.Submit(info, filename); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			err = cln.Submit(info, filename)
		}
	}
	return
}

func loginAgainSio(cln *sio_client.SioClient, err error) error {
	if err != nil && err.Error() == sio_client.ErrorNotLogged {
		color.Red("Not logged. Try to login\n")
		err = cln.Login()
	}
	return err
}
