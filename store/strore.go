package store

import (
	"encoding/json"
	"math/rand"
	"os"
)

func LoadItems[T Event | Task](filename string) ([]T, error) {
	var items []T

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return items, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		items = []T{}
	} else {
		err = json.Unmarshal(data, &items)
		if err != nil {
			return nil, err
		}

	}
	return items, nil
}

func SaveItems[T Event | Task](filename string, items []T) error {
	jsonData, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonData, 0600)
	if err != nil {
		return err
	}
	return nil
}

func IdGen() string {
	lenght := 5
	b := make([]byte, lenght)
	for i := range b {
		randInd := rand.Intn(len(Charset))
		b[i] = Charset[randInd]
	}
	return string(b)
}
