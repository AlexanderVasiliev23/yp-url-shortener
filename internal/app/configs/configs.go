package configs

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultTokenLen              = 10
	defaultStorageFilePath       = "/tmp/short-url-db.json"
	defaultFileStorageBufferSize = 10
	defaultJWTSecretKey          = "V}^7/Y-t;F2*E,G>Tw<$Dd"
	defaultDebug                 = false
)

// Config missing godoc.
type Config struct {
	Addr                  string
	BaseAddress           string
	StorageFilePath       string
	DatabaseDSN           string
	JWTSecretKey          string
	TokenLen              int
	FileStorageBufferSize int
	DeleteWorkerConfig    DeleteWorkerConfig
	Debug                 bool
	EnableHTTPS           bool
}

// DeleteWorkerConfig missing godoc.
type DeleteWorkerConfig struct {
	RepoTimeout time.Duration
}

// MustConfigure missing godoc.
func MustConfigure() *Config {
	conf := &Config{}

	flag.StringVar(&conf.Addr, "a", ":8080", "server address")
	flag.StringVar(&conf.BaseAddress, "b", "http://localhost:8080", "base address for short url")
	flag.IntVar(&conf.TokenLen, "token-len", defaultTokenLen, "length of a token")
	flag.StringVar(&conf.StorageFilePath, "f", defaultStorageFilePath, "storage file path")
	flag.IntVar(&conf.FileStorageBufferSize, "file-storage-buffer-size", defaultFileStorageBufferSize, "size of file storage buffer")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "db data source name")
	flag.BoolVar(&conf.EnableHTTPS, "s", false, "enable HTTPS")

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

	if DatabaseDSN, set := os.LookupEnv("DATABASE_DSN"); set {
		conf.DatabaseDSN = DatabaseDSN
	}

	conf.JWTSecretKey = defaultJWTSecretKey
	if JWTSecretKey, set := os.LookupEnv("JWT_SECRET_KEY"); set {
		conf.JWTSecretKey = JWTSecretKey
	}

	conf.Debug = defaultDebug
	if debug, set := os.LookupEnv("DEBUG"); set {
		asBool, err := strconv.ParseBool(debug)
		if err != nil {
			panic(fmt.Errorf("parsing debug env as bool: %w", err))
		}
		conf.Debug = asBool
	}

	if enableHTTPS, set := os.LookupEnv("ENABLE_HTTPS"); set {
		asBool, err := strconv.ParseBool(enableHTTPS)
		if err != nil {
			panic(fmt.Errorf("parsing enableHTTPS env as bool: %w", err))
		}
		conf.EnableHTTPS = asBool
	}

	conf.DeleteWorkerConfig = DeleteWorkerConfig{
		RepoTimeout: 30 * time.Second,
	}

	return conf
}
