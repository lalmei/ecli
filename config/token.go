package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	tokenFile = ".token"
)

func absTokenFile() string {
	return path.Join(os.Getenv("HOME"), tokenFile)
}

func SaveToken(token string) error {
	f, err := os.Create(absTokenFile())
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	b := &struct {
		Token string `json:"token"`
	}{token}
	return enc.Encode(b)
}

func LoadToken() (string, error) {
	f, err := os.Open(absTokenFile())
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("You are not logged in. Please run log in first.")
		}
		return "", err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	b := &struct {
		Token string `json:"token"`
	}{}
	if err := dec.Decode(b); err != nil {
		return "", err
	}
	return b.Token, nil
}

func DeleteTokenFile() error {
	return os.Remove(absTokenFile())
}
