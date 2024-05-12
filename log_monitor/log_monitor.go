package log_monitor

import (
	"fmt"
	"strings"
	"time"
)

// ----- 抽象层 -----
type Reader interface {
	Read(c chan string)
}
type Writer interface {
	Write(c chan string)
}

// ----- 实现层 -----
type LogMonitor struct {
	rc     chan string
	wc     chan string
	reader Reader
	writer Writer
}

type MonitorReader struct{}

func (m *MonitorReader) Read(c chan string) {
	s := "message"
	c <- s
}

type MonitorWriter struct{}

func (m *MonitorWriter) Write(c chan string) {
	fmt.Println(<-c)
}

func (l *LogMonitor) Parse() {
	s := <-l.rc
	s = strings.ToUpper(s)
	l.wc <- s
}

// ----- 业务逻辑层 -----
func Run() {
	rc := make(chan string)
	wc := make(chan string)

	reader := &MonitorReader{}
	writer := &MonitorWriter{}

	l := &LogMonitor{
		rc:     rc,
		wc:     wc,
		reader: reader,
		writer: writer,
	}

	go l.reader.Read(l.rc)
	go l.Parse()
	go l.writer.Write(l.wc)

	time.Sleep(1 * time.Second)
}
