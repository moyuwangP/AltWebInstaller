package util

import (
	"AltWebServer/app/model"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"os"
)

var config model.Config

func init() {
	var (
		jsonFile []byte
		err      error
	)

	if jsonFile, err = os.ReadFile("./config.json"); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(jsonFile, &config); err != nil {
		panic(err)
	}

	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(config); err != nil {
		panic(err)
	}
}

func Config() model.Config {
	return config
}
