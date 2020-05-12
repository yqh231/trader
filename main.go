package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/yqh231/trader/logger"
	"github.com/yqh231/trader/server"
	toml "github.com/pelletier/go-toml"
)


func main() {
	initLogger()
	toml := initConfig()
	
	srv := server.NewServer(toml)

	srv.Consume(toml.Get("system.worker").(int))

	waitForSignal()
}

func initLogger() {
	_  = logger.GetLogger() // init
}


func initConfig() *toml.Tree {
	content, err := ioutil.ReadFile("default.toml")
	if err != nil {
		panic("Read default toml fail")
	}
	config, err := toml.Load(string(content))

	if err != nil {
		panic("Load default toml fail")
	}
	return config
}


func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
}



