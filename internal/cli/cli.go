package cli

import (
	"bufio"
	"fmt"
	"inmemorykvdb/internal/compute"
	"inmemorykvdb/internal/database"
	"inmemorykvdb/internal/engine"
	"inmemorykvdb/internal/storage"
	"os"

	"go.uber.org/zap"
)

func Run() {
	in := bufio.NewReader(os.Stdin)

	var engineSize uint

	fmt.Print("Enter your db size: ")
	fmt.Fscanln(os.Stdin, &engineSize)

	engine, err := engine.NewInMemoryEngine(zap.NewNop(), engineSize)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	storage, err := storage.NewSimpleStorage(engine, zap.NewNop())

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	compute, err := compute.NewSimpleCompute(zap.NewNop())

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db, err := database.NewInMemoryKvDb(compute, storage, zap.NewNop())

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var req string
	var resp string

	for {

		fmt.Print("ENTER COMMAND: ")
		req, err = in.ReadString('\n')

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		resp, err = db.Request(req)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if resp == "" {
			fmt.Println("OK")
		} else {
			fmt.Println("RESPONSE IS:", resp)
		}
	}
}
