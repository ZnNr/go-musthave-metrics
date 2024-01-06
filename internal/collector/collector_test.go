package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name          string     // name - название тестового случая.
		storage       collector  // storage - экземпляр коллектора, который тестируется.
		metricName    string     // metricName - название собираемой метрики.
		metricType    string     // metricType - тип собираемой метрики.
		metricValue   string     // metricValue - значение собираемой метрики.
		expected      memStorage // expected - ожидаемое состояние storage после сбора метрики.
		expectedError error      // expectedError - ожидаемая ошибка после сбора метрики.
	}{
		{
			name:        "case0",
			storage:     collector{storage: &memStorage{gauges: map[string]string{}, counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "1",
			expected: memStorage{
				gauges: map[string]string{
					"Alloc": "1",
				},
			},
		},
		{
			name: "case1",
			storage: collector{storage: &memStorage{gauges: map[string]string{
				"Alloc":         "3",
				"GCCPUFraction": "5.543",
			}, counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "1",
			expected: memStorage{
				gauges: map[string]string{
					"Alloc":         "1",
					"GCCPUFraction": "5.543",
				},
			},
		},
		{
			name: "case3",
			storage: collector{storage: &memStorage{gauges: map[string]string{
				"Alloc": "3",
				"Sys":   "5",
			}, counters: map[string]int{
				"Counter": 5,
			}}},
			metricName:  "Counter",
			metricType:  "counter",
			metricValue: "10",
			expected: memStorage{
				gauges: map[string]string{
					"Alloc": "3",
					"Sys":   "5",
				},
				counters: map[string]int{
					"Counter": 15,
				},
			},
		},
		{
			name: "case4",
			storage: collector{storage: &memStorage{gauges: map[string]string{
				"Alloc": "3",
				"Sys":   "5",
			}, counters: map[string]int{}}},
			metricName:  "Counter",
			metricType:  "counter",
			metricValue: "10",
			expected: memStorage{
				gauges: map[string]string{
					"Alloc": "3",
					"Sys":   "5",
				},
				counters: map[string]int{
					"Counter": 10,
				},
			},
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{gauges: map[string]string{}, counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "1.0000000",
			expected: memStorage{
				gauges: map[string]string{
					"Alloc": "1.0000000",
				},
			},
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{gauges: map[string]string{}, counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "invalid",
			expected: memStorage{
				map[string]int{},
				map[string]string{},
			},
			expectedError: ErrBadRequest,
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{gauges: map[string]string{}, counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "invalid",
			metricValue: "15",
			expected: memStorage{
				map[string]int{},
				map[string]string{},
			},
			expectedError: ErrNotImplemented,
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{gauges: map[string]string{}, counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "counter",
			metricValue: "17.001",
			expected: memStorage{
				map[string]int{},
				map[string]string{},
			},
			expectedError: ErrBadRequest,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Collect(tt.metricName, tt.metricType, tt.metricValue)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
			assert.Equal(t, tt.expected.gauges, tt.storage.GetGauges())
		})
	}
}

func TestCollector_GetMetric(t *testing.T) {
	_ = Collector.Collect("Counter1", "counter", "0")
	_ = Collector.Collect("Counter2", "counter", "15")
	_ = Collector.Collect("Gauge1", "gauge", "17.01")
	_ = Collector.Collect("Gauge2", "gauge", "18.00000")

	testCases := []struct {
		name          string
		metricName    string
		metricType    string
		expectedValue string
		expectedError error
	}{
		{
			name:          "case0",
			metricType:    "counter",
			metricName:    "Counter1",
			expectedValue: "0",
		},
		{
			name:          "case1",
			metricType:    "counter",
			metricName:    "Counter2",
			expectedValue: "15",
		},
		{
			name:          "case2",
			metricType:    "gauge",
			metricName:    "Gauge1",
			expectedValue: "17.01",
		},
		{
			name:          "case3",
			metricType:    "gauge",
			metricName:    "Gauge2",
			expectedValue: "18.00000",
		},
		{
			name:          "case4",
			metricType:    "gauge",
			metricName:    "Gauge3",
			expectedValue: "",
			expectedError: ErrNotFound,
		},
		{
			name:          "case4",
			metricType:    "invalid",
			metricName:    "Gauge2",
			expectedValue: "",
			expectedError: ErrNotImplemented,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := Collector.GetMetricByName(tt.metricName, tt.metricType)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
			assert.Equal(t, value, tt.expectedValue)
		})
	}
}
