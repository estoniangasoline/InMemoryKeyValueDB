package cli

import (
	"bufio"
	"flag"
	"fmt"
	"inmemorykvdb/internal/config"
	"time"
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

func ParseServerConfig() *config.Config {
	def := flag.Bool("d", false, "No options load")

	engineType := flag.String("et", "in_memory", "Type of engine")

	address := flag.String("na", "127.0.0.1:3223", "Address of server")
	maxConns := flag.Int("nmc", 0, "Network max connections")
	maxMsgSize := flag.String("nms", "2MB", "Network max message size")
	idleTimeout := flag.Int("nt", 0, "Network timeout")
	isSync := flag.Bool("ns", false, "Synchronise the server")

	loggingLevel := flag.String("ll", "info", "Level of logging")
	output := flag.String("lo", "C:/go/InMemoryKeyValueDB/test/log/pretty.log", "Output of logging")

	batchSize := flag.Int("wbs", 5, "Batch size")
	batchTimeout := flag.Int("wbt", 1, "Batch timeout")
	maxSegmentSize := flag.String("wms", "B", "Batch max segment size")
	dataDir := flag.String("wd", "C:/go/InMemoryKeyValueDB/test/wal/", "Dir for WAL")
	fileName := flag.String("wfn", "wal.log", "Name for segment")

	replicaType := flag.String("rt", "master", "Type of replica")
	masterAddres := flag.String("rma", "127.0.0.2:3223", "Master address")
	syncInterval := flag.Int("ri", 1, "Replica interval")

	flag.Parse()

	if *def {
		return &config.Config{}
	}

	return &config.Config{
		Engine: &config.EngineConfig{
			EngineType: *engineType,
		},

		Network: &config.NetworkConfig{
			Address:        *address,
			MaxConnections: *maxConns,
			MaxMessageSize: *maxMsgSize,
			IdleTimeout:    time.Duration(*idleTimeout),
			IsSync:         *isSync,
		},

		Logging: &config.LoggingConfig{
			Level:  *loggingLevel,
			Output: *output,
		},

		WalConfig: &config.WalConfig{
			BatchSize:      *batchSize,
			BatchTimeout:   time.Duration(*batchTimeout),
			MaxSegmentSize: *maxSegmentSize,
			DataDirectory:  *dataDir,
			FileName:       *fileName,
		},

		Replication: &config.ReplicaConfig{
			ReplicaType:   *replicaType,
			MasterAddress: *masterAddres,
			SyncInterval:  time.Duration(*syncInterval),
		},
	}
}

func ParseClientOptions() *config.ClientConfig {
	address := flag.String("a", "127.0.0.1:3223", "Address to connect")
	maxMessageSize := flag.Int("m", 1000, "Max message size in bytes")
	timeOut := flag.Int("t", 0, "Timeout for connection in seconds")

	flag.Parse()

	return &config.ClientConfig{
		Address:        *address,
		MaxMessageSize: *maxMessageSize,
		Timeout:        time.Duration(*timeOut),
	}
}
