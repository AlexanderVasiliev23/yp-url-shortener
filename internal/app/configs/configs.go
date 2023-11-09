package configs

import (
	"flag"
	"os"
)

const defaultTokenLen = 10

type Config struct {
	Addr        string
	BaseAddress string
	TokenLen    int
}

func Configure() *Config {
	conf := &Config{}

	flag.StringVar(&conf.Addr, "a", ":8080", "server address")
	flag.StringVar(&conf.BaseAddress, "b", "http://localhost:8080", "base address for short url")
	flag.IntVar(&conf.TokenLen, "token-len", defaultTokenLen, "length of a token")

	flag.Parse()

	if addr, set := os.LookupEnv("SERVER_ADDRESS"); set {
		conf.Addr = addr
	}

	if baseURL, set := os.LookupEnv("BASE_URL"); set {
		conf.BaseAddress = baseURL
	}

	return conf
}
