package server

import (
	"github.com/yqh231/trader/client"
	toml "github.com/pelletier/go-toml"
)


type TraderServer struct {
	coinexManager *client.CoinExClient
}


func NewServer(toml *toml.Tree) *TraderServer {

	return &TraderServer{
		coinexManager: client.NewClient(toml),
	}
}

func (t *TraderServer) Consume(nums int) {
	for i := 0; i < nums; i++ {
		go func() {
			for {
				select {
				case coinexMsg := <- t.coinexManager.Message:
					t.consumeCoinex(coinexMsg)
				}
			}

		}()
	}
}

func (t *TraderServer) consumeCoinex(m *client.MessageItems) {
	switch m.Type {
	case client.CoinexDepth:
		
	}
}


