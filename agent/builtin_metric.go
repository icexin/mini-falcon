package main

import (
	"time"

	"github.com/icexin/mini-falcon/common"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func CpuMetrics() []*common.Metric {
	var ret []*common.Metric
	percents, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}

	ret = append(ret, NewMetric("cpu.percent", percents[0]))
	return ret
}

func MemMetrics() []*common.Metric {
	var ret []*common.Metric
	stat, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}

	ret = append(ret, NewMetric("mem.percent", stat.UsedPercent))
	return ret
}
