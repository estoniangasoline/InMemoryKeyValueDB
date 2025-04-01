package main

import (
	"bytes"
	"fmt"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/initialization"
	"os"
)

var configFile = os.Getenv("CONFIG_FILE_NAME")

func main() {

	cnfg := &config.Config{}

	if configFile != "" {
		file, err := os.ReadFile(configFile)

		if err != nil {
			fmt.Println("problems with reading config")
			return
		}

		reader := bytes.NewReader(file)

		cnfg, err = config.NewConfig(reader)

		if err != nil {
			cnfg = &config.Config{}
		}
	}

	initializer, err := initialization.NewInitializer(cnfg)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	initializer.StartDatabase()
}
