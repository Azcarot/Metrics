package serverconfigs

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/caarlos0/env/v6"
)

func ParseFlagsAndENV() storage.Flags {
	var Flag storage.Flags
	flag.StringVar(&Flag.FlagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Flag.FlagFileStorage, "f", "/tmp/metrics-db.json", "address of a file-storage")
	flag.IntVar(&Flag.FlagStoreInterval, "i", 300, "interval for storing data")
	flag.BoolVar(&Flag.FlagRestore, "r", true, "reading data from file first")
	flag.Parse()
	var envcfg storage.ServerENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(envcfg.Address) > 0 {
		Flag.FlagAddr = envcfg.Address
	}

	if len(envcfg.FileStorage) > 0 {
		Flag.FlagFileStorage = envcfg.FileStorage
	}
	restore := os.Getenv("RESTORE")
	if len(restore) > 0 {
		envrestore, err := strconv.ParseBool(restore)
		if err != nil {
			Flag.FlagRestore = envrestore
		}
	}
	if len(envcfg.StoreInterval) == 0 {
		storeInterval, err := strconv.Atoi(envcfg.StoreInterval)
		if err == nil {
			Flag.FlagStoreInterval = storeInterval
		}
	}
	return Flag
}
