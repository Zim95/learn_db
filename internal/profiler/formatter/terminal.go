package formatter

import (
	"fmt"
	"io"

	"learn_db/internal/profiler"
)

/*
	Terminal formatting.

	Turns a []profiler.Result into a human-readable summary written to any
	io.Writer (usually os.Stderr, so it stays separate from the command's
	real output on stdout). It performs no measurement and does no storage;
	it only knows how to render each result type.
*/

func Terminal(w io.Writer, results []profiler.Result) {
	if len(results) == 0 {
		return
	}

	fmt.Fprintln(w, "── profiling ─────────────────────────────")

	for _, result := range results {
		switch r := result.(type) {

		case profiler.TimerResult:
			fmt.Fprintf(w, "%-8s time   %s\n", r.Operation(), r.Duration)

		case profiler.MemoryResult:
			fmt.Fprintf(
				w,
				"%-8s mem    alloc=%dB objects=%d heap-delta=%dB\n",
				r.Operation(),
				r.TotalAlloc,
				r.Mallocs,
				r.HeapAllocDelta,
			)

		case profiler.CPUResult:
			fmt.Fprintf(
				w,
				"%-8s cpu    profile=%s  (analyze: go tool pprof %s)\n",
				r.Operation(),
				r.ProfilePath,
				r.ProfilePath,
			)

		default:
			fmt.Fprintf(w, "%-8s (unknown result type %T)\n", r.Operation(), r)
		}
	}
}
