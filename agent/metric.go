package main

import (
	"os"
	"time"

	"github.com/icexin/mini-falcon/common"
)

type MetricGroup interface {
	Metrics() []*common.Metric
}

type MetricGroupFunc func() []*common.Metric

func (m MetricGroupFunc) Metrics() []*common.Metric {
	return m()
}

func NewMetric(name string, value float64) *common.Metric {
	host, _ := os.Hostname()
	return &common.Metric{
		Metric:    name,
		Endpoint:  host,
		Value:     value,
		Tag:       tags,
		Timestamp: time.Now().Unix(),
	}
}

type Scheduler struct {
	output chan *common.Metric
}

func NewScheduler(output chan *common.Metric) *Scheduler {
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
