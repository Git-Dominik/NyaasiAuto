package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Progress struct {
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

type Target struct {
	Title      string   `json:"title"`
	Uploaders  []string `json:"uploaders"`
	LatestSeen Progress `json:"latest_seen"`
}

func loadTargets(configFilePath string) ([]Target, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var targets []Target

	if err := json.Unmarshal(file, &targets); err != nil {
		fmt.Println(err)
	}

	return targets, nil
}

func (t *Target) saveProgress(configFilePath string, targets []Target) error {
	updateJson, err := json.MarshalIndent(targets, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	return os.WriteFile(configFilePath, updateJson, 0644)

}
