package main

import (
	"fmt"
	"time"

	"github.com/seannyphoenix/binarytime/pkg/binarytime"
)

func runClock() {
	ticker := time.NewTicker(time.Millisecond * 250)
	defer ticker.Stop()

	done := make(chan bool)
	defer close(done)

	go func() {
		if ops.timeout > 0 {
			time.Sleep(time.Duration(ops.timeout) * time.Second)
			done <- true
		}
	}()

	var currBTime string
	var currTime string
	var updated bool

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			tStr := formatTime(t)
			if tStr != currTime {
				currTime = tStr
				updated = true
			}

			bt := binarytime.FromTime(t)
			btStr := formatBTime(bt)
			if btStr != currBTime {
				currBTime = btStr
				updated = true
			}

			if updated {
				fmt.Print("\033[H\033[2J")
				fmt.Println(currBTime)
				fmt.Println(currTime)
				updated = false
			}
		}
	}
}

func formatTime(t time.Time) string {
	switch ops.format {
	case "d":
		return t.Format("2006-01-02")
	case "t":
		return t.Format("15:04:05")
	case "dt":
		fallthrough
	default:
		return t.Format("2006-01-02\n15:04:05")
	}
}

func formatBTime(bt binarytime.Date) string {
	switch ops.format {
	case "d":
		return bt.DateGlyphs()
	case "t":
		return bt.TimeGlyphs()
	case "dt":
		fallthrough
	default:
		return bt.DateTimeGlyphs()
	}
}
