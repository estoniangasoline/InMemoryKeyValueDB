package parsing

import (
	"errors"
	"strconv"
)

const (
	byteMultiply     = 1
	kilobyteMultiply = 1024
	megabyteMultiply = 1048576
)

func ParseSize(unparsedSize string) (int, error) {
	unparsedMultiply := unparsedSize[len(unparsedSize)-2:]
	multiply, err := ParseMultiply(unparsedMultiply)

	if err != nil {
		return -1, err
	}

	var endOfCount int

	if multiply == byteMultiply {
		endOfCount = 1
	} else {
		endOfCount = 2
	}

	count, err := strconv.Atoi(unparsedSize[:len(unparsedSize)-endOfCount])

	if err != nil {
		return -1, errors.New("incorrect integer in string")
	}

	size := count * multiply

	return size, nil
}

func ParseMultiply(unparsedMultiply string) (int, error) {
	if len(unparsedMultiply) == 0 {
		return -1, errors.New("empty string")
	}

	if '0' <= unparsedMultiply[0] && unparsedMultiply[0] <= '9' && len(unparsedMultiply) > 1 {
		unparsedMultiply = unparsedMultiply[1:]
	}

	switch unparsedMultiply {
	case "KB":
		return kilobyteMultiply, nil
	case "MB":
		return megabyteMultiply, nil
	case "B":
		return byteMultiply, nil
	}

	return -1, errors.New("unknown or forbidden buffer size")
}
