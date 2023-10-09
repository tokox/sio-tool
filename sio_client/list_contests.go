package sio_client

import (
	"fmt"
	"strings"

	"github.com/Arapak/sio-tool/util"
	"github.com/PuerkitoBio/goquery"
)

type ContestInfo struct {
	Name      string
	Alias     string
	Subheader bool
}

func (c *SioClient) findContests(body []byte) (ret []ContestInfo, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	doc.Find("table tbody").First().Find("tr").Each(func(_ int, s *goquery.Selection) {
		_, ok := s.Attr("class")
		info := ContestInfo{}
		if !ok && c.instanceClient == Staszic {
			info.Subheader = true
			info.Name = strings.TrimSpace(s.Find("a").First().Text())
		} else {
			info.Subheader = false
			info.Name = strings.TrimSpace(s.Find("a").First().Text())
			info.Alias = strings.TrimSpace(s.Find("td").First().Text())
		}
		ret = append(ret, info)
	})
	return
}

func (c *SioClient) ListContests() (problems []ContestInfo, perf util.Performance, err error) {
	perf.StartFetching()
	URL := c.host
	if c.instanceClient == Mimuw {
		URL = fmt.Sprintf("%v/c/oi30-1/contest", c.host)
	}
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}
	perf.StopFetching()

	perf.StartParsing()

	if _, err = findUsername(body); err != nil {
		return
	}

	contests, err := c.findContests(body)
	if err != nil {
		return
	}

	perf.StopParsing()

	return contests, perf, nil
}
