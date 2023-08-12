package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Arapak/sio-tool/util"
	"github.com/mitchellh/go-homedir"

	"github.com/Arapak/sio-tool/codeforces_client"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

func validateTemplatePath(path interface{}) (err error) {
	path, err = homedir.Expand(path.(string))
	if err != nil {
		return
	}
	if !filepath.IsAbs(path.(string)) {
		return fmt.Errorf("this is not an absolute path: %v", path)
	}
	stats, err := os.Stat(path.(string))
	if err != nil {
		return
	}
	if stats.IsDir() {
		return errors.New("this is a directory")
	}
	return
}

func (c *Config) AddTemplate() (err error) {
	color.Cyan("Add a template")
	type kv struct {
		K, V string
	}
	var langs []kv
	for k, v := range codeforces_client.Langs {
		langs = append(langs, kv{k, v})
	}
	sort.Slice(langs, func(i, j int) bool { return langs[i].V < langs[j].V })
	langValues := make([]string, len(langs))
	for i, t := range langs {
		langValues[i] = t.V
	}
	langID := 0
	if err = survey.AskOne(&survey.Select{Message: `Select a language`, Options: langValues}, &langID); err != nil {
		return
	}

	note := `Template:
  You can insert some placeholders into your template code. When generate a code from the
  template, st will replace all placeholders by following rules:

  $%U%$   Handle (e.g. Arapak)
  $%Y%$   Year   (e.g. 2019)
  $%M%$   Month  (e.g. 04)
  $%D%$   Day    (e.g. 09)
  $%h%$   Hour   (e.g. 08)
  $%m%$   Minute (e.g. 05)
  $%s%$   Second (e.g. 00)`
	_, _ = ansi.Println(note)
	color.Cyan(`Template absolute path(e.g. "~/template/io.cpp"): `)
	path := ""
	if err = survey.AskOne(&survey.Input{Message: `Template absolute path(e.g. "~/template/io.cpp"):`}, &path, survey.WithValidator(validateTemplatePath)); err != nil {
		return
	}
	path, err = homedir.Expand(path)
	if err != nil {
		return
	}

	color.Cyan(`The suffix of template above will be added by default.`)
	suffixes := ""
	util.GetValue(`Other suffix? (e.g. "cxx cc"), empty is ok:`, &suffixes, false)
	tmpSuffix := strings.Fields(suffixes)
	tmpSuffix = append(tmpSuffix, strings.Replace(filepath.Ext(path), ".", "", 1))
	suffixMap := map[string]bool{}
	var suffix []string
	for _, s := range tmpSuffix {
		if _, ok := suffixMap[s]; !ok {
			suffixMap[s] = true
			suffix = append(suffix, s)
		}
	}

	alias := ""
	util.GetValue(`Template's alias (e.g. "cpp" "py"):`, &alias, true)

	color.Green("Script in template:")
	note = `Template will run 3 scripts in sequence when you run "st test":
    - before_script   (execute once)
    - script          (execute the number of samples times)
    - after_script    (execute once)
  You could set "before_script" or "after_script" to an empty string, meaning not executing.
  You have to run your program in "script" with standard input/output (no need to redirect).

  You can insert some placeholders in your scripts. When executing a script,
  st will replace all placeholders by the following rules:

  $%path%$   Path to source file (Excluding $%full%$, e.g. "/home/arapak/")
  $%full%$   Full name of source file (e.g. "a.cpp")
  $%file%$   Name of source file (Excluding suffix, e.g. "a")
  $%rand%$   Random string with 8 characters (including "a-z" "0-9")`
	_, _ = ansi.Println(note)

	beforeScript := ""
	util.GetValue(`Before script (e.g. "g++ $%full%$ -o $%file%$.e -std=c++17"), empty is ok:`, &beforeScript, false)

	script := ""
	util.GetValue(`Script (e.g. "./$%file%$.e" "python3 $%full%$"):`, &script, true)

	afterScript := ""
	util.GetValue(`After script (e.g. "rm $%file%$.e"), empty is ok:`, &afterScript, false)

	c.Template = append(c.Template, CodeTemplate{
		alias, langs[langID].K, path, suffix,
		beforeScript, script, afterScript,
	})
	makeItDefault := true
	prompt := &survey.Confirm{Message: `Make it default?`, Default: true}
	if err = survey.AskOne(prompt, &makeItDefault); err != nil {
		return
	}
	if makeItDefault {
		c.Default = len(c.Template) - 1
	}
	return c.save()
}

func (c *Config) RemoveTemplate() (err error) {
	if len(c.Template) == 0 {
		color.Red("There is no template. Please add one")
		return nil
	}

	templates := make([]string, len(c.Template))
	for i, template := range c.Template {
		star := " "
		if i == c.Default {
			star = color.New(color.FgGreen).Sprint("*")
		}
		templates[i] = fmt.Sprintf(`%v "%v" "%v"`, star, template.Alias, template.Path)
	}
	prompt := &survey.Select{
		Message: "Remove a template",
		Options: templates,
	}
	idx := 0
	if err = survey.AskOne(prompt, &idx); err != nil {
		return
	}
	c.Template = append(c.Template[:idx], c.Template[idx+1:]...)
	if idx == c.Default {
		c.Default = 0
	} else if idx < c.Default {
		c.Default--
	}
	return c.save()
}

func (c *Config) SetDefaultTemplate() (err error) {
	if len(c.Template) == 0 {
		color.Red("There is no template. Please add one")
		return nil
	}

	templates := make([]string, len(c.Template))
	for i, template := range c.Template {
		star := " "
		if i == c.Default {
			star = color.New(color.FgGreen).Sprint("*")
		}
		templates[i] = fmt.Sprintf(`%v "%v" "%v"`, star, template.Alias, template.Path)
	}
	prompt := &survey.Select{
		Message: "Set default template",
		Options: templates,
	}
	if err = survey.AskOne(prompt, &c.Default); err != nil {
		return
	}
	return c.save()
}

func (c *Config) TemplateByAlias(alias string) []CodeTemplate {
	var ret []CodeTemplate
	for _, template := range c.Template {
		if template.Alias == alias {
			ret = append(ret, template)
		}
	}
	return ret
}
