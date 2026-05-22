package config

import "github.com/spf13/viper"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	OllamaURL   string
	LlamaModel  string
}

func Load() Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("OLLAMA_URL", "http://ollama:11434")
	viper.SetDefault("LLAMA_MODEL", "llama3")
	return Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		OllamaURL:   viper.GetString("OLLAMA_URL"),
		LlamaModel:  viper.GetString("LLAMA_MODEL"),
	}
}
