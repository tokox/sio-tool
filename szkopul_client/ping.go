package szkopul_client

import (
	"errors"

	"github.com/Arapak/sio-tool/util"
)

const ErrorSzkopulIsUnavailable = "szkopul is unavailable (check your internet connection)"

func (c *SzkopulClient) Ping() (err error) {
	URL := APIPingURL(c.host)
	_, err = util.GetBody(c.client, URL)
	if err != nil {
		return errors.New(ErrorSzkopulIsUnavailable)
	}
	return
}
