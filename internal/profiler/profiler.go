package profiler

import "fmt"

/*
	Profiler is the orchestrator.

	It owns the Collector and coordinates whichever sub-profilers are enabled
	(via Config, which is populated from CLI flags). Callers talk to this type
	and never touch Timer/MemoryProfiler/CPUProfiler directly.

	Usage:

		p := profiler.New(profiler.Config{Timer: true, Memory: true})
		p.Start("set")
		defer p.Stop("set")   // stops enabled profilers, stores results

	Start/Stop can be called again for a second region (e.g. "index" then
	"set") and results accumulate in the Collector. The one exception is CPU
	profiling, which is process-global and can wrap only one region per run.
*/

type Config struct {
	Timer   bool
	Memory  bool
	CPU     bool
	CPUPath string
}

type Profiler struct {
	config    Config
	collector *Collector

	// currently active measurements (nil when their profiler is disabled)
	timer  *Timer
	memory *MemoryProfiler
	cpu    *CPUProfiler
}

func New(config Config) *Profiler {
	return &Profiler{
		config:    config,
		collector: CreateCollector(),
	}
}

// Enabled reports whether any profiler is turned on. Handy for deciding
// whether to print a profiling section at all.
func (p *Profiler) Enabled() bool {
	return p.config.Timer || p.config.Memory || p.config.CPU
}

// Start creates every enabled profiler for the named operation and begins
// measuring. Each profiler follows the same Create -> Start lifecycle; CPU is
// the only one whose Start can fail, so it's the only one we check for error.
func (p *Profiler) Start(op string) error {
	if p.config.Timer {
		p.timer = CreateTimer(op)
		p.timer.Start()
	}

	if p.config.Memory {
		p.memory = CreateMemoryProfiler(op)
		p.memory.Start()
	}

	if p.config.CPU {
		p.cpu = CreateCPUProfiler(op, p.config.CPUPath)
		if err := p.cpu.Start(); err != nil {
			return fmt.Errorf("failed to start cpu profiling: %w", err)
		}
	}

	return nil
}

// Stop ends every active measurement and pushes its result into the Collector.
//
// The timer is stopped first so that the (relatively expensive) memory and CPU
// teardown does not pollute the measured duration.
func (p *Profiler) Stop(op string) {
	if p.timer != nil {
		p.collector.Add(p.timer.Stop())
		p.timer = nil
	}

	if p.memory != nil {
		p.collector.Add(p.memory.Stop())
		p.memory = nil
	}

	if p.cpu != nil {
		p.collector.Add(p.cpu.Stop())
		p.cpu = nil
	}
}

// Results returns everything collected so far, for a formatter or exporter.
func (p *Profiler) Results() []Result {
	return p.collector.Results()
}
