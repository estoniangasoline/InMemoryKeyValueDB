package cli

import (
	"bufio"
	"fmt"
)

func ReadRequest(in *bufio.Reader) (string, error) {

	fmt.Print("ENTER COMMAND: ")
	req, err := in.ReadString('\n')

	if err != nil {
		return "", err
	}

	return req, nil
}

func WriteResponse(resp string) {
	if resp == "" {
		fmt.Println("OK")
	} else {
		fmt.Println("RESPONSE IS:", resp)
	}
}
