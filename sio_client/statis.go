package sio_client

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Arapak/sio-tool/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

type StatisInfo struct {
	ID     string
	Name   string
	Alias  string
	Round  string
	Points string
}

const ErrorContestNotFound = "contest not found"

func (prob *StatisInfo) ParsePoint() string {
	if prob.Points == "" {
		return ""
	}
	points, err := strconv.Atoi(prob.Points)
	if err != nil {
		return ""
	}
	if points == 100 {
		return color.New(color.FgGreen).Sprint(points)
	} else if points > 0 {
		return color.New(color.FgCyan).Sprint(points)
	} else {
		return color.New(color.FgRed).Sprint(points)
	}
}

func findProblems(body []byte, curIns SioInstanceClient) (ret []StatisInfo, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	var round = "none"
	doc.Find("table tbody").First().Find("tr").Each(func(_ int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		if strings.Contains(class, "problemlist-subheader") {
			space := regexp.MustCompile(`\s+`)
			round = space.ReplaceAllString(s.Find("strong").Text(), " ")
			return
		}
		info := StatisInfo{Round: round}
		info.Name = strings.TrimPrefix(strings.TrimSpace(s.Find("a").First().Text()), "Zadanie ")
		info.Alias = strings.TrimSpace(s.Find("td").First().Text())
		info.ID = strings.TrimPrefix(s.Find("div").First().AttrOr("id", ""), "limits_")
		if curIns == Mimuw || curIns == Staszic {
			info.Points = strings.TrimSpace(s.Find(".label").First().Text())
		}
		if curIns == Talent {
			info.Points = strings.TrimSpace(s.Find(".badge").First().Text())
		}

		ret = append(ret, info)
	})
	return
}

func (c *SioClient) getActiveRoundName(info Info) string {
	roundInfo, err := c.status(info)
	if err != nil {
		return ""
	}
	return roundInfo.RoundName
}

func (c *SioClient) Statis(info Info) (problems []StatisInfo, perf util.Performance, err error) {
	URL, err := info.ContestURL(c.host)
	if err != nil {
		return
	}

	var body []byte
	var problemsOnPage []StatisInfo

	activeRoundName := c.getActiveRoundName(info)

	pageNum := 1
	for {
		perf.StartFetching()
		body, err = util.GetBody(c.client, fmt.Sprintf("%v/?page=%v", URL, pageNum))
		if err != nil {
			return
		}
		if bytes.Contains(body, []byte("<p>404 &mdash; Page not found</p>")) {
			err = errors.New(ErrorContestNotFound)
			return
		}
		perf.StopFetching()

		pageNum++

		perf.StartParsing()

		if _, err = findUsername(body); err != nil {
			return
		}

		problemsOnPage, err = findProblems(body, c.instanceClient)
		if err != nil {
			return
		}
		if len(problemsOnPage) == 0 {
			break
		}
		for i := range problemsOnPage {
			if problemsOnPage[i].Round == "" {
				problemsOnPage[i].Round = activeRoundName
			}
		}
		problems = append(problems, problemsOnPage...)

		perf.StopParsing()
	}
	var filteredProblems []StatisInfo
	for _, problem := range problems {
		if info.Round != "" && problem.Round != info.Round && clearString(problem.Round) != info.Round {
			continue
		}
		if info.ProblemID != "" && problem.ID != info.ProblemID {
			continue
		}
		if info.ProblemAlias != "" && problem.Alias != info.ProblemAlias {
			continue
		}

		filteredProblems = append(filteredProblems, problem)
	}
	return filteredProblems, perf, nil
}
