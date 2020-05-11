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

const (
	DEPTH = iota + 1
)


const (
	CARRY = iota + 1
)

type CoinExClient struct {
	client *httpClient

	accessId string
	apiKey string
	wsRoot string
	message chan *MessageItems
}

type MessageItems struct {
	Type int
	Strategy int
	Content []byte
}


func NewClient(toml *toml.Tree) *CoinExClient {
	root := toml.Get("coinex.Root").(string)
	apiKey := toml.Get("coinex.ApiKey").(string)
	accessId := toml.Get("coinex.AccessId").(string)
	wsRoot := toml.Get("coinex.WsRoot").(string)
	c := newHttpClient(root, time.Second * 1)

	return &CoinExClient{
		client: c,

		accessId: accessId,
		apiKey: apiKey,
		wsRoot: wsRoot,
		message: make(chan *MessageItems, 10),
	}
}


func (c *CoinExClient) sign(params map[string]string) string{
	var (
		buf []string
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

	return strings.ToUpper(string(secret[:]))
}

func (c *CoinExClient) addTonce(params map[string]string) {
	params["tonce"] = fmt.Sprintf("%v", time.Now().UnixNano() / 1e6)
}

func (c *CoinExClient) BookDepth(market string) {
	u := url.URL{Scheme: "wss", Host: c.wsRoot, Path: ""}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {

	}

	req, _ := json.Marshal(map[string]interface{}{
		"method": "depth.query",
		"params": []string{market, "20", "0"},
		"id": 0,
	})

	conn.WriteMessage(websocket.TextMessage, req)

	go func() {
		defer conn.Close()
		for {
			_, message ,err := conn.ReadMessage()
			if err != nil {
				continue
			}

			c.message <- &MessageItems{
				Type: DEPTH,
				Content: message,
			}
		}
	}()
}

func (c *CoinExClient) routeFunc(m *MessageItems) {
	switch m.Type {
	case DEPTH:
		c.routeDepth(m)
	}
}

func (c *CoinExClient) routeDepth(m *MessageItems) {
	switch m.Strategy {
	case CARRY:
		
	}
}


func (c *CoinExClient) Consume(nums int) {

	for i := 0; i < nums; i ++{
		go func() {
			for {
				select {
				case msg := <- c.message:
					c.routeFunc(msg)
				}
			}
		}()
	}
}
