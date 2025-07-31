package download

import (
	"time"

	"github.com/beeploop/foorrent/client"
)

func (dm *DownloadManager) peerResponder(c *client.Client) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			c.SendKeepAlive()
		}
	}
}
