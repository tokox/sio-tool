package config

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

var oiTemplate string = `
	\/* 
 * Author: $%U%$
 * Time: $%Y%$-$%M%$-$%D%$ $%h%$:$%m%$:$%s%$
**/

#include <bits/stdc++.h>
using namespace std;

typedef long long ll;

int main() {
  ios::sync_with_stdio(false);
  cin.tie(0);

  return 0;
}
`

var oiTemplatePath = "~/.st/template.cpp"
var oiTemplateCompilation = "g++ -std=c++20 -Wpedantic -O3 -static -o $%path%$$%file%$.e $%path%$$%full%$"
var oiTemplateRun = "./$%path%$$%file%$.e"

func (c *Config) AddOiTemplate() (err error) {
	oiTemplatePath, err = homedir.Expand(oiTemplatePath)
	if err != nil {
		return
	}
	err = os.WriteFile(oiTemplatePath, []byte(oiTemplate), 0644)
	if err != nil {
		return
	}
	c.Template = append(c.Template, CodeTemplate{
		"oi-cpp", "54", oiTemplatePath, []string{"cpp", "cxx", "cc"},
		oiTemplateCompilation, oiTemplateRun, "",
	})
	return c.save()
}
