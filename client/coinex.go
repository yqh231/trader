package client

import (
	"time"

	toml "github.com/pelletier/go-toml"
)

type CoinExClient struct {
	client *httpClient

	accessId string
	apiKey string
}


func NewClient(toml *toml.Tree) *CoinExClient {
	root := toml.Get("coinex.root").(string)
	apiKey := toml.Get("coinex.ApiKey").(string)
	accessId := toml.Get("coinex.AccessId").(string)
	c := newHttpClient(root, time.Second * 1)

	return &CoinExClient{
		client: c,

		accessId: accessId,
		apiKey: apiKey,
	}
}

