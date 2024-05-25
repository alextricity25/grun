package main
import (
	"os"
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	MaxLength int `json:"max_length" validate:"required,number"`
}

func Load(path string) (Config, error) {
	var cfg Config

	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

