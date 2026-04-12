package main

import (
	"fmt"
	"log"

	"github.com/afifn11/gopay-x/api-gateway/config"
	"github.com/afifn11/gopay-x/api-gateway/internal/handler"
)

func main() {
	cfg := config.Load()

	r := handler.NewRouter(cfg)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 %s running on %s", cfg.App.Name, addr)
	log.Printf("📡 Proxying to %d services", 8)

	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}