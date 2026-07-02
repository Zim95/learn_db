package profiler

import "time"

/*
	Timer measures execution time and nothing else.

	Usage:

		t := CreateTimer("set")
		t.Start()
		// ... run the operation ...
		result := t.Stop()   // -> TimerResult

	The caller (usually the Profiler orchestrator) is responsible for pushing
	the returned TimerResult into the Collector. The Timer never prints and
	never stores.
*/

// TimerResult holds how long an operation took.
type TimerResult struct {
	Op       string
	Duration time.Duration
}

func (t TimerResult) Operation() string {
	return t.Op
}

type Timer struct {
	op    string
	start time.Time
}

func CreateTimer(op string) *Timer {
	return &Timer{
		op: op,
	}
}

// Start records the moment measurement begins.
func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() TimerResult {
	return TimerResult{
		Op:       t.op,
		Duration: time.Since(t.start),
	}
}
