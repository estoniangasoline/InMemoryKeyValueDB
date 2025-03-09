package storage

type Storage interface {
	Request(requestType int, arg ...string) (string, error)
}
