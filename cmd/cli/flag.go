package main

import "flag"

type options struct {
	timeout int
	format  string
}

var ops = options{
	timeout: 0,    // Default timeout
	format:  "dt", // Default format DateTime
}

func initFlags() {
	flag.IntVar(&ops.timeout, "timeout", 0, "Timeout in seconds for the operation")
	flag.IntVar(&ops.timeout, "t", 0, "Timeout in seconds for the operation (shorthand)")

	var q bool
	flag.BoolVar(&q, "quick", false, "Quick mode, 10 second timeout")
	flag.BoolVar(&q, "q", false, "Quick mode, 10 second timeout (shorthand)")

	var m bool
	flag.BoolVar(&m, "mid", false, "Mid mode, 30 second timeout")
	flag.BoolVar(&m, "m", false, "Mid mode, 30 second timeout (shorthand)")

	var i bool
	flag.BoolVar(&i, "infinite", false, "Infinite mode, no timeout")
	flag.BoolVar(&i, "i", false, "Infinite mode, no timeout (shorthand)")

	flag.StringVar(&ops.format, "format", "dt", "Output format (default: dt)")
	flag.StringVar(&ops.format, "f", "dt", "Output format (shorthand, default: dt)")

	flag.Parse()

	if i {
		ops.timeout = 0
	} else if q {
		ops.timeout = 10
	} else if m {
		ops.timeout = 30
	}

	switch ops.format {
	case "dt", "d", "t":
	default:
		ops.format = "dt"
	}
}
