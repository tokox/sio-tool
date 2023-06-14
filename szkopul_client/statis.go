package szkopul_client

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Arapak/sio-tool/util"
	"github.com/PuerkitoBio/goquery"
	roman "github.com/StefanSchroeder/Golang-Roman"
)

type StatisInfo struct {
	ID      string
	Name    string
	Alias   string
	Stage   string
	Contest string
	Points  string
}

const idRegExp = `problemgroups-(?P<contestID>\d+)(-e?(?P<stageID>[1-3]|ioi-elem))?`

const ErrorIDNotFound = `problem id not found`

func getInfoFromId(id string) (ret StatisInfo) {
	reg := regexp.MustCompile(idRegExp)
	names := reg.SubexpNames()
	for i, val := range reg.FindStringSubmatch(id) {
		if val == "" {
			continue
		}
		if names[i] == "stageID" {
			ret.Stage = val
		} else if names[i] == "contestID" {
			ret.Contest = val
		}
	}
	return
}

func GetAliasAndName(name string) (alias string, problemName string) {
	reg := regexp.MustCompile(`^(Zadanie |Task )?(?P<problemName>[\s\S]+?)( \((?P<alias>\S+)\))?$`)
	names := reg.SubexpNames()
	for i, val := range reg.FindStringSubmatch(name) {
		if names[i] == "problemName" {
			problemName = val
		} else if names[i] == "alias" {
			alias = val
		}
	}
	return
}

func getIdFromHref(href string) (string, error) {
	reg := regexp.MustCompile(`/problemset/problem/(\S+?)/site/`)
	tmp := reg.FindStringSubmatch(href)
	if len(tmp) < 2 {
		return "", errors.New(ErrorIDNotFound)
	}
	return string(tmp[1]), nil
}

func findProblems(body []byte) (ret []StatisInfo, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	doc.Find("#problemgroups").First().Find("table").Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Parent().Attr("id")
		global_info := getInfoFromId(id)
		contest_id, err := strconv.Atoi(global_info.Contest)
		if err != nil {
			return
		}
		global_info.Contest = roman.Roman(contest_id)
		s.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			info := global_info
			linkElement := tr.Find("a").First()
			info.Alias, info.Name = GetAliasAndName(strings.TrimSpace(linkElement.Text()))
			href, hrefExists := linkElement.Attr("href")
			if hrefExists {
				info.ID, err = getIdFromHref(href)
				if err != nil {
					return
				}
			}
			info.Points = strings.TrimSpace(tr.Find(".badge").First().Text())
			ret = append(ret, info)
		})
	})
	return
}

func (c *SzkopulClient) Statis(info Info) (problems []StatisInfo, perf util.Performance, err error) {
	URL, err := info.ProblemSetURL(c.host)
	if err != nil {
		return
	}

	perf.StartFetching()

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	perf.StopFetching()
	perf.StartParsing()

	if _, err = findUsername(body); err != nil {
		return
	}

	perf.StopParsing()

	problems, err = findProblems(body)
	if err != nil {
		return
	}
	var filteredProblems []StatisInfo
	for _, problem := range problems {
		if info.ProblemID != "" && problem.ID != info.ProblemID {
			continue
		}
		if info.ProblemAlias != "" && problem.Alias != info.ProblemAlias {
			continue
		}
		if info.ContestID != "" && problem.Contest != info.ContestID {
			continue
		}
		if info.StageID != "" && problem.Stage != info.StageID {
			continue
		}

		filteredProblems = append(filteredProblems, problem)
	}
	return filteredProblems, perf, nil
}
