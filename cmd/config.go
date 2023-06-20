package cmd

import (
	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
)

// Config command
func Config() (err error) {
	cfg := config.Instance
	codeforces_cln := codeforces_client.Instance
	szkopul_cln := szkopul_client.Instance
	sio_cln := sio_client.Instance
	color.Cyan("Configure the tool")
	ansi.Println(`0) login`)
	ansi.Println(`1) add a template`)
	ansi.Println(`2) delete a template`)
	ansi.Println(`3) set default template`)
	ansi.Println(`4) run "st gen" after "st parse"`)
	ansi.Println(`5) set host domain`)
	ansi.Println(`6) set proxy`)
	ansi.Println(`7) set folders' name`)
	ansi.Println(`8) set default naming`)
	ansi.Println(`9) set database path`)
	index := util.ChooseIndex(10)
	if index == 0 {
		color.Cyan("Select client")
		ansi.Println(`0) Codeforces`)
		ansi.Println(`1) Szkopul`)
		ansi.Println(`2) Sio2 (staszic.waw.pl)`)
		index = util.ChooseIndex(3)
		if index == 0 {
			return codeforces_cln.ConfigLogin()
		} else if index == 1 {
			return szkopul_cln.ConfigLogin()
		} else if index == 2 {
			return sio_cln.ConfigLogin()
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
