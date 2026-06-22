package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type DataBase struct {
	filePath string
}

const FILEPATH string = "db.log"

func CreateDatabase(path string) *DataBase {
	/**
	Constructor for DataBase struct
	:params:
		:path: string
	*/
	return &DataBase{
		filePath: path,
	}
}

func (db *DataBase) SetKey(key string, value string) error {
	/**
	Write to the DataBase file
	:params:
		:key: string: The key to write
		:value: string: The value to write
	*/
	file, err := os.OpenFile(
		db.filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return fmt.Errorf("Failed to open file %w", err)
	}

	defer file.Close() // Will close the file after this function exits

	record := fmt.Sprintf("%s, %s\n", key, value)
	_, err = file.WriteString(record) // := means declare and assign. Since err is already declared we cannot redeclare it. := is only for new values. We need to reassign. So we use = instead of :=

	if err != nil {
		return fmt.Errorf("Failed to write to file %w", err)
	}

	return nil
}

func (db *DataBase) GetKey(key string) (string, error) {
	/**
	Read from the DataBase file
	:params:
		:key: string: The key to read
	*/
	file, err := os.OpenFile(
		db.filePath,
		os.O_RDONLY,
		0644,
	)

	if err != nil {
		return "", fmt.Errorf("Error reading file %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	// as long as scanner has a value
	var value string
	for scanner.Scan() == true {
		// read the value
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		k := parts[0]
		v := parts[1]
		// update the latest value for the matching key
		if k == key {
			value = v
		}
	}
	return value, nil
}

type CommandHanlder func(
	args []string,
) (string, error)

func (db *DataBase) SetCommand(
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

func (db *DataBase) GetCommand(
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

func RunCommand(db *DataBase) (string, error) {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("./learn_db set <key> <value>")
		fmt.Println("./learn_db get <key>")
	}

	command := os.Args[1]

	handlers := map[string]CommandHanlder{
		"set": db.SetCommand,
		"get": db.GetCommand,
	}

	handler, exists := handlers[command]
	if !exists {
		return "", fmt.Errorf("Unknown Command %s", command)
	}

	return handler(os.Args[2:])
}

func main() {
	db := CreateDatabase("db.log")
	result, err := RunCommand(db)
	if err != nil {
		fmt.Printf("Failure\n: %v", err)
	}
	fmt.Printf("Result: %v", result)
}
