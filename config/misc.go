package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/AlecAivazis/survey/v2"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

func (c *Config) SetGenAfterParse() (err error) {
	prompt := &survey.Confirm{Message: `Run "st gen" after "st parse"?`, Default: false}
	if err = survey.AskOne(prompt, &c.GenAfterParse); err != nil {
		return
	}
	return c.save()
}

func validateHost(host interface{}) error {
	reg := regexp.MustCompile(`https?://[\w\-]+(\.[\w\-]+)+/?`)
	if !reg.MatchString(host.(string)) {
		return fmt.Errorf(`invalid host "%v"`, host)
	}
	return nil
}

func formatHost(host string) string {
	for host[len(host)-1:] == "/" {
		host = host[:len(host)-1]
	}
	return host
}

func validateProxy(proxy interface{}) error {
	reg := regexp.MustCompile(`^$|[\w\-]+?://[\w\-]+(\.[\w\-]+)*(:\d+)?`)
	if !reg.MatchString(proxy.(string)) {
		return fmt.Errorf(`invalid proxy "%v"`, proxy)
	}
	return nil
}

func (c *Config) SetCodeforcesHost() (err error) {
	var host string
	if validateHost(c.CodeforcesHost) == nil {
		host = formatHost(c.CodeforcesHost)
	} else {
		host = "https://codeforces.com"
	}
	color.Green("Current host domain is %v", host)
	color.Cyan(`Set a new host domain (e.g. "https://codeforces.com")`)
	color.Cyan(`Note: Don't forget the "http://" or "https://"`)
	if err = survey.AskOne(&survey.Input{Message: `host:`}, &host, survey.WithValidator(validateHost)); err != nil {
		return
	}
	c.CodeforcesHost = formatHost(host)
	color.Green("New host domain is %v", host)
	return c.save()
}

func (c *Config) SetProxy() (err error) {
	proxy := c.Proxy
	if validateProxy(c.Proxy) != nil {
		proxy = ""
	}
	if len(proxy) == 0 {
		color.Green("Current proxy is based on environment")
	} else {
		color.Green("Current proxy is %v", proxy)
	}
	color.Cyan(`Set a new proxy (e.g. "http://127.0.0.1:2333", "socks5://127.0.0.1:1080")`)
	color.Cyan(`Enter empty line if you want to use default proxy from environment`)
	color.Cyan(`Note: Proxy URL should match "protocol://host[:port]"`)
	if err = survey.AskOne(&survey.Input{Message: `proxy:`}, &proxy, survey.WithValidator(validateProxy)); err != nil {
		return
	}
	c.Proxy = proxy
	if len(proxy) == 0 {
		color.Green("Current proxy is based on environment")
	} else {
		color.Green("Current proxy is %v", proxy)
	}
	return c.save()
}

func validateAbsolutePath(path interface{}) (err error) {
	if path.(string) != "" {
		path, err = homedir.Expand(path.(string))
		if err != nil {
			return
		}
		if !filepath.IsAbs(path.(string)) {
			return fmt.Errorf("this is not an absolute path: %v", path)
		}
	}
	return nil
}

func inputDontOverwriteEmpty(message string, value string, validator survey.Validator) (newValue string, err error) {
	message = fmt.Sprintf("%v (current: %v)", message, value)
	if validator == nil {
		if err = survey.AskOne(&survey.Input{Message: message}, &newValue); err != nil {
			return
		}
	} else {
		if err = survey.AskOne(&survey.Input{Message: message}, &newValue, survey.WithValidator(validator)); err != nil {
			return
		}
	}
	if newValue == "" {
		newValue = value
	}
	return
}

func (c *Config) SetFolderName() (err error) {
	color.Cyan(`Set folders' name`)
	color.Cyan(`Enter empty line if you don't want to change the value`)
	if c.FolderName["codeforces-root"], err = inputDontOverwriteEmpty(`Codeforces root path (absolute)`, c.FolderName["codeforces-root"], validateAbsolutePath); err != nil {
		return
	}
	if c.FolderName["codeforces-root"], err = homedir.Expand(c.FolderName["codeforces-root"]); err != nil {
		return
	}
	for _, problemType := range codeforces_client.ProblemTypes {
		if c.FolderName[fmt.Sprintf("codeforces-%v", problemType)], err = inputDontOverwriteEmpty(fmt.Sprintf(`Codeforces %v path`, problemType), c.FolderName[fmt.Sprintf("codeforces-%v", problemType)], nil); err != nil {
			return
		}
	}
	if c.FolderName["szkopul-root"], err = inputDontOverwriteEmpty(`Szkopul root path (absolute)`, c.FolderName["szkopul-root"], validateAbsolutePath); err != nil {
		return
	}
	if c.FolderName["szkopul-root"], err = homedir.Expand(c.FolderName["szkopul-root"]); err != nil {
		return
	}
	for _, archive := range szkopul_client.Archives {
		if c.FolderName[fmt.Sprintf("szkopul-%v", archive)], err = inputDontOverwriteEmpty(fmt.Sprintf(`Szkopul %v archive path`, archive), c.FolderName[fmt.Sprintf("szkopul-%v", archive)], nil); err != nil {
			return
		}
	}
	if c.FolderName["sio-staszic-root"], err = inputDontOverwriteEmpty(`Sio staszic root path (absolute)`, c.FolderName["sio-staszic-root"], validateAbsolutePath); err != nil {
		return
	}
	if c.FolderName["sio-staszic-root"], err = homedir.Expand(c.FolderName["sio-staszic-root"]); err != nil {
		return
	}

	if c.FolderName["sio-mimuw-root"], err = inputDontOverwriteEmpty(`Sio mimuw root path (absolute)`, c.FolderName["sio-mimuw-root"], validateAbsolutePath); err != nil {
		return
	}
	if c.FolderName["sio-mimuw-root"], err = homedir.Expand(c.FolderName["sio-mimuw-root"]); err != nil {
		return
	}
	return c.save()
}

func (c *Config) SetDefaultNaming() (err error) {
	color.Cyan(`Set default naming (for stress testing purposes)`)
	color.Cyan(`Enter empty line if you don't want to change the value`)
	fmt.Printf(`You can insert $%%task%%$ placeholder in your filenames, which you will provide when using st stress-test command.`)
	if c.DefaultNaming["solve"], err = inputDontOverwriteEmpty(`Solution file name`, c.DefaultNaming["solve"], nil); err != nil {
		return
	}
	if c.DefaultNaming["brute"], err = inputDontOverwriteEmpty(`Bruteforce solution filename`, c.DefaultNaming["brute"], nil); err != nil {
		return
	}
	if c.DefaultNaming["gen"], err = inputDontOverwriteEmpty(`Tests generator filename`, c.DefaultNaming["gen"], nil); err != nil {
		return
	}
	fmt.Printf(`Here you can also insert $%%test%%$ placeholder in your filename, which will indicate the test number.`)
	if c.DefaultNaming["test_in"], err = inputDontOverwriteEmpty(`Generated test filename`, c.DefaultNaming["test_in"], nil); err != nil {
		return
	}
	return c.save()
}

func validateDbPath(path interface{}) (err error) {
	err = validateAbsolutePath(path)
	if filepath.Ext(path.(string)) != ".db" {
		err = errors.New("wrong file extension")
	}
	return
}

func (c *Config) SetDbPath() (err error) {
	dbPath, err := homedir.Expand(c.DbPath)
	if err != nil {
		dbPath = "~/.st/tasks.db"
	}
	color.Green("Current database path is %v", dbPath)
	color.Cyan(`Set a new db path (e.g. "~/.st/tasks.db")`)
	color.Cyan(`Note: Don't forget the ".db" extension`)
	if err = survey.AskOne(&survey.Input{Message: `db path:`}, &dbPath, survey.WithValidator(validateDbPath)); err != nil {
		return
	}
	if c.DbPath, err = homedir.Expand(dbPath); err != nil {
		return
	}
	color.Green("New database path is %v", dbPath)
	return c.save()
}
