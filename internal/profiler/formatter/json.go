package formatter

import (
	"encoding/json"
	"io"

	"learn_db/internal/profiler"
)

/*
	JSON formatting.

	Emits the collected results as a JSON array, one object per result, tagged
	with a "type" field so a consumer can tell timer/memory/cpu records apart.
	Useful for scripts, dashboards, and diffing runs.

	We build plain maps rather than marshalling the Result values directly,
	because a heterogeneous []Result would otherwise lose the type distinction
	and expose Go field names instead of stable snake_case keys.
*/

func JSON(w io.Writer, results []profiler.Result) error {
	out := make([]map[string]any, 0, len(results))

	for _, result := range results {
		record := map[string]any{
			"operation": result.Operation(),
		}

		switch r := result.(type) {

		case profiler.TimerResult:
			record["type"] = "timer"
			record["duration_ns"] = r.Duration.Nanoseconds()
			record["duration"] = r.Duration.String()

		case profiler.MemoryResult:
			record["type"] = "memory"
			record["total_alloc_bytes"] = r.TotalAlloc
			record["mallocs"] = r.Mallocs
			record["heap_alloc_delta_bytes"] = r.HeapAllocDelta

		case profiler.CPUResult:
			record["type"] = "cpu"
			record["profile_path"] = r.ProfilePath

		default:
			record["type"] = "unknown"
		}

		out = append(out, record)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}
