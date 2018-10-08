package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lonord/sse"
)

func main() {
	opt := handleCmdArgs()
	addr := fmt.Sprintf("%s:%d", opt.Host, opt.Port)
	s := createWebServer()
	go func() {
		log.Printf("server listens at http://%s\n", addr)
		if err := s.Start(addr); err != nil {
			log.Println("shutting down the server")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	for {
		isExit := false
		select {
		case sig := <-signalChan:
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				log.Printf("got singal %s, exit", sig.String())
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := s.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
				isExit = true
				break
			}
		}
		if isExit {
			break
		}
	}
}

func createWebServer() *echo.Echo {
	ec := echo.New()
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())
	ec.Use(middleware.CORS())
	sb := NewSSESystemBoradcast()
	ec.Any("/system", func(c echo.Context) error {
		sb.handleClient(sse.GenerateClientID(), c.Response().Writer)
		return nil
	})
	return ec
}
