package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
	"github.com/Arapak/sio-tool/szkopul_client"
)

func Config() (err error) {
	cfg := config.Instance
	codeforcesCln := codeforces_client.Instance
	szkopulCln := szkopul_client.Instance
	sioCln := sio_client.Instance

	index := 0
	prompt := &survey.Select{
		Message: "Configure the tool",
		Options: []string{
			`login`,
			`add a template`,
			`delete a template`,
			`set default template`,
			`run "st gen" after "st parse"`,
			`set codeforces host domain`,
			`set proxy`,
			`set folders' name`,
			`set default naming`,
			`set database path`,
		},
		PageSize: 10,
	}
	if err = survey.AskOne(prompt, &index); err != nil {
		return
	}
	if index == 0 {
		prompt := &survey.Select{
			Message: "Select client",
			Options: []string{
				`Codeforces`,
				`Szkopul`,
				`Sio2 (staszic.waw.pl)`,
			},
		}
		if err = survey.AskOne(prompt, &index); err != nil {
			return
		}
		if index == 0 {
			return codeforcesCln.ConfigLogin()
		} else if index == 1 {
			return szkopulCln.ConfigLogin()
		} else if index == 2 {
			return sioCln.ConfigLogin()
		}
	} else if index == 1 {
		return cfg.AddTemplate()
	} else if index == 2 {
		return cfg.RemoveTemplate()
	} else if index == 3 {
		return cfg.SetDefaultTemplate()
	} else if index == 4 {
		return cfg.SetGenAfterParse()
	} else if index == 5 {
		return cfg.SetCodeforcesHost()
	} else if index == 6 {
		return cfg.SetProxy()
	} else if index == 7 {
		return cfg.SetFolderName()
	} else if index == 8 {
		return cfg.SetDefaultNaming()
	} else if index == 9 {
		return cfg.SetDbPath()
	}
	return
}
