package cli

import (
	"bufio"
	"fmt"
	"inmemorykvdb/internal/database"
	"inmemorykvdb/internal/database/compute"
	"inmemorykvdb/internal/database/storage"
	"inmemorykvdb/internal/database/storage/engine"
	"os"

	"go.uber.org/zap"
)

type Database interface {
	HandleRequest(data string) (string, error)
}

func Run() {

	in := bufio.NewReader(os.Stdin)

	db, err := BuildDatabase()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {

		req, err := ReadRequest(in)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		resp, err := db.HandleRequest(req)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		WriteResponse(resp)
	}
}

func BuildDatabase() (Database, error) {
	engine, err := engine.NewInMemoryEngine(zap.NewNop())

	if err != nil {
		return nil, err
	}

	storage, err := storage.NewStorage(engine, zap.NewNop())

	if err != nil {
		return nil, err
	}

	compute, err := compute.NewCompute(zap.NewNop())

	if err != nil {
		return nil, err
	}

	db, err := database.NewInMemoryKvDb(compute, storage, zap.NewNop())

	if err != nil {
		return nil, err
	}

	return db, err
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
