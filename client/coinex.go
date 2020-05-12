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
	"github.com/yqh231/trader/logger"
)

type CoinExClient struct {
	client *httpClient

	accessId string
	apiKey   string
	wsRoot   string
	Message  chan *MessageItems
}

type MarketOrderDepth struct {
	Code    int       `json:"code"`
	Data    OrderData `json:"data"`
	Message string    `json:"message"`
}

type OrderData struct {
	Amount       string `json:"amount"`
	AvgPrice     string `json:"avg_price"`
	CreateTime   int    `json:"create_time"`
	DealAmount   string `json:"deal_amount"`
	DealMoney    string `json:"deal_money"`
	Id           int    `json:"id"`
	Left         string `json:"left"`
	MakerFeeRate string `json:"string"`
	Market       string `json:"market"`
	OrderType    string `json:"order_type"`
	Price        string `json:"price"`
	SourceId     string `json:"source_id"`
	Status       string `json:"status"`
	TakerFeeRate string `json:"taker_fee_rate"`
	Type         string `json:"type"`
	ClientId     string `json:"client_id"`
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

func (c *CoinExClient) sign(params map[string]interface{}) {
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
		secretStr += fmt.Sprintf("&%s=%v", key, params[key])
	}
	secretStr += fmt.Sprintf("&secret_key=%s", c.apiKey)

	secret := md5.Sum([]byte(secretStr))

	c.client.setHeaders(map[string]string{
		"AUTHORIZATION": strings.ToUpper(string(secret[:])),
		"Content-Type":  "application/json",
	})
}

func (c *CoinExClient) addTonce(params map[string]interface{}) {
	params["tonce"] = fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
}

func (c *CoinExClient) BookDepth(market string) error {
	u := url.URL{Scheme: "wss", Host: c.wsRoot, Path: ""}
	l := logger.GetLogger()
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		l.Zap.Warnw(err.Error())
		return err
	}

	req, _ := json.Marshal(map[string]interface{}{
		"method": "depth.subscribe",
		"params": []interface{}{market, 5, "0"},
		"id":     0,
	})

	conn.WriteMessage(websocket.TextMessage, req)

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				l.Zap.Warnw(err.Error())
				continue
			}
			c.Message <- &MessageItems{
				Type:    CoinexDepth,
				Content: message,
			}
		}
	}()

	return nil
}

func (c *CoinExClient) PutMarketOrder(params map[string]interface{}) (*MarketOrderDepth, error) {
	var (
		err  error
		resp *httpResp
		url  = "/v1/order/market"
		l    = logger.GetLogger()
	)
	c.sign(params)
	resp, err = c.client.Post(url, params)
	if err != nil {
		l.Zap.Warnw(err.Error())
		return nil, nil
	}

	var m MarketOrderDepth
	resp.unmarshal(&m)

	return &m, nil
}
