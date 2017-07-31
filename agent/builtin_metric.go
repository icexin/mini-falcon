package main

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CpuMetrics() []*MetricEntry {
	var ret []*MetricEntry
	percents, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}

	ret = append(ret, NewMetricEntry("cpu.percent", percents[0]))
	return ret
}

func MemMetrics() []*MetricEntry {
	return nil
}
