package handlers

import (
	"fmt"
	"github.com/ZnNr/go-musthave-metrics.git/internal/collector"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveMetric(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", SaveMetric)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name          string
		mType         string
		mName         string
		mValue        string
		expectedCode  int
		expectedError error
	}{
		{
			name:         "case0",
			mType:        "counter",
			mName:        "Counter1",
			mValue:       "15",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case1",
			mType:        "gauge",
			mName:        "Gauge1",
			mValue:       "12.282",
			expectedCode: http.StatusOK,
		},
		{
			name:          "case2",
			mType:         "invalid",
			mName:         "Gauge1",
			mValue:        "12.282",
			expectedCode:  http.StatusNotImplemented,
			expectedError: collector.ErrNotImplemented,
		},
		{
			name:          "case3",
			mType:         "counter",
			mName:         "Counter1",
			mValue:        "15.2562",
			expectedCode:  http.StatusBadRequest,
			expectedError: collector.ErrNotFound,
		},
		{
			name:          "case4",
			mType:         "gauge",
			mName:         "Gauge1",
			mValue:        "12.282dgh",
			expectedCode:  http.StatusBadRequest,
			expectedError: collector.ErrNotFound,
		},
		{
			name:          "case5",
			mType:         "gauge",
			mName:         "Gauge1",
			mValue:        "",
			expectedCode:  http.StatusNotFound,
			expectedError: collector.ErrNotFound,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("%s/update/%s/%s/%s", srv.URL, tt.mType, tt.mName, tt.mValue))

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)

			value, err := collector.Collector.GetMetricByName(tt.mName, tt.mType)
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, value, tt.mValue)
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", SaveMetric)
	r.Get("/value/{type}/{name}", GetMetric)
	srv := httptest.NewServer(r)
	defer srv.Close()

	client := resty.New()
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter3/15", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter2/0", srv.URL))

	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge1/100500.2780001", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge2/100500.278000100", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge3/100500", srv.URL))

	testCases := []struct {
		name          string
		mType         string
		mName         string
		mValue        string
		expectedCode  int
		expectedError error
	}{
		{
			name:         "case0",
			mType:        "counter",
			mName:        "Counter3",
			mValue:       "15",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case1",
			mType:        "counter",
			mName:        "Counter2",
			mValue:       "0",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case2",
			mType:        "gauge",
			mName:        "Gauge1",
			mValue:       "100500.2780001",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case3",
			mType:        "gauge",
			mName:        "Gauge2",
			mValue:       "100500.278000100",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case4",
			mType:        "gauge",
			mName:        "Gauge3",
			mValue:       "100500",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case5",
			mType:        "gauge",
			mName:        "Gauge4",
			mValue:       "",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "case5",
			mType:        "invalid",
			mName:        "Gauge4",
			mValue:       "",
			expectedCode: http.StatusNotImplemented,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Get(fmt.Sprintf("%s/value/%s/%s", srv.URL, tt.mType, tt.mName))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, string(resp.Body()), tt.mValue)
		})
	}
}

func TestShowMetrics(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", SaveMetric)
	r.Get("/", ShowMetrics)
	srv := httptest.NewServer(r)
	defer srv.Close()

	client := resty.New()
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter3/15", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter2/0", srv.URL))

	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge1/100500.2780001", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge2/100500.278000100", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge3/100500", srv.URL))

	testCases := []struct {
		name         string
		expectedPage string
		expectedCode int
	}{
		{
			name:         "case0",
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Get(fmt.Sprintf("%s/", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
		})
	}
}
