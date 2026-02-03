package main

import (
	"log"
	"task-system-go/internal/config"
	"task-system-go/internal/db"
	"task-system-go/internal/http/router"
)

func main() {
	cfg := config.Load()
	if err := db.Init(cfg); err != nil {
		log.Fatal(err)
	}

	r := router.New(cfg)
	log.Printf("Gin server listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
