package index

import (
	"learn_db/internal/record"
)

type Index interface {
	Set(key string, pointer record.RecordPointer)
	Get(key string) (record.RecordPointer, error)
	Delete(key string)
}
