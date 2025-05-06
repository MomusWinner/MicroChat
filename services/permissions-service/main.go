package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/MomusWinner/MicroChat/internal/proxyproto"
	"github.com/MomusWinner/MicroChat/services/permissions-service/internal/config"
	"github.com/MomusWinner/MicroChat/services/permissions-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	log.Print("Start App")
	conf, err := config.Load()

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp4", ":10000")
	if err != nil {
		log.Fatalln(err)
	}

	errChan := make(chan error)

	srv := grpc.NewServer()

	svc, err := service.New(conf)

	if err != nil {
		log.Fatalln(err)
	}

	log.Print(svc)

	proxyproto.RegisterCentrifugoProxyServer(srv, svc)

	exitCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

		cancel()

		srv.GracefulStop()

		close(errChan)

		if err := listener.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		errChan <- srv.Serve(listener)
	}()

	select {
	case err := <-errChan:
		log.Fatalln(err)
	case <-exitCtx.Done():
		log.Println("exit")
	}
}
