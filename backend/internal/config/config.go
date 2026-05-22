package config

import "github.com/spf13/viper"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	OpenAIURL   string
	OpenAIKey   string
	OpenAIModel string
}

func Load() Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("OPENAI_BASE_URL", "https://api.openai.com/v1")
	viper.SetDefault("OPENAI_MODEL", "openai/gpt-4o")
	return Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		OpenAIURL:   viper.GetString("OPENAI_BASE_URL"),
		OpenAIKey:   viper.GetString("OPENAI_API_KEY"),
		OpenAIModel: viper.GetString("OPENAI_MODEL"),
	}
}
