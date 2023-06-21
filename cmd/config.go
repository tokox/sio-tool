package cmd

import (
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

func Config() (err error) {
	cfg := config.Instance
	codeforcesCln := codeforces_client.Instance
	szkopulCln := szkopul_client.Instance
	sioCln := sio_client.Instance
	color.Cyan("Configure the tool")
	_, _ = ansi.Println(`0) login`)
	_, _ = ansi.Println(`1) add a template`)
	_, _ = ansi.Println(`2) delete a template`)
	_, _ = ansi.Println(`3) set default template`)
	_, _ = ansi.Println(`4) run "st gen" after "st parse"`)
	_, _ = ansi.Println(`5) set host domain`)
	_, _ = ansi.Println(`6) set proxy`)
	_, _ = ansi.Println(`7) set folders' name`)
	_, _ = ansi.Println(`8) set default naming`)
	_, _ = ansi.Println(`9) set database path`)
	index := util.ChooseIndex(10)
	if index == 0 {
		color.Cyan("Select client")
		_, _ = ansi.Println(`0) Codeforces`)
		_, _ = ansi.Println(`1) Szkopul`)
		_, _ = ansi.Println(`2) Sio2 (staszic.waw.pl)`)
		index = util.ChooseIndex(3)
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
		return cfg.SetHost()
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
