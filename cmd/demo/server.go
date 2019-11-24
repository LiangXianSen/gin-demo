package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LiangXianSen/gin-demo/config"
	"github.com/LiangXianSen/gin-demo/service/web"
)

func main() {
	// Load config from toml file
	conf, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	webSrv := web.NewServer(conf)

	go webSrv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	webSrv.Shutdown()
}
