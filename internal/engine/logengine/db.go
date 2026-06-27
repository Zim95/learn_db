package logengine

import (
	"bufio"
	"fmt"
	"io"
	"learn_db/internal/index"
	"learn_db/internal/record"
	"os"
	"strings"
)

type DataBase struct {
	filePath string
	index    index.Index
}

const FILEPATH string = "db.log"

func CreateDatabase(path string, idx index.Index) *DataBase {
	/**
	Constructor for DataBase struct
	:params:
		:path: string
	*/
	return &DataBase{
		filePath: path,
		index:    idx,
	}
}

func (db *DataBase) BuildIndex() error {
	file, err := os.Open(db.filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var offset int64 = 0

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])

		db.index.Set(
			key,
			record.RecordPointer{
				Offset: offset,
			},
		)

		// Move our logical cursor forward
		offset += int64(len(scanner.Bytes()) + 1)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	return nil
}

func (db *DataBase) setKeyLinear(
	key string,
	value string,
) error {

	file, err := os.OpenFile(
		db.filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	record := fmt.Sprintf("%s, %s\n", key, value)

	_, err = file.WriteString(record)
	if err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

func (db *DataBase) setKeyIndexed(
	key string,
	value string,
) error {
	if db.index == nil {
		return fmt.Errorf("no index configured")
	}

	file, err := os.OpenFile(
		db.filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to get EOF offset: %w", err)
	}

	defer file.Close()

	line := fmt.Sprintf("%s, %s\n", key, value)

	_, err = file.WriteString(line)
	if err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	db.index.Set(
		key,
		record.RecordPointer{
			Offset: offset,
		},
	)

	return nil
}

func (db *DataBase) getKeyLinear(
	key string,
) (string, error) {

	file, err := os.Open(db.filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var value string
	var found bool

	for scanner.Scan() {

		line := scanner.Text()

		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			continue
		}

		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		if k == key {
			value = v
			found = true
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if !found {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	return value, nil
}

func (db *DataBase) getKeyIndexed(
	key string,
) (string, error) {
	if db.index == nil {
		return "", fmt.Errorf("no index configured")
	}

	file, err := os.Open(db.filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	pointer, err := db.index.Get(key)
	if err != nil {
		return "", err
	}

	_, err = file.Seek(pointer.Offset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek: %w", err)
	}

	reader := bufio.NewReader(file)

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read record: %w", err)
	}

	parts := strings.SplitN(line, ",", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid record")
	}

	return strings.TrimSpace(parts[1]), nil
}

func (db *DataBase) SetKey(
	key,
	value string,
) error {
	if db.index != nil {
		return db.setKeyIndexed(key, value)
	}
	return db.setKeyLinear(key, value)
}

func (db *DataBase) GetKey(
	key string,
) (string, error) {
	if db.index != nil {
		return db.getKeyIndexed(key)
	}
	return db.getKeyLinear(key)
}
