package log_monitor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ----- 抽象层 -----
type Reader interface {
	Read(rc chan []byte)
}
type Writer interface {
	Write(wc chan string)
}

// ----- 实现层 -----
type LogMonitor struct {
	rc     chan []byte
	wc     chan string
	reader Reader
	writer Writer
}

func (l *LogMonitor) Parse() {
	for s := range l.rc {
		str := strings.ToUpper(string(s))
		l.wc <- str
	}
}

type MonitorReader struct {
	path string
}

func (m *MonitorReader) Read(rc chan []byte) {
	f, err := os.Open(m.path)
	if err != nil {
		panic(fmt.Sprintf("open file err: %s", err.Error()))
	}

	// 偏移 从末尾读取
	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes err: %s", err.Error()))
		}

		rc <- line
	}
}

type MonitorWriter struct {
	path string
}

func (m *MonitorWriter) Write(c chan string) {
	f, err := os.OpenFile(m.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("open file err: %s", err.Error()))
	}
	defer f.Close()

	for str := range c {
		_, err := f.WriteString(str)
		if err != nil {
			fmt.Printf("write err: %s", err.Error())
			continue
		}
	}
}

// ----- 业务逻辑层 -----
func Run() {
	rc := make(chan []byte)
	wc := make(chan string)

	reader := &MonitorReader{
		path: "./access.log",
	}
	writer := &MonitorWriter{
		path: "./output.log",
	}

	l := &LogMonitor{
		rc:     rc,
		wc:     wc,
		reader: reader,
		writer: writer,
	}

	go l.reader.Read(l.rc)
	go l.Parse()
	go l.writer.Write(l.wc)

	for {
		time.Sleep(time.Second) // 阻塞主线程
	}
}
