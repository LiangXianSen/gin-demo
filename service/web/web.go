package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/LiangXianSen/gin-demo/config"
)

type Server struct {
	Port   string
	Router *gin.Engine
	exit   chan struct{}
}

func NewServer(conf *config.Config) *Server {
	return &Server{
		Port:   conf.General.Port,
		Router: LoadRouter(routeGroups, conf),
		exit:   make(chan struct{}),
	}
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:    s.Port,
		Handler: s.Router,
	}

	// Run and listen server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Running: %s\n", err)
		}
	}()

	fmt.Println("Server Listen on port", s.Port)

	defer close(s.exit)
	<-s.exit
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: %s\n", err)
	}
	fmt.Println("Server exited.")
}

func (s *Server) Shutdown() {
	s.exit <- struct{}{}

	if _, open := <-s.exit; !open {
		return
	}
}
