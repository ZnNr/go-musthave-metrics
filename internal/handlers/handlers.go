package handlers

import (
	"errors"
	"fmt"
	"github.com/ZnNr/go-musthave-metrics.git/internal/collector"
	"github.com/go-chi/chi/v5"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

// SaveMetric Функция сохранения метрики
func SaveMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(r, "type")   // Получает значение параметра "type" из URL
	metricName := chi.URLParam(r, "name")   // Получает значение параметра "name" из URL
	metricValue := chi.URLParam(r, "value") // Получает значение параметра "value" из URL

	err := collector.Collector.Collect(metricName, metricType, metricValue)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest) // Устанавливает код ответа 400 Bad Request
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented) // Устанавливает код ответа 501 Not Implemented
		return
	}

	if _, err = io.WriteString(w, ""); err != nil {
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(metricName)))
	w.WriteHeader(http.StatusOK)
}

// GetMetric Функция получения метрики
func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type") // Получает значение параметра "type" из URL
	metricName := chi.URLParam(r, "name")

	value, err := collector.Collector.GetMetric(metricName, metricType)
	if errors.Is(err, collector.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if _, err = io.WriteString(w, ""); err != nil {
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(value)))
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, value); err != nil {
		return
	}
}

// ShowMetrics Функция отображения всех метрик
func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, fmt.Sprintf("wrong path %q", r.URL.Path), http.StatusNotFound)
		return
	}
	page := ""
	for _, n := range collector.Collector.GetAvailableMetrics() {
		page += fmt.Sprintf("<h1>	%s</h1>", n)
	}
	tmpl, _ := template.New("data").Parse("<h1>AVAILABLE METRICS</h1>{{range .}}<h3>{{ .}}</h3>{{end}}")
	if err := tmpl.Execute(w, collector.Collector.GetAvailableMetrics()); err != nil {
		return
	}
	w.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
