package api

import (
	"os"
	"strconv"
)

type Config struct {
	Port   int
	Ollama OllamaConfig
}

type OllamaConfig struct {
	Domain string
}

func NewConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	return &Config{
		Port: port,
		Ollama: OllamaConfig{
			Domain: os.Getenv("OLLAMA_DOMAIN"),
		},
	}
}
