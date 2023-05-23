package config

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/szkopul_client"
	"github.com/Arapak/sio-tool/util"

	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/mitchellh/go-homedir"
)

// SetGenAfterParse set it yes or no
func (c *Config) SetGenAfterParse() (err error) {
	c.GenAfterParse = util.YesOrNo(`Run "st gen" after "st parse" (y/n)? `)
	return c.save()
}

func formatHost(host string) (string, error) {
	reg := regexp.MustCompile(`https?://[\w\-]+(\.[\w\-]+)+/?`)
	if !reg.MatchString(host) {
		return "", fmt.Errorf(`invalid host "%v"`, host)
	}
	for host[len(host)-1:] == "/" {
		host = host[:len(host)-1]
	}
	return host, nil
}

func formatProxy(proxy string) (string, error) {
	reg := regexp.MustCompile(`[\w\-]+?://[\w\-]+(\.[\w\-]+)*(:\d+)?`)
	if !reg.MatchString(proxy) {
		return "", fmt.Errorf(`invalid proxy "%v"`, proxy)
	}
	return proxy, nil
}

// SetHost set host for Codeforces
func (c *Config) SetHost() (err error) {
	host, err := formatHost(c.CodeforcesHost)
	if err != nil {
		host = "https://codeforces.com"
	}
	color.Green("Current host domain is %v", host)
	color.Cyan(`Set a new host domain (e.g. "https://codeforces.com"`)
	color.Cyan(`Note: Don't forget the "http://" or "https://"`)
	for {
		host, err = formatHost(util.ScanlineTrim())
		if err == nil {
			break
		}
		color.Red(err.Error())
	}
	c.CodeforcesHost = host
	color.Green("New host domain is %v", host)
	return c.save()
}

// SetProxy set proxy for codeforces client
func (c *Config) SetProxy() (err error) {
	proxy, err := formatProxy(c.Proxy)
	if err != nil {
		proxy = ""
	}
	if len(proxy) == 0 {
		color.Green("Current proxy is based on environment")
	} else {
		color.Green("Current proxy is %v", proxy)
	}
	color.Cyan(`Set a new proxy (e.g. "http://127.0.0.1:2333", "socks5://127.0.0.1:1080"`)
	color.Cyan(`Enter empty line if you want to use default proxy from environment`)
	color.Cyan(`Note: Proxy URL should match "protocol://host[:port]"`)
	for {
		proxy, err = formatProxy(util.ScanlineTrim())
		if err == nil {
			break
		}
		color.Red(err.Error())
	}
	c.Proxy = proxy
	if len(proxy) == 0 {
		color.Green("Current proxy is based on environment")
	} else {
		color.Green("Current proxy is %v", proxy)
	}
	return c.save()
}

// SetFolderName set folder name
func (c *Config) SetFolderName() (err error) {
	color.Cyan(`Set folders' name`)
	color.Cyan(`Enter empty line if you don't want to change the value`)
	color.Green(`Codeforces root path (absolute) (current: %v)`, c.FolderName["codeforces-root"])
	if value := util.ScanlineTrim(); value != "" {
		value, err = homedir.Expand(value)
		if err != nil {
			color.Red(err.Error())
			return
		}
		if filepath.IsAbs(value) {
			c.FolderName["codeforces-root"] = value
		} else {
			color.Red("this is not an absolute path (leaving current)")
		}
	}
	for _, problemType := range codeforces_client.ProblemTypes {
		color.Green(`%v path (current: %v)`, problemType, c.FolderName[fmt.Sprintf("codeforces-%v", problemType)])
		if value := util.ScanlineTrim(); value != "" {
			c.FolderName[fmt.Sprintf("codeforces-%v", problemType)] = value
		}
	}
	color.Green(`Szkopul root path (absolute) (current: %v)`, c.FolderName["szkopul-root"])
	if value := util.ScanlineTrim(); value != "" {
		value, err = homedir.Expand(value)
		if err != nil {
			color.Red(err.Error())
			return
		}
		if filepath.IsAbs(value) {
			c.FolderName["szkopul-root"] = value
		} else {
			color.Red("this is not an absolute path (leaving current)")
		}
	}
	for _, archive := range szkopul_client.Archives {
		color.Green(`%v path (current: %v)`, archive, c.FolderName[fmt.Sprintf("szkopul-%v", archive)])
		if value := util.ScanlineTrim(); value != "" {
			c.FolderName[fmt.Sprintf("szkopul-%v", archive)] = value
		}
	}
	return c.save()
}

func (c *Config) SetDefaultNaming() (err error) {
	color.Cyan(`Set default naming (for stress testing purposes)`)
	color.Cyan(`Enter empty line if you don't want to change the value`)
	fmt.Printf(`You can insert $%%task%%$ placeholder in your filenames, which you will provide when using st stress-test command.`)
	color.Green(`Solution file name (current: %v)`, c.DefaultNaming["solve"])
	if value := util.ScanlineTrim(); value != "" {
		c.DefaultNaming["solve"] = value
	}
	color.Green(`Brute forces solution filename (current: %v)`, c.DefaultNaming["brute"])
	if value := util.ScanlineTrim(); value != "" {
		c.DefaultNaming["brute"] = value
	}
	color.Green(`Tests generator filename (current: %v)`, c.DefaultNaming["gen"])
	if value := util.ScanlineTrim(); value != "" {
		c.DefaultNaming["gen"] = value
	}
	fmt.Printf(`Here you can also insert $%%test%%$ placeholder in your filename, which will indicate the test number.`)
	color.Green(`Generated test filename (current: %v)`, c.DefaultNaming["test_in"])
	if value := util.ScanlineTrim(); value != "" {
		c.DefaultNaming["test_in"] = value
	}
	return c.save()
}
