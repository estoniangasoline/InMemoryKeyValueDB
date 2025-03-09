package database

type Database interface {
	Request(data string) (string, error)
}
