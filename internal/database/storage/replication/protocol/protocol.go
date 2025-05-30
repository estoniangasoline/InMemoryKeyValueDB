package protocol

import "encoding/json"

const (
	OkStatus      = 0
	ErrorStatus   = 1
	UnfoundStatus = 2

	ReadLast = 0
	ReadAll  = 1
)

type Request struct {
	Type         int
	LastFileName string `json:"last_file_name"`
}

type Response struct {
	Status    int      `json:"status"`
	FileNames []string `json:"file_name"`
	Data      [][]byte `json:"data"`
}

func newRequest(reqType int, lastFileName string) *Request {
	return &Request{Type: reqType, LastFileName: lastFileName}
}

func ReadAllRequest() *Request {
	return newRequest(ReadAll, "")
}

func ReadLastRequest(lastFileName string) *Request {
	return newRequest(ReadLast, lastFileName)
}

func OkResponseOneFile(fileName string, data []byte) *Response {
	return newResponse(OkStatus, []string{fileName}, [][]byte{data})
}

func OkResponseAllFiles(fileNames []string, data [][]byte) *Response {
	return newResponse(OkStatus, fileNames, data)
}

func newResponse(status int, fileName []string, data [][]byte) *Response {
	return &Response{Status: status, FileNames: fileName, Data: data}
}

func UnfoundResponse() *Response {
	return newResponse(UnfoundStatus, []string(nil), [][]byte(nil))
}

func ErrorResponse(err error) *Response {
	return newResponse(ErrorStatus, []string{"error"}, [][]byte{[]byte(err.Error())})
}

func Marshal[ProtocolObject Request | Response](po *ProtocolObject) ([]byte, error) {
	data, err := json.Marshal(po)
	return data, err
}

func Unmarshal[ProtocolObject Request | Response](po *ProtocolObject, data []byte) error {
	err := json.Unmarshal(data, po)

	return err
}
