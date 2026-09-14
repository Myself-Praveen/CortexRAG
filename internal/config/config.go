package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`

	Vector struct {
		Backend string `mapstructure:"backend"`
	} `mapstructure:"vector"`

	Embedding struct {
		Provider string `mapstructure:"provider"`
		Model    string `mapstructure:"model"`
	} `mapstructure:"embedding"`

	LLM struct {
		Provider string `mapstructure:"provider"`
		Model    string `mapstructure:"model"`
	} `mapstructure:"llm"`

	RAG struct {
		ChunkSize    int `mapstructure:"chunk_size"`
		ChunkOverlap int `mapstructure:"chunk_overlap"`
		CacheSize    int `mapstructure:"cache_size"`
		TopK         int `mapstructure:"top_k"`
	} `mapstructure:"rag"`

	Worker struct {
		Count int `mapstructure:"count"`
	} `mapstructure:"worker"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	v.SetDefault("port", 8080)
	v.SetDefault("vector.backend", "memory")
	v.SetDefault("embedding.provider", "gemini")
	v.SetDefault("embedding.model", "gemini-embedding-001")
	v.SetDefault("llm.provider", "gemini")
	v.SetDefault("llm.model", "gemini-3.5-flash-lite")
	v.SetDefault("rag.chunk_size", 512)
	v.SetDefault("rag.chunk_overlap", 64)
	v.SetDefault("rag.cache_size", 10000)
	v.SetDefault("rag.top_k", 5)
	v.SetDefault("worker.count", 0)

	// Set Env prefixes
	v.SetEnvPrefix("CORTEXRAG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			fmt.Printf("Warning: error reading config file: %v\n", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
