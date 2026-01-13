package timer

import (
	"fmt"
	"time"
)

func ExampleTimer() {
	print := func(t Timer) {
		fmt.Printf(
			"Duration: %d\nElapsed: %d\nStarted: %t\nRunning: %t\nFinished: %t\nProgress: %.2f\n\n",
			t.Duration()/time.Second,
			t.Elapsed()/time.Second,
			t.Started(),
			t.Running(),
			t.Finished(),
			t.Progress(),
		)
	}

	var t Timer
	print(t)

	t.Set(10 * time.Second)
	print(t)

	start := time.Now()
	t.Start(start)
	print(t)

	t.Tick(start.Add(3 * time.Second))
	print(t)

	t.Stop(start.Add(5 * time.Second))
	print(t)

	t.Start(start.Add(6 * time.Second))
	print(t)

	t.Tick(start.Add(11 * time.Second))
	print(t)

	t.Reset()
	print(t)

	// Output:
	// Duration: 0
	// Elapsed: 0
	// Started: false
	// Running: false
	// Finished: false
	// Progress: 0.00
	//
	// Duration: 10
	// Elapsed: 0
	// Started: false
	// Running: false
	// Finished: false
	// Progress: 0.00
	//
	// Duration: 10
	// Elapsed: 0
	// Started: true
	// Running: true
	// Finished: false
	// Progress: 0.00
	//
	// Duration: 10
	// Elapsed: 3
	// Started: true
	// Running: true
	// Finished: false
	// Progress: 0.30
	//
	// Duration: 10
	// Elapsed: 5
	// Started: true
	// Running: false
	// Finished: false
	// Progress: 0.50
	//
	// Duration: 10
	// Elapsed: 5
	// Started: true
	// Running: true
	// Finished: false
	// Progress: 0.50
	//
	// Duration: 10
	// Elapsed: 10
	// Started: true
	// Running: false
	// Finished: true
	// Progress: 1.00
	//
	// Duration: 10
	// Elapsed: 0
	// Started: false
	// Running: false
	// Finished: false
	// Progress: 0.00
	//
}
