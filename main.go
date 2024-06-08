package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mickamy/go_todo_app/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %s %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %v", url)
	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}
