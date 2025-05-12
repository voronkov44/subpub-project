package main

import (
	"log"
	"subpub-project/configs"
	"subpub-project/internal/subpub"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}

	// Создаем новую шину событий
	pubSub := subpub.NewSubPub()

	// Запускаем gRPC сервер
	if err := subpub.StartGRPCServer(cfg.Server.Address, pubSub); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
