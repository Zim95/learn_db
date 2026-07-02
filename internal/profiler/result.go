package profiler

/*
	The Result contract.

	Every profiler (Timer, Memory, CPU, ...) produces a value that satisfies
	Result. That single interface is the polymorphic glue of the package:

	  - the Collector can store one []Result regardless of what produced them,
	  - formatters/exporters type-switch on the concrete type to render/send it.

	The concrete result structs live next to the profiler that produces them
	(TimerResult in timer.go, MemoryResult in memory.go, CPUResult in cpu.go),
	so everything about one kind of measurement stays in one file. Adding a new
	profiler later means adding its struct in its own file and satisfying this
	interface — nothing here has to change.
*/

type Result interface {
	// Operation returns the name of the database operation that was
	// measured, e.g. "set" or "get".
	Operation() string
}
