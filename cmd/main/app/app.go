package main

import (
	"bufio"
	"fmt"
	"inmemorykvdb/internal/cli"
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

		req, err := cli.ReadRequest(in)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		resp, err := client.Send([]byte(req))

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		cli.WriteResponse(string(resp))
	}
}
