// Package serverconfigs - функции общего назначения, парсинг флагов и ключей шифрования
package serverconfigs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"flag"
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
	flag.Parse()
	var envcfg storage.ServerENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
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

func GetPrivateKey(pth string) ([]byte, error) {
	data, err := os.ReadFile(pth)
	return data, err

}

func DecypherData(key []byte, data []byte) ([]byte, error) {
	var x509Key *rsa.PrivateKey
	x509Key, _ = x509.ParsePKCS1PrivateKey(key)
	dencryptedData, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		x509Key,
		data,
		nil,
	)
	return dencryptedData, err
}
