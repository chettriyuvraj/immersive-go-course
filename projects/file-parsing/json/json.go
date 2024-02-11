package json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type PlayerData struct {
	Name      string `json:"name"`
	HighScore int    `json:"high_score"`
}

func DecodePlayerDataSetJSONFile(path string) ([]PlayerData, error) {
	f, err := os.Open(path)
	if err != nil {
		return []PlayerData{}, fmt.Errorf("error opening file: %w", err)
	}

	return DecodePlayerDataSet(f)
}

func DecodePlayerDataSet(r io.Reader) ([]PlayerData, error) {
	var pdset []PlayerData

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&pdset)
	if err != nil {
		return []PlayerData{}, fmt.Errorf("error decoding json data: %w", err)
	}

	return pdset, nil
}
