package diagnostics

import "time"

type Stopwatch struct {
	startTime time.Time
}

func (stw *Stopwatch) Reset() {
	stw.startTime = time.Now()
}

func (stw *Stopwatch) Microseconds() int64 {
	return time.Since(stw.startTime).Microseconds()
}

func (stw *Stopwatch) Seconds() float64 {
	return time.Since(stw.startTime).Seconds()
}

func (stw *Stopwatch) Milliseconds() int64 {
	return time.Since(stw.startTime).Milliseconds()
}
