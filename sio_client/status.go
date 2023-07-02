package sio_client

import (
	"encoding/json"
	"errors"
	"github.com/Arapak/sio-tool/util"
)

const ErrorNoActiveRound = "there is no active round in this contest"

type RoundInfo struct {
	Time           float64 `json:"time"`
	RoundStartDate float64 `json:"round_start_date"`
	RoundName      string  `json:"round_name"`
	Username       string  `json:"user"`
}

func (c *SioClient) status(info Info) (roundInfo RoundInfo, err error) {
	URL, err := info.StatusURL(c.host)
	if err != nil {
		return
	}
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &roundInfo)
	if err != nil {
		return
	}
	if roundInfo.Username == "" {
		err = errors.New(ErrorNotLogged)
		return
	}
	if roundInfo.RoundName == "" {
		err = errors.New(ErrorNoActiveRound)
		return
	}
	return
}
