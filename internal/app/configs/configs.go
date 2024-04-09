package configs

import (
	"encoding/json"
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
	GRPCServerAddr        string
	BaseAddress           string
	StorageFilePath       string
	DatabaseDSN           string
	JWTSecretKey          string
	TokenLen              int
	FileStorageBufferSize int
	DeleteWorkerConfig    DeleteWorkerConfig
	Debug                 bool
	EnableHTTPS           bool
	TrustedSubnet         string
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
	flag.StringVar(&conf.TrustedSubnet, "t", "", "trusted subnet")
	flag.StringVar(&conf.GRPCServerAddr, "grpc-server-addr", "", "GRPC server address")

	var configFilePath string
	flag.StringVar(&configFilePath, "c", "", "path to config JSON file")

	flag.Parse()

	if path, set := os.LookupEnv("CONFIG"); set {
		configFilePath = path
	}

	if addr, set := os.LookupEnv("SERVER_ADDRESS"); set {
		conf.Addr = addr
	}

	if baseURL, set := os.LookupEnv("BASE_URL"); set {
		conf.BaseAddress = baseURL
	}

	if storageFilePath, set := os.LookupEnv("FILE_STORAGE_PATH"); set {
		conf.StorageFilePath = storageFilePath
	}

	if databaseDSN, set := os.LookupEnv("DATABASE_DSN"); set {
		conf.DatabaseDSN = databaseDSN
	}

	if trustedSubnet, set := os.LookupEnv("TRUSTED_SUBNET"); set {
		conf.TrustedSubnet = trustedSubnet
	}

	if grpcServerAddress, set := os.LookupEnv("GRPC_SERVER_ADDRESS"); set {
		conf.GRPCServerAddr = grpcServerAddress
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

	if configFilePath != "" {
		fileContent, err := os.ReadFile(configFilePath)
		if err != nil {
			panic(fmt.Errorf("read config file %s: %w", configFilePath, err))
		}

		configFromFile := struct {
			ServerAddress         string `json:"server_address"`
			GRPCServerAddress     string `json:"grpc_server_address"`
			BaseURL               string `json:"base_url"`
			FileStoragePath       string `json:"file_storage_path"`
			DatabaseDSN           string `json:"database_dsn"`
			EnableHTTPS           bool   `json:"enable_https"`
			TokenLen              int    `json:"token_len"`
			FileStorageBufferSize int    `json:"file_storage_buffer_size"`
			TrustedSubnet         string `json:"trusted_subnet"`
		}{}

		if err := json.Unmarshal(fileContent, &configFromFile); err != nil {
			panic(fmt.Errorf("unmarshal config file %s: %w", configFilePath, err))
		}

		if configFromFile.ServerAddress != "" {
			conf.Addr = configFromFile.ServerAddress
		}

		if configFromFile.BaseURL != "" {
			conf.BaseAddress = configFromFile.BaseURL
		}

		if configFromFile.FileStoragePath != "" {
			conf.StorageFilePath = configFromFile.FileStoragePath
		}

		if configFromFile.DatabaseDSN != "" {
			conf.DatabaseDSN = configFromFile.DatabaseDSN
		}

		if configFromFile.TokenLen > 0 {
			conf.TokenLen = configFromFile.TokenLen
		}

		if configFromFile.FileStorageBufferSize > 0 {
			conf.FileStorageBufferSize = configFromFile.FileStorageBufferSize
		}

		if configFromFile.TrustedSubnet != "" {
			conf.TrustedSubnet = configFromFile.TrustedSubnet
		}

		if configFromFile.GRPCServerAddress != "" {
			conf.GRPCServerAddr = configFromFile.GRPCServerAddress
		}

		conf.EnableHTTPS = configFromFile.EnableHTTPS
	}

	return conf
}
