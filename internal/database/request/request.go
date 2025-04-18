package request

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"strings"
)

const (
	DelimElement = " "
	EndElement   = "\n"

	minRequestLen        = 2
	minSetRequestArgsLen = 2

	commandIndex = 0
	argsIndex    = 1
)

type Request struct {
	RequestType int
	Args        []string
}

func (r *Request) ParseToBytes() ([]byte, error) {
	var command string

	switch r.RequestType {
	case commands.SetCommand:
		command = "SET"
	case commands.DelCommand:
		command = "DEL"
	default:
		return []byte(nil), errors.New("incorrect command type")
	}

	parsed := []byte(command)

	for _, arg := range r.Args {
		parsed = append(parsed, []byte(DelimElement+arg)...)
	}

	parsed = append(parsed, []byte(EndElement)...)

	return parsed, nil
}

func NewRequest(data string) (*Request, error) {
	splittedData := strings.Split(data, DelimElement)

	if len(splittedData) < minRequestLen {
		return nil, errors.New("incorrect data")
	}

	req := &Request{}
	req.Args = splittedData[argsIndex:]

	lastArg := req.Args[len(req.Args)-1]
	lastElement := lastArg[len(lastArg)-1]

	if string([]byte{lastElement}) == EndElement {
		lastArg = lastArg[:len(lastArg)-1]
	}

	req.Args[len(req.Args)-1] = lastArg

	switch splittedData[commandIndex] {
	case "SET":
		req.RequestType = commands.SetCommand
		if len(req.Args) < minSetRequestArgsLen {
			return nil, errors.New("set command in data has less than two arguments")
		}

	case "DEL":
		req.RequestType = commands.DelCommand

	default:
		return nil, errors.New("incorrect command")
	}

	return req, nil
}
