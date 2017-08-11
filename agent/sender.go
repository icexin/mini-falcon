package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/icexin/mini-falcon/common"
)

type Sender struct {
	ch chan *common.Metric

	addr string
	w    *bufio.Writer
	conn net.Conn
}

func NewSender(addr string) *Sender {
	return &Sender{
		ch:   make(chan *common.Metric, 1000),
		addr: addr,
	}
}

func monConn(conn net.Conn) {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		conn.Close()
	}
}

func (s *Sender) connect() error {
	conn, err := net.Dial("tcp", s.addr)
	if err != nil {
		return err
	}
	s.conn = conn
	s.w = bufio.NewWriter(s.conn)
	go monConn(conn)
	return nil
}

func (s *Sender) flush() {
	if s.w == nil {
		return
	}

	err := s.w.Flush()
	if err != nil {
		s.conn.Close()
		s.w = nil
		s.conn = nil
	}
}

func (s *Sender) send(m *common.Metric) error {
	if s.conn == nil {
		err := s.connect()
		if err != nil {
			return err
		}
	}

	buf, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	s.w.Write(buf)
	err = s.w.WriteByte('\n')

	if err != nil {
		s.conn.Close()
		s.w = nil
		s.conn = nil
	}
	return err
}

func (s *Sender) loopsend() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.flush()
		case m := <-s.ch:
			err := s.send(m)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (s *Sender) Channel() chan *common.Metric {
	return s.ch
}

func (s *Sender) Start() {
	go s.loopsend()
}
