package app

import (
	"context"
	"github.com/evgen1067/anti-bruteforce/internal/bucket"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/evgen1067/anti-bruteforce/internal/repository/psql"
	"github.com/evgen1067/anti-bruteforce/internal/rest"
	"github.com/evgen1067/anti-bruteforce/internal/service"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errs := make(chan error)

	// Подключаемся к БД
	db := psql.NewRepo(cfg)
	err := db.Connect(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	// Запускаем дырявые ведра :)
	leakyBucket := bucket.NewLeakyBucket(cfg)
	go func() {
		leakyBucket.Repeat(ctx)
	}()

	// Собираем сервисы
	services := service.NewServices(ctx, db, leakyBucket)

	// Запускаем сервер АПИ
	server := rest.NewServer(services, cfg)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errs <- err
		}
	}()

	// Выползаем при ошибке или завершении программы
	select {
	case err := <-errs:
		return err
	case <-ctx.Done():

	}

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
