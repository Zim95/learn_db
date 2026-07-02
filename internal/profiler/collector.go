package profiler

/*
	Collector is the central storage for profiling data.

	It does not measure anything and it does not print anything. It only holds
	the Results produced by the various profilers so that a formatter or
	exporter can consume them later.
*/

type Collector struct {
	results []Result
}

func CreateCollector() *Collector {
	return &Collector{
		results: make([]Result, 0),
	}
}

// Add stores a single result.
func (c *Collector) Add(result Result) {
	c.results = append(c.results, result)
}

// Results returns all collected results in the order they were added.
func (c *Collector) Results() []Result {
	return c.results
}

// Clear drops all collected results.
func (c *Collector) Clear() {
	c.results = make([]Result, 0)
}
