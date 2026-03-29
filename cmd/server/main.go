package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"xray_server/internal/delivery/http"
	"xray_server/internal/domain"
	"xray_server/internal/repository"
	"xray_server/internal/usecase/auth"
	"xray_server/internal/usecase/config"
	"xray_server/internal/usecase/stats"
	pkgconfig "xray_server/pkg/config"
)

const (
	defaultConfigPath = "config.json"
	httpPort          = 8080
	dbPath            = "users.db"
)

func main() {
	if err := godotenv.Overload(); err != nil {
		fmt.Println("warning: Не удалось загрузить переменные окружения из .env")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cfg, err := pkgconfig.LoadEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	configUseCase := config.NewConfigUseCase(config.ConfigOptions{
		Mldsa65Seed:   cfg.Mldsa65Seed,
		Mldsa65Public: cfg.Mldsa65Public,
		PrivateKey:    cfg.PrivateKey,
		PublicKey:     cfg.PublicKey,
		ShortIdsRaw:   cfg.ShortIdsRaw,
		ConfigPath:    configPath,
		ServerIP:      cfg.ServerIP,
	})

	if err := configUseCase.SaveConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сохранения конфигурации: %v\n", err)
		os.Exit(1)
	}

	userRepo, err := repository.NewUserRepositorySQLite(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации БД пользователей: %v\n", err)
		os.Exit(1)
	}
	defer userRepo.Close()

	_, err = userRepo.FindByUsername(cfg.AdminUsername)
	if err != nil {
		adminUser := &domain.User{
			Username: cfg.AdminUsername,
			Password: cfg.AdminPassword,
			Role:     domain.RoleAdmin,
		}
		if err := userRepo.Create(adminUser); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка создания админа: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Админ %s создан\n", cfg.AdminUsername)
	}

	jwtService := repository.NewJWTService(cfg.JWTSecret, 24*time.Hour)
	loginUseCase := auth.NewLoginUseCase(userRepo, jwtService)
	authHandler := http.NewAuthHandler(loginUseCase)
	authMiddleware := http.NewAuthMiddleware(jwtService)

	xrayAPIEndpoint := os.Getenv("XRAY_API_ENDPOINT")
	if xrayAPIEndpoint == "" {
		xrayAPIEndpoint = "xray:54321"
	}

	statsRepo, err := repository.NewStatsRepositorygRPC(xrayAPIEndpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации сервиса статистики: %v\n", err)
	}
	defer func() {
		if statsRepo != nil {
			statsRepo.Close()
		}
	}()

	statsUseCase := stats.NewStatsUseCase(statsRepo)
	statsHandler := http.NewStatsHandler(statsUseCase)
	clientHandler := http.NewClientHandler(configUseCase)

	server := http.NewServer(authHandler, clientHandler, statsHandler, authMiddleware, httpPort)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка HTTP сервера: %v\n", err)
		os.Exit(1)
	}
}
