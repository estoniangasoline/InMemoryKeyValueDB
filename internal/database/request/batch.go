package request

import (
	"errors"
	"fmt"
	"strings"
)

const (
	intSize      = 4
	delimElement = ' '
)

type Batch struct {
	Data     []*Request
	ByteSize int
	MaxSize  int
}

func NewBatch(maxSize int) *Batch {
	return &Batch{Data: make([]*Request, 0, maxSize), MaxSize: maxSize}
}

func (b *Batch) Add(req *Request) {
	b.ByteSize += intSize

	for _, arg := range req.Args {
		b.ByteSize += len(arg)
	}

	b.Data = append(b.Data, req)
}

func (b *Batch) LoadFilesToBatch(fileNames []string, files [][]byte) error {

	errorFiles := &strings.Builder{}

	for i, filedata := range files {
		err := b.LoadData(filedata)

		if err != nil {
			errorFiles.Write([]byte(fileNames[i]))
			errorFiles.WriteByte(delimElement)
			continue
		}
	}

	if errorFiles.String() != "" {
		strErr := errorFiles.String()
		strErr = strErr[:len(strErr)-1]
		return fmt.Errorf("has files with errors: %s", strErr)
	}

	return nil
}

func (b *Batch) LoadData(data []byte) error {
	var startIndex int

	var hasUnparsedRequests bool

	for i, elem := range data {
		if string([]byte{elem}) == EndElement {
			unparsed, err := NewRequest(string(data[startIndex:i]))
			startIndex = i + 1

			if err == nil {
				b.Data = append(b.Data, unparsed)
			} else if !hasUnparsedRequests {
				hasUnparsedRequests = true
			}
		}
	}

	if hasUnparsedRequests {
		return errors.New("has unparsed requests")
	}

	return nil
}

func (b *Batch) ParseBatch() ([]byte, error) {
	byteBuffer := make([]byte, 0, b.ByteSize)

	hasUnparsedRequests := false

	for _, req := range b.Data {
		requestInBytes, err := req.ParseToBytes()

		if err != nil {
			hasUnparsedRequests = true
		} else {
			byteBuffer = append(byteBuffer, requestInBytes...)
		}
	}

	if hasUnparsedRequests {
		return byteBuffer, errors.New("some requests has not parsed")
	}

	return byteBuffer, nil
}

func (b *Batch) IsFilled() bool {
	return b.ByteSize > b.MaxSize
}

func (b *Batch) Clear() {
	b.Data = make([]*Request, 0, b.MaxSize) // allocate more memory than necessary to avoid trouble
	b.ByteSize = 0
}
