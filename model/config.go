package model

import "fmt"

type Config struct {
	naziv     string
	verzija   float64
	parametri map[string]string
}

func NewConfig(naziv string, verzija float64, parametri map[string]string) *Config {
	return &Config{
		naziv:     naziv,
		verzija:   verzija,
		parametri: parametri,
	}
}

func (c Config) GenerateKey() string {
	return c.naziv + fmt.Sprintf("%.2f", c.verzija)
}

type ConfigRepository interface {
	Get(key string) Config
	Add(c Config)
	Delete(key string)
}
