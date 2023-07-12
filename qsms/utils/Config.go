package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	SimpleMessageFee   int      `json:"simple_message_fee"`
	PeriodicMessageFee int      `json:"periodic_message_fee"`
	TemplateFee        int      `json:"template_fee"`
	BadWords           []string `json:"bad_words"`
}

func LoadConfig() *Config {
	file, err := os.Open("./config.json")
	if err != nil {
		return nil
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil
	}
	return &config
}

func SaveConfig(config *Config) error {
	file, err := os.Create("./config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err = encoder.Encode(config); err != nil {
		return err
	}
	return nil
}

func AddBadWord(config *Config, words []string) {
	config.BadWords = append(config.BadWords, words...)
}
