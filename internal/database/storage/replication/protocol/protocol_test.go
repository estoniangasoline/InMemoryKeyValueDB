package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRequest(t *testing.T) {
	t.Parallel()

	lastFileName := "wal1.log"
	expectedReq := &Request{Type: ReadAll, LastFileName: lastFileName}

	actualReq := newRequest(ReadAll, lastFileName)

	assert.Equal(t, expectedReq, actualReq)
}

func Test_NewResponse(t *testing.T) {
	t.Parallel()

	status := OkStatus
	fileName := []string{"wal1.log"}
	data := [][]byte{[]byte("set biba boba")}

	expectedResp := &Response{Status: status, FileNames: fileName, Data: data}

	actualResp := newResponse(status, fileName, data)

	assert.Equal(t, expectedResp, actualResp)
}

func Test_Response(t *testing.T) {
	t.Parallel()

	status := OkStatus
	fileName := []string{"wal1.log"}
	data := [][]byte{[]byte("set biba boba")}

	resp := newResponse(status, fileName, data)

	marshaled, err := Marshal(resp)
	assert.Nil(t, err)

	unmarshaled := &Response{}
	err = Unmarshal(unmarshaled, marshaled)
	assert.Nil(t, err)

	assert.Equal(t, resp, unmarshaled)
}

func Test_Request(t *testing.T) {
	t.Parallel()

	lastFileName := "wal1.log"

	req := newRequest(ReadLast, lastFileName)

	marshaled, err := Marshal(req)
	assert.Nil(t, err)

	unmarshaled := &Request{}
	err = Unmarshal(unmarshaled, marshaled)
	assert.Nil(t, err)

	assert.Equal(t, req, unmarshaled)
}
