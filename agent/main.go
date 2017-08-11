package main

import (
	"flag"
	"time"
)

var (
	transAddr = flag.String("trans", "127.0.0.1:7070", "transfer address")
)

func main() {
	flag.Parse()

	sender := NewSender(*transAddr)
	go sender.Start()
	sched := NewScheduler(sender.Channel())
	sched.AddMetric(MetricGroupFunc(CpuMetrics), time.Second*3)
	sched.AddMetric(MetricGroupFunc(MemMetrics), time.Second*3)
	//sched.AddMetric(NewUserMetric("/Users/fanbingxin/metric.sh"), time.Second*3)
	select {}
}
