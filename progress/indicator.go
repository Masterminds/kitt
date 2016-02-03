package progress

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Indicator struct {
	Interval time.Duration
	Frames   []string
	Writer   io.Writer

	initialMessage string
	done           chan bool
	msg            chan string
}

func (i *Indicator) Start(message string) {
	i.done = make(chan bool)
	i.msg = make(chan string, 1)
	go i.run()
	i.msg <- message
}

func (i *Indicator) Message(message string) {
	i.msg <- message
}

func (i *Indicator) Done() {
	i.done <- true
}

func (i *Indicator) run() {
	tt := time.NewTicker(i.Interval)
	fr := 0
	mm := i.initialMessage
	maxlen := len(mm)

	for {
		select {
		case m := <-i.msg:
			mm = m
		case <-tt.C:
			pad := maxlen - len(mm)
			if pad < 0 {
				maxlen = len(mm)
			} else if pad > 0 {
				mm += strings.Repeat(" ", pad)
			}
			fmt.Fprintf(i.Writer, "%s %s\r", i.Frames[fr], mm)
			fr++
			if fr >= len(i.Frames) {
				fr = 0
			}
		case <-i.done:
			tt.Stop()
			close(i.done)
			close(i.msg)
		}
	}
}
