package main

import (
	"os"
	"time"
)

type MetricGroup interface {
	Metrics() []*MetricEntry
}

type MetricGroupFunc func() []*MetricEntry

func (m MetricGroupFunc) Metrics() []*MetricEntry {
	return m()
}

type MetricEntry struct {
	Metric    string  `json:"metric"`
	Endpoint  string  `json:"endpoint"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

func NewMetricEntry(name string, value float64) *MetricEntry {
	host, _ := os.Hostname()
	return &MetricEntry{
		Metric:    name,
		Endpoint:  host,
		Value:     value,
		Timestamp: time.Now().Unix(),
	}
}

type Scheduler struct {
	output chan *MetricEntry
}

func NewScheduler(output chan *MetricEntry) *Scheduler {
	return &Scheduler{
		output: output,
	}
}

func (s *Scheduler) runMetric(m MetricGroup, circle time.Duration) {
	ticker := time.NewTicker(circle)
	defer ticker.Stop()
	for range ticker.C {
		metrics := m.Metrics()
		for _, m := range metrics {
			s.output <- m
		}
	}
}

func (s *Scheduler) AddMetric(m MetricGroup, circle time.Duration) {
	go s.runMetric(m, circle)
}
