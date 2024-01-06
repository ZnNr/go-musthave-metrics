package main

import (
	"github.com/ZnNr/go-musthave-metrics.git/internal/flags"
	"github.com/ZnNr/go-musthave-metrics.git/internal/handlers"
	log "github.com/ZnNr/go-musthave-metrics.git/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment() // добавляем предустановленный логер NewDevelopment
	if err != nil {                     // вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	log.SugarLogger = *logger.Sugar()

	params := flags.Init(flags.WithAddr())
	r := chi.NewRouter()
	r.Use(log.RequestLogger)
	r.Post("/update/", handlers.SaveMetricFromJSON)
	r.Post("/value/", handlers.GetMetricFromJSON)
	r.Post("/update/{type}/{name}/{value}", handlers.SaveMetric)
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Get("/", handlers.ShowMetrics)

	log.SugarLogger.Infow(
		"Starting server",
		"addr", params.FlagRunAddr,
	)
	if err := http.ListenAndServe(params.FlagRunAddr, r); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}
