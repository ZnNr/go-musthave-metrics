package storage

import (
	"fmt"
	"github.com/ZnNr/go-musthave-metrics.git/internal/collector"
	"math/rand"
	"runtime"
	"strconv"
)

// Store функуия используется для сбора метрик и сохранения их в хранилище.
func (a *storage) Store() {
	metrics := runtime.MemStats{}  //создается переменная metrics типа runtime.MemStats, которая представляет собой статистику памяти
	runtime.ReadMemStats(&metrics) //вызывается функциякоторая заполняет структуру metrics актуальными данными о памяти.
	//Происходит вызов метода Collect() объекта a.metricsCollector для каждого из собранных показателей памяти.
	//Каждый вызов передает имя метрики, тип и значение метрики, преобразованное в строку.
	a.metricsCollector.Collect("Alloc", "gauge", strconv.FormatUint(metrics.Alloc, 10)) //собирает метрику "Alloc" и передает ее в хранилище. Значение метрики образуется из поля metrics.Alloc, которое представляет количество выделенной памяти для объектов Go
	a.metricsCollector.Collect("BuckHashSys", "gauge", strconv.FormatUint(metrics.BuckHashSys, 10))
	a.metricsCollector.Collect("Frees", "gauge", strconv.FormatUint(metrics.Frees, 10))
	a.metricsCollector.Collect("GCCPUFraction", "gauge", fmt.Sprintf("%.3f", metrics.GCCPUFraction))
	a.metricsCollector.Collect("GCSys", "gauge", strconv.FormatUint(metrics.GCSys, 10))
	a.metricsCollector.Collect("HeapAlloc", "gauge", strconv.FormatUint(metrics.HeapAlloc, 10))
	a.metricsCollector.Collect("HeapIdle", "gauge", strconv.FormatUint(metrics.HeapIdle, 10))
	a.metricsCollector.Collect("HeapInuse", "gauge", strconv.FormatUint(metrics.HeapInuse, 10))
	a.metricsCollector.Collect("HeapObjects", "gauge", strconv.FormatUint(metrics.HeapObjects, 10))
	a.metricsCollector.Collect("HeapReleased", "gauge", strconv.FormatUint(metrics.HeapReleased, 10))
	a.metricsCollector.Collect("HeapSys", "gauge", strconv.FormatUint(metrics.HeapSys, 10))
	a.metricsCollector.Collect("Lookups", "gauge", strconv.FormatUint(metrics.Lookups, 10))
	a.metricsCollector.Collect("MCacheInuse", "gauge", strconv.FormatUint(metrics.MCacheInuse, 10))
	a.metricsCollector.Collect("MCacheSys", "gauge", strconv.FormatUint(metrics.MCacheSys, 10))
	a.metricsCollector.Collect("MSpanInuse", "gauge", strconv.FormatUint(metrics.MSpanInuse, 10))
	a.metricsCollector.Collect("MSpanSys", "gauge", strconv.FormatUint(metrics.MSpanSys, 10))
	a.metricsCollector.Collect("Mallocs", "gauge", strconv.FormatUint(metrics.Mallocs, 10))
	a.metricsCollector.Collect("NextGC", "gauge", strconv.FormatUint(metrics.NextGC, 10))
	a.metricsCollector.Collect("NumForcedGC", "gauge", strconv.Itoa(int(metrics.NumForcedGC)))
	a.metricsCollector.Collect("NumGC", "gauge", strconv.FormatUint(uint64(metrics.NumGC), 10))
	a.metricsCollector.Collect("OtherSys", "gauge", strconv.Itoa(int(metrics.OtherSys)))
	a.metricsCollector.Collect("PauseTotalNs", "gauge", strconv.Itoa(int(metrics.PauseTotalNs)))
	a.metricsCollector.Collect("StackInuse", "gauge", strconv.Itoa(int(metrics.StackInuse)))
	a.metricsCollector.Collect("StackSys", "gauge", strconv.Itoa(int(metrics.StackSys)))
	a.metricsCollector.Collect("Sys", "gauge", strconv.Itoa(int(metrics.Sys)))
	a.metricsCollector.Collect("TotalAlloc", "gauge", strconv.Itoa(int(metrics.TotalAlloc)))
	a.metricsCollector.Collect("RandomValue", "gauge", strconv.Itoa(rand.Int()))
	a.metricsCollector.Collect("LastGC", "gauge", strconv.Itoa(int(metrics.LastGC)))

	cnt, _ := collector.Collector.GetMetricByName("PollCount", "counter")
	v, _ := strconv.Atoi(cnt)
	collector.Collector.Collect("PollCount", "counter", strconv.Itoa(v+1))
}

// New - это конструктор, который создает и возвращает новый экземпляр структуры storage.
// Он принимает аргумент metricsCollector, который должен быть реализацией интерфейса collectorImpl
func New(metricsCollector collectorImpl) *storage {
	return &storage{
		metricsCollector: metricsCollector,
	}
}

// В структуре storage определены два поля:
//
// metricsCollector - тип этого поля задан как collectorImpl, это поле будет использоваться для сбора и хранения метрик.
// полю metricsCollector можно присвоить любое значение, которое соответствует интерфейсу collectorImpl.
type storage struct {
	metricsCollector collectorImpl
}

// Интерфейс collectorImpl определяет только один метод Collect, который принимает три аргумента: metricName (имя метрики), metricType (тип метрики) и metricValue (значение метрики), и возвращает ошибку
type collectorImpl interface {
	Collect(metricName string, metricType string, metricValue string) error
}
