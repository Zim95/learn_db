package main

import (
	"flag"
	"fmt"
	"learn_db/internal/engine"
	"learn_db/internal/engine/logengine"
	"learn_db/internal/index"
	"learn_db/internal/index/offset"
	"learn_db/internal/profiler"
	"learn_db/internal/profiler/formatter"
	"os"
)

type CommandHandler func(
	db engine.Engine,
	args []string,
) (string, error)

func SetCommand(
	db engine.Engine,
	args []string,
) (string, error) {
	/*
		Wrapper around the setkey function to match the function signature.
	*/
	if len(args) != 2 {
		return "", fmt.Errorf("usage: set <key> <value>")
	}

	err := db.SetKey(
		args[0],
		args[1],
	)

	if err != nil {
		return "", err
	}

	return "OK", nil
}

func GetCommand(
	db engine.Engine,
	args []string,
) (string, error) {
	/*
		Wrapper around the getkey function to match the function signature.
	*/
	if len(args) != 1 {
		return "", fmt.Errorf("usage: get <key>")
	}

	return db.GetKey(args[0])
}

func RunCommand(
	db engine.Engine,
	args []string,
) (string, error) {

	if len(args) == 0 {
		return "", fmt.Errorf("usage:\nset <key> <value>\nget <key>")
	}

	command := args[0]

	handlers := map[string]CommandHandler{
		"set": SetCommand,
		"get": GetCommand,
	}

	handler, exists := handlers[command]
	if !exists {
		return "", fmt.Errorf("unknown command %s", command)
	}

	return handler(db, args[1:])
}

func main() {
	dbType := flag.String("db", "log", "Database engine (log, lsm, memory)")
	indexType := flag.String("index", "", "Index implementation (offset, btree, art)")
	timer := flag.Bool("timer", false, "Measure execution time")
	mem := flag.Bool("mem", false, "Measure memory allocation")
	cpu := flag.Bool("cpu", false, "Write a CPU profile (see -cpu-profile)")
	cpuPath := flag.String("cpu-profile", "cpu.prof", "Path for the CPU profile file")
	profileFormat := flag.String("profile-format", "terminal", "Profiling output format (terminal, json)")

	flag.Parse()

	prof := profiler.New(profiler.Config{
		Timer:   *timer,
		Memory:  *mem,
		CPU:     *cpu,
		CPUPath: *cpuPath,
	})

	var idx index.Index
	switch *indexType {
	case "":
		idx = nil
	case "offset":
		idx = offset.CreateOffsetIndex()
	default:
		fmt.Printf("Unknown index: %s\n", *indexType)
		return
	}

	var db engine.Engine
	switch *dbType {
	case "log":
		db = logengine.CreateDatabase(
			logengine.FILEPATH,
			idx,
		)
	default:
		fmt.Printf("Unknown database engine: %s\n", *dbType)
		return
	}

	if idx != nil {
		db.BuildIndex()
	}

	if err := prof.Start("command"); err != nil {
		fmt.Printf("Failed to start profiler: %v\n", err)
		return
	}
	result, err := RunCommand(db, flag.Args())
	prof.Stop("command")

	if err != nil {
		fmt.Printf("Failure\n: %v", err)
	}
	fmt.Printf("Result: %v", result)

	if prof.Enabled() {
		// Profiling output goes to stderr so it stays separate from the
		// command's real result on stdout.
		switch *profileFormat {
		case "json":
			formatter.JSON(os.Stderr, prof.Results())
		default:
			formatter.Terminal(os.Stderr, prof.Results())
		}
	}
}
