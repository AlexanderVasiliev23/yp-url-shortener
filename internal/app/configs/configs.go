package configs

import (
	"flag"
	"os"
)

const (
	defaultTokenLen        = 10
	defaultStorageFilePath = "/tmp/short-url-db.json"
)

type Config struct {
	Addr            string
	BaseAddress     string
	TokenLen        int
	StorageFilePath string
}

func Configure() *Config {
	conf := &Config{}

	flag.StringVar(&conf.Addr, "a", ":8080", "server address")
	flag.StringVar(&conf.BaseAddress, "b", "http://localhost:8080", "base address for short url")
	flag.IntVar(&conf.TokenLen, "token-len", defaultTokenLen, "length of a token")
	flag.StringVar(&conf.StorageFilePath, "f", defaultStorageFilePath, "storage file path")

	flag.Parse()

	if addr, set := os.LookupEnv("SERVER_ADDRESS"); set {
		conf.Addr = addr
	}

	if baseURL, set := os.LookupEnv("BASE_URL"); set {
		conf.BaseAddress = baseURL
	}

	if storageFilePath, set := os.LookupEnv("FILE_STORAGE_PATH"); set {
		conf.StorageFilePath = storageFilePath
	}

	return conf
}
