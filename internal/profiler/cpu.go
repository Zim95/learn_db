package profiler

import (
	"fmt"
	"os"
	"runtime/pprof"
)

/*
	CPUProfiler wraps runtime/pprof CPU profiling.

	Usage:

		c := CreateCPUProfiler("set", "cpu.prof")
		if err := c.Start(); err != nil { ... }
		// ... run the operation ...
		result := c.Stop()   // -> CPUResult, then: go tool pprof cpu.prof

	IMPORTANT CAVEAT (read this before trusting the output):

	pprof CPU profiling works by *sampling* the call stack ~100 times per
	second. A single SetKey/GetKey on this database finishes in microseconds,
	so it will almost certainly capture ZERO samples and the profile will be
	empty.

	CPU profiling is only meaningful when it wraps something that runs for a
	meaningful stretch of wall-clock time, for example:
	  - building the index over a large log file, or
	  - a benchmark loop (`go test -bench`) that runs the op thousands of times.

	Also note pprof CPU profiling is process-global: only one can be active at
	a time. So the Profiler may wrap exactly one region in CPU profiling per run.
*/

// CPUResult points at the pprof profile written to disk for an operation.
// There is no in-memory summary: you analyze it with `go tool pprof <path>`.
type CPUResult struct {
	Op          string
	ProfilePath string
}

func (c CPUResult) Operation() string {
	return c.Op
}

type CPUProfiler struct {
	op   string
	path string
	file *os.File
}

func CreateCPUProfiler(op string, path string) *CPUProfiler {
	return &CPUProfiler{
		op:   op,
		path: path,
	}
}

// Start creates the profile file and begins CPU profiling. It can fail (the
// file may not be creatable), so unlike the other profilers' Start it returns
// an error.
func (c *CPUProfiler) Start() error {
	file, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("failed to create cpu profile file: %w", err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		file.Close()
		return fmt.Errorf("failed to start cpu profile: %w", err)
	}

	c.file = file
	return nil
}

func (c *CPUProfiler) Stop() CPUResult {
	pprof.StopCPUProfile()
	c.file.Close()

	return CPUResult{
		Op:          c.op,
		ProfilePath: c.file.Name(),
	}
}
