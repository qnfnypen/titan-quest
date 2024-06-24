package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gnasnik/titan-quest/api"
	"github.com/gnasnik/titan-quest/config"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/spf13/viper"
)

func main() {
	OsSignal := make(chan os.Signal, 1)

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading config file: %v\n", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshaling config file: %v\n", err)
	}

	config.Cfg = cfg

	if err := dao.Init(&cfg); err != nil {
		log.Fatalf("initital: %v\n", err)
	}

	api.InitCaptcha()

	go api.ServerAPI(&cfg)

	api.InitBot()

	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM)
	_ = <-OsSignal

	fmt.Printf("Exiting received OsSignal\n")
}
