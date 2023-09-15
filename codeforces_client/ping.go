package codeforces_client

import (
	"errors"

	"github.com/Arapak/sio-tool/util"
)

const ErrorCodeforcesIsUnavailable = "codeforces is unavailable (check your internet connection)"

func (c *CodeforcesClient) Ping() (err error) {
	_, err = util.GetBody(c.client, c.host)
	if err != nil {
		return errors.New(ErrorCodeforcesIsUnavailable)
	}
	return
}
