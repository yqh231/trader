package server

import (
	"github.com/yqh231/trader/client"
)

func (t *TraderServer) consumeCoinex(m *client.MessageItems) {
	switch m.Type {
	case client.CoinexDepth:
		
	}
}