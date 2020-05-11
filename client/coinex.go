package client

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	toml "github.com/pelletier/go-toml"
)

type CoinExClient struct {
	client *httpClient

	accessId string
	apiKey   string
	wsRoot   string
	Message  chan *MessageItems
}

type MarketOrderDepth struct {
	Code int `json: "code"`
	Data OrderData `json: "data"`
	Message string `json: "message"`
}

type OrderData struct {
	Amount string `json: "amount"`
	AvgPrice string `json: "avg_price"`

}

func NewClient(toml *toml.Tree) *CoinExClient {
	root := toml.Get("coinex.Root").(string)
	apiKey := toml.Get("coinex.ApiKey").(string)
	accessId := toml.Get("coinex.AccessId").(string)
	wsRoot := toml.Get("coinex.WsRoot").(string)
	c := newHttpClient(root, time.Second*1)

	return &CoinExClient{
		client: c,

		accessId: accessId,
		apiKey:   apiKey,
		wsRoot:   wsRoot,
		Message:  make(chan *MessageItems, 10),
	}
}

func (c *CoinExClient) sign(params map[string]string) {
	var (
		buf       []string
		secretStr string
	)

	c.addTonce(params)
	for param := range params {
		buf = append(buf, param)
	}

	sort.Strings(buf)

	secretStr += fmt.Sprintf("access_id=%s", c.accessId)
	for _, key := range buf {
		secretStr += fmt.Sprintf("&%s=%s", key, params[key])
	}
	secretStr += fmt.Sprintf("&secret_key=%s", c.apiKey)

	secret := md5.Sum([]byte(secretStr))

	c.client.setHeaders(map[string]string{
		"AUTHORIZATION": strings.ToUpper(string(secret[:])),
		"Content-Type": "application/json",
	})
}

func (c *CoinExClient) addTonce(params map[string]string) {
	params["tonce"] = fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
}

func (c *CoinExClient) BookDepth(market string) {
	u := url.URL{Scheme: "wss", Host: c.wsRoot, Path: ""}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {

	}

	req, _ := json.Marshal(map[string]interface{}{
		"method": "depth.query",
		"params": []string{market, "20", "0"},
		"id":     0,
	})

	conn.WriteMessage(websocket.TextMessage, req)

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				continue
			}

			c.Message <- &MessageItems{
				Type:    CoinexDepth,
				Content: message,
			}
		}
	}()
}

func (c *CoinExClient) PutMarketOrder(params map[string]string) {
	var (
		err error
		resp *httpResp
		url = "/v1/order/market"
	)
	c.sign(params)	
	resp, err = c.client.Post(url, params)
	if err != nil {

	}

	// resp.unmarshal()
}