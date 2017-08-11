package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"text/template"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/icexin/mini-falcon/common"
)

var (
	kafkaAddr = flag.String("kafka", "127.0.0.1:9092", "kafka address list")
	topic     = flag.String("topic", "falcon-alarm", "kafka alarm topic")
)

var mailTpl = template.Must(template.New("alarm").Parse(`From: <falcon@reboot.com>
To: <{{.To}}>
Subject: "mail alarm"
Content-type:text/html;charset=utf8

<table border=1>
<tr>
<th>是否恢复</th>
<th>机器</th>
<th>表达式</th>
<th>值</th>
<th>标签</th>
</tr>

<tr>
<td>{{.Recovery}}</td>
<td>{{.Endpoint}}</td>
<td>{{.Expr}}</td>
<td>{{.Value}}</td>
<td>{{.Tag}}</td>
</tr>

</table>
`))

type alarmBody struct {
	To       string
	Recovery bool
	Endpoint string
	Expr     string
	Value    float64
	Tag      string
}

func mail(body *alarmBody) error {
	buf := new(bytes.Buffer)
	err := mailTpl.Execute(buf, body)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("sendmail", "-t")
	cmd.Stdin = buf
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s:%s", err, out)
	}

	return nil
}

func sendAlarm(alarm *common.Alarm) {
	body := &alarmBody{
		Recovery: alarm.Recovery,
		Endpoint: alarm.Endpoint,
		Expr:     alarm.Expr,
		Value:    alarm.Value,
		Tag:      strings.Join(alarm.Tag, ","),
	}

	log.Print(body)
	for _, addr := range alarm.Mail {
		body.To = addr
		err := mail(body)
		if err != nil {
			log.Print(err)
		}
	}
}

func main() {
	flag.Parse()
	kafkaList := strings.Split(*kafkaAddr, ",")

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := cluster.NewConsumer(kafkaList, "falcon-alarm", []string{*topic}, config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-consumer.Messages():
			var alarm common.Alarm
			err = json.Unmarshal(msg.Value, &alarm)
			if err != nil {
				log.Printf("%s:%s", msg.Value, err)
			}
			sendAlarm(&alarm)
		case err := <-consumer.Errors():
			log.Print(err)
		}
	}

}
