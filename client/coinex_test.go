package client

import (
	"fmt"
	"io/ioutil"
	"testing"

	toml "github.com/pelletier/go-toml"
)

func initConfig() *toml.Tree {
	content, err := ioutil.ReadFile("../default.toml")
	if err != nil {
		panic("Read default toml fail")
	}
	config, err := toml.Load(string(content))

	if err != nil {
		panic("Load default toml fail")
	}
	return config
}

func TestWs(t *testing.T) {

	toml := initConfig()
	coinex := NewClient(toml)
	coinex.BookDepth("BTCUSDT")

	go func() {
		for {
			select {
			case msg := <-coinex.Message:
				fmt.Println(string(msg.Content))
			}

		}
	}()
    select {}
}
