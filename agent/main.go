package main

import "time"

func main() {
	sender := NewSender("127.0.0.1:8080")
	go sender.Start()
	sched := NewScheduler(sender.Channel())
	sched.AddMetric(MetricGroupFunc(CpuMetrics), time.Second*3)
	sched.AddMetric(NewUserMetric("/Users/fanbingxin/metric.sh"), time.Second*3)
	select {}
}
