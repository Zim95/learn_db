package main

import (
	"fmt"
	"learn_db/internal/engine/logengine"
	"os"
)

type CommandHanlder func(
	db *logengine.DataBase,
	args []string,
) (string, error)

func SetCommand(
	db *logengine.DataBase,
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
	db *logengine.DataBase,
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

func RunCommand(db *logengine.DataBase) (string, error) {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("./learn_db set <key> <value>")
		fmt.Println("./learn_db get <key>")
	}

	command := os.Args[1]

	handlers := map[string]CommandHanlder{
		"set": SetCommand,
		"get": GetCommand,
	}

	handler, exists := handlers[command]
	if !exists {
		return "", fmt.Errorf("Unknown Command %s", command)
	}

	return handler(db, os.Args[2:])
}

func main() {
	db := logengine.CreateDatabase(logengine.FILEPATH)
	result, err := RunCommand(db)
	if err != nil {
		fmt.Printf("Failure\n: %v", err)
	}
	fmt.Printf("Result: %v", result)
}
