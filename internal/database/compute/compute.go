package compute

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"
	"strings"

	"go.uber.org/zap"
)

const (
	minDataLen      = 2
	enterSymbolsLen = 2
)

type Compute struct {
	logger *zap.Logger
}

func NewCompute(logger *zap.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errors.New("could not make compute without logger")
	}
	return &Compute{logger: logger}, nil
}

func (c *Compute) Parse(data string) (request.Request, error) {

	splittedData := strings.Split(data, " ")

	if len(splittedData) < 2 {
		c.logger.Error("could not to parse less than two arguments")
		return request.Request{RequestType: commands.IncorrectCommand}, errors.New("could not to parse less than two arguments")
	}

	stringCommand := splittedData[0]

	parsedCommand, err := c.parseCommand(stringCommand)

	if err != nil {
		return request.Request{RequestType: parsedCommand}, err
	}

	c.logger.Debug("parsing of command is correct")

	parsedArgs, err := c.parseArguments(parsedCommand, splittedData[1:])

	if err != nil {
		c.logger.Error("parsing arguments has error")
		return request.Request{RequestType: parsedCommand, Args: parsedArgs}, err
	}

	c.logger.Debug("parsing of arguments is correct")

	return request.Request{RequestType: parsedCommand, Args: parsedArgs}, nil
}

func (c *Compute) parseCommand(stringCommand string) (int, error) {

	c.logger.Debug("started parse command")

	var parsedCommand int
	var err error

	switch strings.ToUpper(stringCommand) {

	case "GET":

		parsedCommand = commands.GetCommand

		c.logger.Debug("command parsed as get")

	case "DEL":

		parsedCommand = commands.DelCommand

		c.logger.Debug("command parsed as del")

	case "SET":

		parsedCommand = commands.SetCommand

		c.logger.Debug("command parsed as set")

	default:

		parsedCommand = commands.IncorrectCommand
		err = errors.New("incorrect command")

		c.logger.Error("incorrect command")
	}

	return parsedCommand, err
}

func (c *Compute) parseArguments(command int, arguments []string) ([]string, error) {

	c.logger.Debug("started to parse args")

	var parsedArgs []string

	if command == commands.SetCommand {

		if len(arguments) < minDataLen {
			return nil, errors.New("set command has two arguments")
		}

		parsedArgs = []string{arguments[0], arguments[1]}

	} else {
		parsedArgs = []string{arguments[0]}
	}

	lastArg := arguments[len(parsedArgs)-1]

	if len(lastArg) >= enterSymbolsLen && string(lastArg[len(lastArg)-enterSymbolsLen:]) == "\r\n" {
		parsedArgs[len(arguments)-1] = lastArg[:len(lastArg)-enterSymbolsLen]

		if parsedArgs[len(arguments)-1] == "" {
			return nil, errors.New("set command has two arguments")
		}
	}

	return parsedArgs, nil
}
