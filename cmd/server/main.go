package main

import (
	"github.com/ZnNr/go-musthave-metrics.git/internal/metricshandlers"
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/update/", metricshandlers.SaveMetric)
	// запуск сервера с адресом localhost, порт 8080
	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		//Если при запуске сервера возникает какая-либо ошибка, она фиксируется и поднимается паника.
		panic(err)
	}
}
