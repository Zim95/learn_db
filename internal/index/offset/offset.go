package offset

import (
	"fmt"
	"learn_db/internal/record"
)

type OffsetIndex struct {
	offsetIndex map[string]record.RecordPointer
}

func CreateOffsetIndex() *OffsetIndex {
	return &OffsetIndex{
		offsetIndex: make(map[string]record.RecordPointer),
	}
}

func (oi *OffsetIndex) Set(
	key string,
	pointer record.RecordPointer,
) {
	oi.offsetIndex[key] = pointer
}

func (oi *OffsetIndex) Get(
	key string,
) (record.RecordPointer, error) {

	pointer, exists := oi.offsetIndex[key]
	if !exists {
		return record.RecordPointer{}, fmt.Errorf("key %s does not exist", key)
	}

	return pointer, nil
}

func (oi *OffsetIndex) Delete(key string) {
	delete(oi.offsetIndex, key)
}
