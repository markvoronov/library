package main

import (
	"fmt"
	"library"
	"library/cmd/common"
	"library/internal/config"
	handler "library/internal/http-server/handlers"
	"library/internal/lib/logger/sl"
	"library/internal/service"
	"library/internal/storage"
	"library/internal/storage/postgres"
	"log/slog"
	"os"
)

func main() {

	cfg := config.MustLoad()

	log := common.SetupLogger(cfg.Env)
	log.Info("Запуск приложения 'библиотека'", slog.String("env", cfg.Env))
	log.Debug("Включены отладочные сообщения")

	db, err := postgres.NewPostgresDB(&cfg.DbServer)
	if err != nil {
		log.Error("Не удалось подключиться к базе postgres: ", sl.Err(err))
		log.Debug(fmt.Sprint(&cfg.DbServer))
		os.Exit(1)
	}

	log.Info("Приложение подключилось к базе данных postgres", slog.String("env", cfg.Env))

	repos := storage.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(library.Server)
	if err := srv.Run(cfg.HTTPServer.Port, handlers.InitRoutes()); err != nil {
		log.Error("Ошибка при запуске http сервера %s", err.Error())
	}

}
