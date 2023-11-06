package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr        string
	BaseAddress string
}

func Configure() Config {
	conf := Config{}

	flag.StringVar(&conf.Addr, "a", ":8080", "server address")
	flag.StringVar(&conf.BaseAddress, "b", "http://localhost:8080", "base address for short url")

	flag.Parse()

	if addr, set := os.LookupEnv("SERVER_ADDRESS"); set {
		conf.Addr = addr
	}

	if baseURL, set := os.LookupEnv("BASE_URL"); set {
		conf.BaseAddress = baseURL
	}

	return conf
}
