package main

import (
	"flag"
	"runtime"
	"strings"
	"time"
)

var (
	transAddr = flag.String("trans", "127.0.0.1:7070", "transfer address")
	tagflag   = flag.String("tags", "", "tags")

	tags []string
)

func main() {
	flag.Parse()

	if *tagflag != "" {
		tags = strings.Split(*tagflag, ",")
	}
	tags = append(tags, runtime.GOOS)

	sender := NewSender(*transAddr)
	go sender.Start()
	sched := NewScheduler(sender.Channel())
	sched.AddMetric(MetricGroupFunc(CpuMetrics), time.Second*3)
	sched.AddMetric(MetricGroupFunc(MemMetrics), time.Second*3)
	select {}
}
