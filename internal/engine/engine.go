package engine

type Engine interface {
	SET(key string, value string) error
	GET(key string) (string, error)
	DEL(key string) error
}
