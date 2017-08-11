package main

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
	var ret []*MetricEntry
	stat, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}

	ret = append(ret, NewMetricEntry("mem.used", stat.UsedPercent))
	return ret
}
