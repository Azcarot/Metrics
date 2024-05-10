// Package serverconfigs - функции общего назначения, парсинг флагов и ключей шифрования
package serverconfigs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/caarlos0/env/v6"
)

var PrivateKey []byte

// ParseFlagsAndENV - Обрабатывает переменные окружения и флаги, приоритет у переменных окружения
func ParseFlagsAndENV() storage.Flags {
	var Flag storage.Flags
	flag.StringVar(&Flag.FlagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Flag.FlagDBAddr, "d", "", "address for db")
	flag.StringVar(&Flag.FlagFileStorage, "f", "/tmp/metrics-db.json", "address of a file-storage")
	flag.IntVar(&Flag.FlagStoreInterval, "i", 300, "interval for storing data")
	flag.BoolVar(&Flag.FlagRestore, "r", true, "reading data from file first")
	flag.StringVar(&Flag.FlagKey, "k", "", "Hash key")
	flag.StringVar(&Flag.FlagCrypto, "crypto-key", "", "path to private key")
	flag.StringVar(&Flag.FlagConfig, "config", "", "path to server config file")
	flag.StringVar(&Flag.FlagSubnet, "t", "", "CIDR for trusted subnet")
	flag.Parse()
	var envcfg storage.ServerENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}
	//Обработка флагов из файла и пересечения их с флагами командной строки и переменных окружения
	if Flag.FlagConfig != "" {
		if envcfg.ConfigPath != "" {
			Flag.FlagConfig = envcfg.ConfigPath
		}
		isFlagSet := make(map[string]bool)
		fileData, err := parseFile(Flag.FlagConfig)
		if err != nil {
			log.Fatal(err)
		}
		fileFlag := *fileData
		flag.Visit(func(f *flag.Flag) {
			isFlagSet[f.Name] = true
		})
		if len(fileFlag.FlagAddr) == 0 || isFlagSet["a"] {
			fileFlag.FlagAddr = Flag.FlagAddr
			if envcfg.Address != "" {
				fileFlag.FlagAddr = envcfg.Address
			}
		}
		if len(fileFlag.FlagDBAddr) == 0 || isFlagSet["d"] {
			fileFlag.FlagDBAddr = Flag.FlagDBAddr
			if envcfg.DBAddress != "" {
				fileFlag.FlagDBAddr = envcfg.DBAddress
			}
		}
		if len(fileFlag.FlagSubnet) == 0 || isFlagSet["t"] {
			fileFlag.FlagSubnet = Flag.FlagSubnet
			if envcfg.TrustedSubnet != "" {
				fileFlag.FlagSubnet = envcfg.TrustedSubnet
			}
		}
		if len(fileFlag.FlagFileStorage) == 0 || isFlagSet["f"] {
			fileFlag.FlagFileStorage = Flag.FlagFileStorage
			if envcfg.FileStorage != "" {
				fileFlag.FlagFileStorage = envcfg.FileStorage
			}
		}
		if fileFlag.FlagStoreInterval == 0 || isFlagSet["i"] {
			fileFlag.FlagStoreInterval = Flag.FlagStoreInterval
			if len(envcfg.StoreInterval) == 0 {
				storeInterval, err := strconv.Atoi(envcfg.StoreInterval)
				if err == nil {
					fileFlag.FlagStoreInterval = storeInterval
				}
			}
		}
		if !fileFlag.FlagRestore || isFlagSet["r"] {
			fileFlag.FlagRestore = Flag.FlagRestore
			restore := os.Getenv("RESTORE")
			if len(restore) > 0 {
				envrestore, err := strconv.ParseBool(restore)
				if err != nil {
					fileFlag.FlagRestore = envrestore
				}
			}

		}
		if len(fileFlag.FlagKey) == 0 || isFlagSet["k"] {
			fileFlag.FlagKey = Flag.FlagKey
			if envcfg.Key != "" {
				fileFlag.FlagKey = envcfg.Key
			}
		}

		if len(fileFlag.FlagCrypto) == 0 || isFlagSet["crypto-key"] {
			fileFlag.FlagCrypto = Flag.FlagCrypto
			if envcfg.CryptoKey != "" {
				fileFlag.FlagCrypto = envcfg.CryptoKey
			}

		}

		return fileFlag
	}
	if len(envcfg.Address) > 0 {
		Flag.FlagAddr = envcfg.Address
	}

	if len(envcfg.CryptoKey) > 0 {
		Flag.FlagCrypto = envcfg.CryptoKey
	}

	if len(envcfg.DBAddress) > 0 {
		Flag.FlagDBAddr = envcfg.DBAddress
	}

	if len(envcfg.FileStorage) > 0 {
		Flag.FlagFileStorage = envcfg.FileStorage
	}
	if len(envcfg.Key) > 0 {
		Flag.FlagKey = envcfg.Key
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

func parseFile(path string) (*storage.Flags, error) {
	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	var flagData storage.Flags
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON
	var envData storage.ServerENV
	err = json.Unmarshal(data, &envData)
	if err != nil {
		return nil, err
	}
	flagData.FlagStoreInterval, _ = strconv.Atoi(envData.StoreInterval)
	flagData.FlagDBAddr = envData.DBAddress
	flagData.FlagAddr = envData.Address
	flagData.FlagKey = envData.Key
	flagData.FlagCrypto = envData.CryptoKey
	flagData.FlagFileStorage = envData.FileStorage
	flagData.FlagRestore = envData.Restore
	return &flagData, nil
}

func GetPrivateKey(pth string) ([]byte, error) {
	data, err := os.ReadFile(pth)
	return data, err

}

func DecypherData(key []byte, data []byte) ([]byte, error) {
	var x509Key *rsa.PrivateKey
	x509Key, _ = x509.ParsePKCS1PrivateKey(key)
	msgLen := len(data)
	step := x509Key.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedData, err := rsa.DecryptOAEP(
			sha256.New(),
			rand.Reader,
			x509Key,
			data[start:finish],
			nil,
		)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedData...)
	}

	return decryptedBytes, nil
}
