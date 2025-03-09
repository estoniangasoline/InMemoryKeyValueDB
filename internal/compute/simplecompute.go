package compute

import (
	"errors"
	"inmemorykvdb/internal/commands"
	"strings"

	"go.uber.org/zap"
)

type SimpleCompute struct {
	logger *zap.Logger
}

func NewSimpleCompute(logger *zap.Logger) (Compute, error) {
	if logger == nil {
		return nil, errors.New("could not make compute without logger")
	}
	return &SimpleCompute{logger: logger}, nil
}

func (c *SimpleCompute) Parse(data string) (int, []string, error) {

	splittedData := strings.Split(data, " ")

	if len(splittedData) < 2 {
		c.logger.Error("could not to parse less than two arguments")
		return commands.UncorrectCommand, []string{}, errors.New("could not to parse less than two arguments")
	}

	stringCommand := splittedData[0]

	var intCommand int
	var arguments []string

	switch strings.ToUpper(stringCommand) {

	case "GET":
		intCommand = commands.GetCommand
		arguments = []string{splittedData[1]}

		c.logger.Debug("command parsed as get")

	case "DEL":
		intCommand = commands.DelCommand
		arguments = []string{splittedData[1]}

		c.logger.Debug("command parsed as del")

	case "SET":

		if len(splittedData) < 3 {
			c.logger.Error("set command could not be realized with one argument")
			return commands.UncorrectCommand, []string{}, errors.New("set command could not be realized with one argument")
		}

		intCommand = commands.SetCommand
		arguments = []string{splittedData[1], splittedData[2]}

		c.logger.Debug("command parsed as set")

	default:
		c.logger.Error("unknown command")
		return commands.UncorrectCommand, []string{}, errors.New("unknown command")
	}

	lastArg := arguments[len(arguments)-1]

	if len(lastArg) >= 4 && string(lastArg[len(lastArg)-2:]) == "\r\n" {
		arguments[len(arguments)-1] = lastArg[:len(lastArg)-2]
	}

	return intCommand, arguments, nil
}
