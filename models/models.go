package models

import (
	"encoding/json"
	"log"
	"os"
)

type Models struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name          string  `json:"name"`
	ModelString   string  `json:"modelString"`
	ParamB        int     `json:"paramB"`
	ContextLength int     `json:"contextLength"`
	Cost          float32 `json:"cost"`
}

func LoadModels() Models {
	jsonModels, err := os.Open("./models/models.json")
	if err != nil {
		log.Fatal(err)
	}
	var models Models
	err = json.NewDecoder(jsonModels).Decode(&models)

	defer jsonModels.Close()

	return models
}
