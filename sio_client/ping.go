package sio_client

import (
	"errors"

	"github.com/Arapak/sio-tool/util"
)

const ErrorSioIsUnavailable = "sio is unavailable (check your internet connection)"

func (c *SioClient) Ping() (err error) {
	_, err = util.GetBody(c.client, c.host)
	if err != nil {
		return errors.New(ErrorSioIsUnavailable)
	}
	return
}
