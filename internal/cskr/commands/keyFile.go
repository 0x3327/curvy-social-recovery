package commands

import (
	"encoding/json"
	"os"
)

type KeyFile struct {
	SpendingKey string `json:"k"`
	ViewingKey  string `json:"v"`
}

func (k *KeyFile) LoadFromFile(path string) error {
	keyFileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(keyFileBytes, k)
	if err != nil {
		return err
	}

	return nil
}

func (k *KeyFile) WriteFile(path string) error {
	data, err := json.MarshalIndent(k, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0666)
}
