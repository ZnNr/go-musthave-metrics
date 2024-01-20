package file

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/ZnNr/go-musthave-metrics.git/internal/collector"
	"github.com/ZnNr/go-musthave-metrics.git/internal/flags"
	"os"
)

func (m *manager) Restore(ctx context.Context) ([]collector.MetricJSON, error) {
	file, err := os.OpenFile(m.fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, err
	}

	data := scanner.Bytes()
	var metricsFromFile []collector.MetricJSON
	if err = json.Unmarshal(data, &metricsFromFile); err != nil {
		return nil, err
	}
	return metricsFromFile, nil
}

func (m *manager) Save(ctx context.Context, metrics []collector.MetricJSON) error {
	var saveError error
	file, err := os.OpenFile(m.fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			saveError = err
		}
	}()

	writer := bufio.NewWriter(file)

	data, err := json.Marshal(&metrics)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return saveError
}

func New(params *flags.Params) *manager {
	return &manager{fileName: params.FileStoragePath}
}

type manager struct {
	fileName string
}
