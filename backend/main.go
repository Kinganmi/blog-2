package main

import (
	"backend/middleware"
	"backend/routers"
	"backend/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gin.SetMode(utils.ServerInfo.RunMode)
	router := routers.InitRouter()

	defer middleware.CloseLogFile()

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", utils.ServerInfo.Host, utils.ServerInfo.Port),
		Handler:        router,
		ReadTimeout:    utils.ServerInfo.ReadTimeout,
		WriteTimeout:   utils.ServerInfo.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if e := s.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", e)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
