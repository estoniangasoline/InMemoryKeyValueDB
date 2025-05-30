package main

import (
	"bufio"
	"fmt"
	"inmemorykvdb/internal/cli"
	"inmemorykvdb/internal/network"
	"os"
)

func main() {

	clientCnfg := cli.ParseClientOptions()

	client, err := network.NewClient(clientCnfg.Address,
		network.WithClientTimeout(clientCnfg.Timeout),
		network.WithClientMaxBufferSize(clientCnfg.MaxMessageSize))

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
