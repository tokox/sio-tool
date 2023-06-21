package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"path/filepath"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/sio_client"
)

func SioParse() (err error) {
	cfg := config.Instance
	cln := sio_client.Instance
	info := Args.SioInfo
	source := ""
	ext := ""
	if cfg.GenAfterParse {
		if len(cfg.Template) == 0 {
			return errors.New("you have to add at least one code template by `st config`")
		}
		path := cfg.Template[cfg.Default].Path
		ext = filepath.Ext(path)
		if source, err = readTemplateSource(path, cln.Username); err != nil {
			return
		}
	}

	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		fmt.Printf("failed to open database connection: %v\n", err)
		return
	}
	defer db.Close()

	work := func() error {
		_, paths, err := cln.Parse(info, db)
		if err != nil {
			return err
		}
		if cfg.GenAfterParse {
			for _, path := range paths {
				err = GenFiles(source, path, ext)
				if err != nil {
					color.Red(err.Error())
				}
			}
		}
		return nil
	}
	if err = work(); err != nil {
		if err = loginAgainSio(cln, err); err == nil {
			err = work()
		}
	}
	return
}
