package main

type Config struct {
	Port    int
	AofPath string
}

func buildConfig() Config {
	return Config{
		Port:    6379,
		AofPath: "dump.aof",
	}
}
