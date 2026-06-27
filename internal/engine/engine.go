package engine

type Engine interface {
	SetKey(key string, value string) error
	GetKey(key string) (string, error)
	BuildIndex() error
}
