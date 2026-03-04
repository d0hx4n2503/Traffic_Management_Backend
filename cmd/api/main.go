package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/server"
	"github.com/adohong4/driving-license/pkg/db/postgres"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
)

// @title Traffic License REST API
// @version 1.0
// @description REST API for Traffic License Management
// @contact.url https://github.com/adohong4
// @BasePath /v1/api

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Type 'Bearer {your-jwt-token}' to authenticate. This is required for protected endpoints.
func main() {
	log.Println("Starting driving license API server")
	loadDotEnvIfExists(".env")

	configEnv := os.Getenv("config")
	if configEnv == "" {
		configEnv = os.Getenv("CONFIG")
	}
	configPath := utils.GetConfigPath(configEnv)

	// Read & Analyst Config
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	if renderPort := os.Getenv("PORT"); renderPort != "" {
		if !strings.HasPrefix(renderPort, ":") {
			renderPort = ":" + renderPort
		}
		cfg.Server.Port = renderPort
	}

	// Initialize Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Connect PostgreSQL
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %v", err)
	}
	defer psqlDB.Close()
	appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())

	// Run Server
	s := server.NewServer(cfg, psqlDB, appLogger)
	if err := s.Run(); err != nil {
		log.Fatalf("Server run: %v", err)
	}
}

func loadDotEnvIfExists(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if len(val) >= 2 {
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) || (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				val = val[1 : len(val)-1]
			}
		}
		if key == "" {
			continue
		}
		_ = os.Setenv(key, val)
	}
}
