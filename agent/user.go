package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/icexin/mini-falcon/common"
)

type UserMetric struct {
	script string
}

func NewUserMetric(script string) *UserMetric {
	return &UserMetric{
		script: script,
	}
}

func (u *UserMetric) Metrics() []*common.Metric {
	var ret []*common.Metric
	cmd := exec.Command("bash", "-c", u.script)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Print(err)
		return ret
	}

	r := bufio.NewReader(stdout)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		name := fields[0]
		value, _ := strconv.ParseFloat(fields[1], 64)
		metric := NewMetric(name, value)
		ret = append(ret, metric)
	}

	err = cmd.Wait()
	if err != nil {
		out, _ := ioutil.ReadAll(stderr)
		log.Printf("%s:%s", out, err)
	}

	return ret
}
