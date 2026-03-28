package main

import (
	"fmt"
	"os"

	"xray_server/internal/app"
	"xray_server/internal/env"
	"xray_server/internal/http"
)

const (
	configFileName = "config.json"
	httpPort       = 8080
)

func main() {
	cfg, err := env.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	configService := app.NewConfigService(app.ConfigServiceOptions{
		Mldsa65Sign:   cfg.Mldsa65Sign,
		Mldsa65Public: cfg.Mldsa65Public,
		PrivateKey:    cfg.PrivateKey,
		PublicKey:     cfg.PublicKey,
		ShortIdsRaw:   cfg.ShortIdsRaw,
		ConfigPath:    configFileName,
		ServerIP:      cfg.ServerIP,
	})

	if err := configService.SaveConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сохранения конфигурации: %v\n", err)
		os.Exit(1)
	}

	server := http.NewServer(configService, httpPort)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка HTTP сервера: %v\n", err)
		os.Exit(1)
	}
}
