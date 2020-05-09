package client

import (
	//"fmt"
	"fmt"
	"testing"
	"time"
)

type DepthResponse struct {
	Code int
	Message string
	Data *DepthData
}

type DepthData struct {
	Last string
	Time int64
	Asks [][]string
	Bids [][]string
}
func TestHttpGet(t *testing.T) {
	var (
		root = "https://api.coinex.com"
		url = "/v1/market/depth"
	)

	c := newHttpClient(root, 2 * time.Second)

	c.setHeaders(map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	})

	resp, _ := c.Get(url, map[string]string{
		"market": "bchbtc",
		"limit": "5",
		"merge": "0",
	})

	depth := new(DepthResponse)
	fmt.Println(depth.Data.Asks)

	resp.unmarshal(depth)

	fmt.Println(depth.Data.Asks)
}
