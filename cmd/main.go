package main

import (
	"github.com/Futturi/Raspisanie/internal/handler"
	"github.com/Futturi/Raspisanie/internal/repository"
	"github.com/Futturi/Raspisanie/internal/service"
	"github.com/Futturi/Raspisanie/internal/ws"
	"github.com/Futturi/Raspisanie/pkg"
	"github.com/Futturi/Raspisanie/script"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

func main() {
	logg := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logg)
	err := InitConfig()
	if err != nil {
		slog.Error("error with config", slog.Any("error", err))
	}
	cfg := pkg.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Dbname:   viper.GetString("db.namedb"),
		Sslmode:  viper.GetString("db.sslmode"),
	}
	db, err := pkg.InitPostgres(cfg)
	if err != nil {
		slog.Error("error with initializing db", slog.Any("error", err))
	}
	script.InitScript()
	rep := repository.NewRepository(db)
	ser := service.NewSerivce(rep)
	hub := ws.NewHub()
	wsHan := ws.NewHandler(hub, ser)
	han := handler.NewHandler(ser)
	server := new(service.Server)
	go hub.Run()
	if err := server.InitServer("8080", han.InitRoutes(wsHan)); err != nil {
		slog.Error("error with running in 8080", slog.Any("error", err))
	}
}

func InitConfig() error {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	return viper.ReadInConfig()
}
