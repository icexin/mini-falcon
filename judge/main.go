package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/icexin/mini-falcon/common"
)

var (
	kafkaAddr = flag.String("kafka", "127.0.0.1:9092", "kafka address list")
	topic     = flag.String("topic", "falcon", "kafka topic")
	ruleFile  = flag.String("f", "rule.toml", "rule file")
)

type Rule struct {
	Expr     string
	Tags     []string
	Mails    []string
	AlarmMax int
	MatchMax int

	metric string
	op     string
	value  float64

	matchCouter int
	alarmCouter int
}

func (r *Rule) MatchExpr(metric *common.Metric) bool {
	switch r.op {
	case ">":
		return metric.Value > r.value
	case "<":
		return metric.Value < r.value
	default:
		return false
	}
}

func (r *Rule) MatchTag(metric *common.Metric) bool {
	set := make(map[string]struct{})
	for _, tag := range metric.Tag {
		set[tag] = struct{}{}
	}
	for _, tag := range r.Tags {
		_, ok := set[tag]
		if !ok {
			return false
		}
	}

	return true
}

func parseRule(r *Rule) error {
	fields := strings.Fields(r.Expr)
	if len(fields) != 3 {
		return errors.New("bad expr")
	}
	r.metric, r.op = fields[0], fields[1]
	var err error
	r.value, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return err
	}

	return nil
}

type config struct {
	Rules []*Rule
}

func loadRules() (map[string][]*Rule, error) {
	var cfg config
	_, err := toml.DecodeFile(*ruleFile, &cfg)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]*Rule)
	for _, r := range cfg.Rules {
		log.Printf("%#v", r)
		err = parseRule(r)
		if err != nil {
			return nil, err
		}
		m[r.metric] = append(m[r.metric], r)
	}
	return m, nil
}

func alarm(r *Rule, metric *common.Metric, recovery bool) {
	if recovery {
		log.Printf("[recovery] %s %v", r.Expr, metric)
	} else {
		log.Printf("[alarm] %s %v", r.Expr, metric)
	}
}

func judge(metric *common.Metric) {
	rules, ok := rulemap[metric.Metric]
	if !ok {
		return
	}

	for _, r := range rules {
		if !r.MatchTag(metric) {
			continue
		}
		if !r.MatchExpr(metric) {
			if r.matchCouter >= r.MatchMax {
				alarm(r, metric, true)
			}
			r.matchCouter = 0
			r.alarmCouter = 0
			continue
		}
		r.matchCouter++
		if r.matchCouter >= r.MatchMax && r.alarmCouter < r.AlarmMax {
			r.alarmCouter++
			alarm(r, metric, false)
		}
	}
}

var (
	rulemap = make(map[string][]*Rule)
)

func main() {
	flag.Parse()
	var err error

	rulemap, err = loadRules()
	if err != nil {
		log.Fatal(err)
	}

	kafkaList := strings.Split(*kafkaAddr, ",")
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := cluster.NewConsumer(kafkaList, "falcon-judge", []string{*topic}, config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-consumer.Messages():
			var metric common.Metric
			err := json.Unmarshal(msg.Value, &metric)
			if err != nil {
				log.Printf("%s:%s", err, msg.Value)
			}
			judge(&metric)
		case err := <-consumer.Errors():
			log.Print(err)
		}
	}
}
