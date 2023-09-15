package cmd

import (
	"errors"
	"regexp"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/szkopul_client"

	"github.com/fatih/color"
)

const ErrorProblemIDNotFound = "problem id not found"

func getProblemIDFromLink(link string) (problemID string, err error) {
	reg := regexp.MustCompile(SzkopulLinkRegStr)
	names := reg.SubexpNames()
	for i, val := range reg.FindStringSubmatch(link) {
		if names[i] == "problemSecretKey" && val != "" {
			return val, nil
		}
	}
	return "", errors.New(ErrorProblemIDNotFound)
}

func SzkopulSubmit() (err error) {
	cln := szkopul_client.Instance
	err = cln.Ping()
	if err != nil {
		return
	}
	cfg := config.Instance
	filename, _, err := getOneCode(Args.File, cfg.Template, szkopul_client.AcceptedExtensions)
	if err != nil {
		return
	}
	if Args.SzkopulInfo.ProblemID == "" {
		link, err := searchForLinkSzkopul()
		if err != nil {
			return err
		}
		Args.SzkopulInfo.ProblemID, err = getProblemIDFromLink(link)
		if err != nil {
			return err
		}
		color.Green("Found problem secret key")
	}
	info := Args.SzkopulInfo

	if err = cln.Submit(info, filename); err != nil {
		if err = loginAgainSzkopul(cln, err); err == nil {
			err = cln.Submit(info, filename)
		}
	}
	return
}

func loginAgainSzkopul(cln *szkopul_client.SzkopulClient, err error) error {
	if err != nil && err.Error() == szkopul_client.ErrorNotLogged {
		color.Red("Not logged. Try to login\n")
		err = cln.Login()
	}
	return err
}
