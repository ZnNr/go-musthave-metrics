package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZnNr/go-musthave-metrics.git/internal/collector"
	"github.com/go-chi/chi/v5"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

func SaveMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	err := collector.Collector.Collect(metricName, metricType, metricValue)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if _, err = io.WriteString(w, ""); err != nil {
		return
	}

	if _, err = io.WriteString(w, ""); err != nil {
		return
	}
	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(metricName)))
	w.WriteHeader(http.StatusOK)
}

func SaveMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metric Metrics
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metricValue := ""
	switch metric.MType {
	case "counter":
		metricValue = strconv.Itoa(int(*metric.Delta))
	case "gauge":
		metricValue = fmt.Sprintf("%.11f", *metric.Value)
	}

	err := collector.Collector.Collect(metric.ID, metric.MType, metricValue)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	updated, err := collector.Collector.GetMetricByName(metric.ID, metric.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := Metrics{
		ID:    metric.ID,
		MType: metric.MType,
	}
	switch result.MType {
	case "counter":
		c, err := strconv.Atoi(updated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c64 := int64(c)
		result.Delta = &c64
	case "gauge":
		g, err := strconv.ParseFloat(updated, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result.Value = &g
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	if _, err = w.Write(resultJSON); err != nil {
		return
	}
	w.Header().Set("content-length", strconv.Itoa(len(metric.ID)))
	w.WriteHeader(http.StatusOK)
}

func GetMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metric Metrics
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := collector.Collector.GetMetricByName(metric.ID, metric.MType)
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
	switch metric.MType {
	case "counter":
		c, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c64 := int64(c)
		metric.Delta = &c64
	case "gauge":
		g, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metric.Value = &g
	}
	resultJSON, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	if _, err = w.Write(resultJSON); err != nil {
		return
	}
	w.Header().Set("content-length", strconv.Itoa(len(value)))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	value, err := collector.Collector.GetMetricByName(metricName, metricType)
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
	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(value)))
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, value); err != nil {
		return
	}
}

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

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
