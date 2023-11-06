package config

import "flag"

type Config struct {
	Addr        string
	BaseAddress string
}

func Configure() Config {
	conf := Config{}

	flag.StringVar(&conf.Addr, "a", ":8080", "server address")
	flag.StringVar(&conf.BaseAddress, "b", ":8080", "base address for short url")

	flag.Parse()

	return conf
}
