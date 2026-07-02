package profiler

import "runtime"

/*
	MemoryProfiler captures allocation activity across an operation.

	Usage:

		m := CreateMemoryProfiler("set")
		m.Start()
		// ... run the operation ...
		result := m.Stop()   // -> MemoryResult

	How it works:

	runtime.ReadMemStats fills a MemStats snapshot. We take one snapshot before
	the operation and one after, then subtract.

	  - TotalAlloc and Mallocs are monotonic counters (they only ever go up),
	    so their deltas are always >= 0 and directly answer "how much did this
	    operation allocate?".
	  - HeapAlloc is *live* heap, which the garbage collector can shrink at any
	    time. Its delta is therefore signed and can be negative if a GC ran
	    mid-operation.

	Note: ReadMemStats stops the world briefly, so this is not free. That is
	fine for a one-shot CLI, but you would not call it inside a hot loop.
*/

// MemoryResult holds allocation activity across an operation.
//
// TotalAlloc and Mallocs are cumulative counters in the Go runtime, so their
// deltas are always >= 0 and are the honest numbers to reason about.
// HeapAllocDelta is signed on purpose: the garbage collector can run mid
// operation and free memory, which would make a "live heap" delta negative.
type MemoryResult struct {
	Op             string
	TotalAlloc     uint64 // bytes allocated by this op (cumulative counter delta)
	Mallocs        uint64 // number of heap objects allocated by this op
	HeapAllocDelta int64  // change in live heap bytes (may be negative after a GC)
}

func (m MemoryResult) Operation() string {
	return m.Op
}

type MemoryProfiler struct {
	op     string
	before runtime.MemStats
}

func CreateMemoryProfiler(op string) *MemoryProfiler {
	return &MemoryProfiler{
		op: op,
	}
}

// Start captures the "before" memory snapshot.
func (m *MemoryProfiler) Start() {
	runtime.ReadMemStats(&m.before) // initialize the stats value
}

// Stop captures the "after" snapshot and returns the difference.
func (m *MemoryProfiler) Stop() MemoryResult {
	var after runtime.MemStats
	runtime.ReadMemStats(&after) // initialize the stats value

	return MemoryResult{
		Op:             m.op,
		TotalAlloc:     after.TotalAlloc - m.before.TotalAlloc,
		Mallocs:        after.Mallocs - m.before.Mallocs,
		HeapAllocDelta: int64(after.HeapAlloc) - int64(m.before.HeapAlloc),
	}
}
