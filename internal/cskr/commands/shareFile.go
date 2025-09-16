package commands

import (
	"encoding/json"
	"os"

	keyrecovery "github.com/0x3327/curvy-social-recovery/key_recovery"
)

type ShareFile struct {
	Point        string `json:"x"`
	SpendingEval string `json:"spendingEval"`
	ViewingEval  string `json:"viewingEval"`
}

func (sf *ShareFile) FromShare(s keyrecovery.Share) {
	sf.Point = s.Point
	sf.ViewingEval = s.ViewingEval
	sf.SpendingEval = s.SpendingEval
}

func (sf *ShareFile) ReadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return json.Unmarshal(data, sf)
}

func (sf *ShareFile) WriteFile(path string) error {
	data, err := json.MarshalIndent(sf, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0666)
}
