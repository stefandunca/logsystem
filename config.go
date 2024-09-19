package logsystem

import (
	"encoding/json"
	"os"
)

type Config struct {
	Drivers map[DriverID]json.RawMessage `json:"drivers"`
}

func LoadConfigFromFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	return loadConfig(data)
}

func loadConfig(data []byte) (Config, error) {
	config := Config{}
	err := json.Unmarshal(data, &config)
	return config, err
}
