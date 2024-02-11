package repeatedjson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type PlayerData struct { // repeated here + package json, bring to one place
	Name      string `json:"name"`
	HighScore int    `json:"high_score"`
}

func ParseRepeatedJSONFile(path string) ([]PlayerData, error) {
	var pdset []PlayerData

	f, err := os.Open(path)
	if err != nil {
		return pdset, err
	}

	sc := bufio.NewScanner(f)
	for {
		b, err := ScanNextLine(sc)
		if err != nil {
			if err != io.EOF {
				return pdset, err
			}
			break // if EOF reached then break
		}

		if IsLineAComment(b) {
			continue
		}

		var pd PlayerData
		err = json.Unmarshal(b, &pd)
		if err != nil {
			return pdset, err
		}

		pdset = append(pdset, pd)
	}

	return pdset, nil
}

func ScanNextLine(sc *bufio.Scanner) ([]byte, error) {
	if nextLineExists := sc.Scan(); !nextLineExists {
		err := sc.Err()
		if err != nil {
			return nil, err
		}
		return nil, io.EOF // if next line doesn't exist and no error, then eof reached
	}

	return sc.Bytes(), nil
}

func IsLineAComment(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return string(b[0]) == "#"
}

func DecodePlayerData(b []byte) (PlayerData, error) {
	var pd PlayerData
	err := json.Unmarshal(b, &pd)
	if err != nil {
		return PlayerData{}, fmt.Errorf("error un-marshalling player data: %w", err)
	}
	return pd, nil
}
