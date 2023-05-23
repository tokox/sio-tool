package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
)

func parseTemplate(source string, cln *codeforces_client.CodeforcesClient) string {
	now := time.Now()
	source = strings.ReplaceAll(source, "$%U%$", cln.Handle)
	source = strings.ReplaceAll(source, "$%Y%$", fmt.Sprintf("%v", now.Year()))
	source = strings.ReplaceAll(source, "$%M%$", fmt.Sprintf("%02v", int(now.Month())))
	source = strings.ReplaceAll(source, "$%D%$", fmt.Sprintf("%02v", now.Day()))
	source = strings.ReplaceAll(source, "$%h%$", fmt.Sprintf("%02v", now.Hour()))
	source = strings.ReplaceAll(source, "$%m%$", fmt.Sprintf("%02v", now.Minute()))
	source = strings.ReplaceAll(source, "$%s%$", fmt.Sprintf("%02v", now.Second()))
	return source
}

func readTemplateSource(path string, cln *codeforces_client.CodeforcesClient) (source string, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	source = parseTemplate(string(b), cln)
	return
}

func GenFiles(source, currentPath, ext string) error {
	path := filepath.Join(currentPath, filepath.Base(currentPath))

	savePath := path + ext
	i := 1
	for _, err := os.Stat(savePath); err == nil; _, err = os.Stat(savePath) {
		tmpPath := fmt.Sprintf("%v%v%v", path, i, ext)
		fmt.Printf("%v exists. Rename to %v\n", filepath.Base(savePath), filepath.Base(tmpPath))
		savePath = tmpPath
		i++
	}

	err := os.WriteFile(savePath, []byte(source), 0644)
	if err == nil {
		color.Green("Generated! See %v", filepath.Base(savePath))
	}
	return err
}

// Gen command
func Gen() (err error) {
	cfg := config.Instance
	if len(cfg.Template) == 0 {
		return errors.New("you have to add at least one code template by `st config`")
	}
	alias := Args.Alias
	var path string

	if alias != "" {
		templates := cfg.TemplateByAlias(alias)
		if len(templates) < 1 {
			return fmt.Errorf("cannot find any template with alias %v", alias)
		} else if len(templates) == 1 {
			path = templates[0].Path
		} else {
			fmt.Printf("There are multiple templates with alias %v\n", alias)
			for i, template := range templates {
				fmt.Printf(`%3v: "%v"`, i, template.Path)
				fmt.Println()
			}
			i := util.ChooseIndex(len(templates))
			path = templates[i].Path
		}
	} else {
		path = cfg.Template[cfg.Default].Path
	}

	cln := codeforces_client.Instance
	source, err := readTemplateSource(path, cln)
	if err != nil {
		return
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return
	}

	ext := filepath.Ext(path)
	return GenFiles(source, currentPath, ext)
}
