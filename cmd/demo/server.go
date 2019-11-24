package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/LiangXianSen/gin-demo/service/web"
)

var (
	port = flag.String("port", ":8080", "listen on port")
)

func main() {
	flag.Parse()

	webSrv := web.NewServer(*port)

	go webSrv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	webSrv.Shutdown()
}
