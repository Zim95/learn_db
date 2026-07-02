# learn_db

Learning to build a database — a minimal append-only, log-structured key/value
store, with an optional in-memory offset index and a small built-in profiler.

## Build & run

Quickest, no build step:

```bash
go run ./cmd/learn_db <flags> <command> [args]
```

Or build a binary:

```bash
go build -o learn_db ./cmd/learn_db
./learn_db <flags> <command> [args]

# via the Makefile (runs fmt + vet first); EXECUTABLE names the output file:
make build EXECUTABLE=learn_db
```

## Commands

| Command | Usage | Description |
|---------|-------|-------------|
| `set`   | `set <key> <value>` | Append a key/value record to the log |
| `get`   | `get <key>`         | Return the latest value for a key |

```bash
go run ./cmd/learn_db set name namah
go run ./cmd/learn_db get name
# Result: namah
```

Data is stored in an append-only log at `db.log` in the working directory. A
`get` returns the **most recent** value written for a key.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-db` | `log` | Storage engine (`log` is the only one implemented) |
| `-index` | *(none)* | Index implementation. `offset` builds an in-memory key→byte-offset map on startup so reads seek directly instead of scanning the whole file |
| `-timer` | `false` | Measure wall-clock execution time |
| `-mem` | `false` | Measure memory allocation |
| `-cpu` | `false` | Write a CPU profile (see `-cpu-profile`) |
| `-cpu-profile` | `cpu.prof` | Path for the CPU profile file |
| `-profile-format` | `terminal` | Profiling output format: `terminal` or `json` |

Using the offset index:

```bash
go run ./cmd/learn_db -index offset get name
```

## Profiling

Profiling is off by default. Enable any combination of `-timer`, `-mem`, and
`-cpu`. **Results print to stderr**, so they stay separate from the command's
result on stdout (you can redirect one without the other).

### Time + memory

```bash
go run ./cmd/learn_db -timer -mem set name namah
```

```
Result: OK
── profiling ─────────────────────────────
command  time   620µs
command  mem    alloc=2232B objects=13 heap-delta=2232B
```

- `alloc` — bytes allocated during the operation (cumulative counter delta; always ≥ 0)
- `objects` — number of heap objects allocated
- `heap-delta` — change in *live* heap bytes; may be negative if the garbage collector ran mid-operation

### JSON output

For scripts/dashboards, use `-profile-format json`:

```bash
go run ./cmd/learn_db -timer -mem -profile-format json get name
```

```json
[
  { "operation": "command", "type": "timer",  "duration": "62µs", "duration_ns": 62375 },
  { "operation": "command", "type": "memory", "total_alloc_bytes": 4768, "mallocs": 28, "heap_alloc_delta_bytes": 4768 }
]
```

### CPU profiling

```bash
go run ./cmd/learn_db -cpu -cpu-profile cpu.prof get name
go tool pprof cpu.prof        # then: top, list, web, ...
```

> ⚠️ **Caveat:** pprof samples the call stack ~100 times/second. A single
> `set`/`get` finishes in microseconds, so the profile will almost always show
> **0 samples**. CPU profiling only produces signal when it wraps something
> long-running (building the index over a large log, or a `go test -bench`
> loop). The plumbing is correct; the workload is just too short.

## Package layout

```
cmd/learn_db/          CLI entry point + command dispatch
internal/
  engine/              Engine interface (SetKey/GetKey/BuildIndex)
    logengine/         append-only log implementation
  index/               Index interface
    offset/            in-memory key→offset map
  record/              RecordPointer (offset into the log)
  profiler/            measure / store / model / format / export
    result.go            Result interface (shared contract)
    timer.go             Timer + TimerResult
    memory.go            MemoryProfiler + MemoryResult
    cpu.go               CPUProfiler + CPUResult
    collector.go         Collector — stores []Result
    profiler.go          orchestrator (Start/Stop, reads CLI flags)
    formatter/           terminal + json output
    exporter/            prometheus + otel (future; no-op stubs)
```

Each profiler *measures only*, the collector *stores only*, formatters
*display only* — every result type satisfies the `Result` interface, so new
profilers can be added without changing existing components.
