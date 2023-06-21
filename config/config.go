package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/Arapak/sio-tool/szkopul_client"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

type CodeTemplate struct {
	Alias        string   `json:"alias"`
	Lang         string   `json:"lang"`
	Path         string   `json:"path"`
	Suffix       []string `json:"suffix"`
	BeforeScript string   `json:"before_script"`
	Script       string   `json:"script"`
	AfterScript  string   `json:"after_script"`
}

type Config struct {
	Template       []CodeTemplate    `json:"template"`
	Default        int               `json:"default"`
	GenAfterParse  bool              `json:"gen_after_parse"`
	CodeforcesHost string            `json:"codeforces_host"`
	SzkopulHost    string            `json:"szkopul_host"`
	SioHost        string            `json:"sio_host"`
	Proxy          string            `json:"proxy"`
	FolderName     map[string]string `json:"folder_name"`
	DefaultNaming  map[string]string `json:"default_naming"`
	DbPath         string            `json:"db_path"`
	path           string
}

var Instance *Config

func Init(path string) {
	c := &Config{path: path, CodeforcesHost: "https://codeforces.com", SzkopulHost: "https://szkopul.edu.pl", SioHost: "https://sio2.staszic.waw.pl", DbPath: "~/.st/tasks.db", Proxy: ""}
	if err := c.load(); err != nil {
		color.Red(err.Error())
		color.Green("Create a new configuration in %v", path)
	}
	if c.Default < 0 || c.Default >= len(c.Template) {
		c.Default = 0
	}
	if c.FolderName == nil {
		c.FolderName = map[string]string{}
	}
	if _, ok := c.FolderName["sio-root"]; !ok {
		c.FolderName["sio-root"] = "~/st/sio"
	}
	if _, ok := c.FolderName["codeforces-root"]; !ok {
		c.FolderName["codeforces-root"] = "~/st/codeforces"
	}
	for _, problemType := range codeforces_client.ProblemTypes {
		if _, ok := c.FolderName[fmt.Sprintf("codeforces-%v", problemType)]; !ok {
			c.FolderName[fmt.Sprintf("codeforces-%v", problemType)] = problemType
		}
	}
	if _, ok := c.FolderName["szkopul-root"]; !ok {
		c.FolderName["szkopul-root"] = "~/st/szkopul"
	}
	for _, archive := range szkopul_client.Archives {
		if _, ok := c.FolderName[fmt.Sprintf("szkopul-%v", archive)]; !ok {
			c.FolderName[fmt.Sprintf("szkopul-%v", archive)] = archive
		}
	}

	if c.DefaultNaming == nil {
		c.DefaultNaming = map[string]string{}
	}
	if _, ok := c.DefaultNaming["solve"]; !ok {
		c.DefaultNaming["solve"] = "$%task%$.cpp"
	}
	if _, ok := c.DefaultNaming["brute"]; !ok {
		c.DefaultNaming["brute"] = "$%task%$-brute.cpp"
	}
	if _, ok := c.DefaultNaming["gen"]; !ok {
		c.DefaultNaming["gen"] = "$%task%$-gen.cpp"
	}
	if _, ok := c.DefaultNaming["test_in"]; !ok {
		c.DefaultNaming["test_in"] = "$%task%$GenTest$%test%$.in"
	}
	err := c.save()
	if err != nil {
		color.Red(err.Error())
		return
	}
	c.FolderName["sio-root"], err = homedir.Expand(c.FolderName["sio-root"])
	if err != nil {
		color.Red(err.Error())
	}
	c.FolderName["codeforces-root"], err = homedir.Expand(c.FolderName["codeforces-root"])
	if err != nil {
		color.Red(err.Error())
	}
	c.FolderName["szkopul-root"], err = homedir.Expand(c.FolderName["szkopul-root"])
	if err != nil {
		color.Red(err.Error())
	}
	c.DbPath, err = homedir.Expand(c.DbPath)
	if err != nil {
		color.Red(err.Error())
	}
	Instance = c
}

func (c *Config) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) save() (err error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(c)
	if err == nil {
		err = os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		if err == nil {
			err = os.WriteFile(c.path, data.Bytes(), 0644)
		}
	}
	if err != nil {
		color.Red("Cannot save config to %v\n%v", c.path, err.Error())
	}
	return
}
