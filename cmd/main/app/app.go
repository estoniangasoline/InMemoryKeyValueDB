package main

import (
	"bufio"
	"fmt"
	"inmemorykvdb/internal/network"
	"os"
	"time"
)

const (
	address        = ":8080"
	maxMessageSize = 4096
	timeOut        = 0 * time.Second
)

func main() {

	client, err := network.NewClient(address, network.WithClientTimeout(timeOut), network.WithClientMaxBufferSize(maxMessageSize))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer client.Close()

	in := bufio.NewReader(os.Stdin)

	for {

		req, err := ReadRequest(in)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		resp, err := client.Send([]byte(req))

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		WriteResponse(string(resp))
	}
}

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
