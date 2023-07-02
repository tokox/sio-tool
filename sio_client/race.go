package sio_client

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"math"
	"time"
)

func (c *SioClient) RaceContest(info Info) (round string, err error) {
	color.Cyan("Race " + info.Hint())
	roundInfo, err := c.status(info)
	color.Cyan("Round %v", roundInfo.RoundName)
	timeLeft := int64(math.Round(roundInfo.RoundStartDate - roundInfo.Time))
	if timeLeft <= 0 {
		return roundInfo.RoundName, nil
	}
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
