package sio_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Arapak/sio-tool/util"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"math"
	"strings"
	"time"
)

const ErrorNoRoundToRace = "there is no round to race in this contest"

type RoundInfo struct {
	Time           float64 `json:"time"`
	RoundStartDate float64 `json:"round_start_date"`
	RoundName      string  `json:"round_name"`
	Username       string  `json:"user"`
}

// RaceContest wait for contest starting
func (c *SioClient) RaceContest(info Info) (round string, err error) {
	color.Cyan("Race " + info.Hint())

	URL, err := info.ContestURL(c.host)
	if err != nil {
		return
	}

	body, err := util.GetBody(c.client, strings.TrimSuffix(URL, "/p")+"/status")
	if err != nil {
		return
	}
	roundInfo := RoundInfo{}
	err = json.Unmarshal(body, &roundInfo)
	if err != nil {
		return
	}
	if roundInfo.Username == "" {
		err = errors.New(ErrorNotLogged)
		return
	}
	if roundInfo.RoundName == "" {
		err = errors.New(ErrorNoRoundToRace)
		return
	}
	color.Cyan("Round %v", roundInfo.RoundName)
	timeLeft := int64(math.Round(roundInfo.RoundStartDate - roundInfo.Time))
	color.Green("Countdown: ")
	for timeLeft > 0 {
		h := timeLeft / 60 / 60
		m := timeLeft/60 - h*60
		s := timeLeft - h*60*60 - m*60
		fmt.Printf("%02d:%02d:%02d\n", h, m, s)
		ansi.CursorUp(1)
		timeLeft--
		time.Sleep(time.Second)
	}
	time.Sleep(900 * time.Millisecond)
	return roundInfo.RoundName, nil
}
