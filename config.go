package main

type Config struct {
	Port int
}

func buildConfig() Config {
	return Config{
		Port: 6379,
	}
}
