package config

import (
	"encoding/json"
	"fmt"
	"os"

	"impulse/internal/timeclock"
)

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Floors < 2 {
		return fmt.Errorf("Floors must be at least 2 (regular floors + boss)")
	}
	if c.Monsters < 1 {
		return fmt.Errorf("Monsters must be at least 1")
	}
	if c.Duration < 1 {
		return fmt.Errorf("Duration must be at least 1 hour")
	}
	if _, err := timeclock.Parse(c.OpenAt); err != nil {
		return fmt.Errorf("invalid OpenAt: %w", err)
	}
	return nil
}

func (c Config) RegularFloors() int {
	return c.Floors - 1
}
