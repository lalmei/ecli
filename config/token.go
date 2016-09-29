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

func StoreSession(url, token string) error {
	f, err := os.Create(absTokenFile())
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	b := &struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	}{url, token}
	return enc.Encode(b)
}

// Loads token and url endpoint information.
func LoadSession() (string, string, error) {
	f, err := os.Open(absTokenFile())
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("You are not logged in. Please run log in first.")
		}
		return "", "", err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	b := &struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	}{}
	if err := dec.Decode(b); err != nil {
		return "", "", err
	}
	return b.Url, b.Token, nil
}

func DeleteTokenFile() error {
	return os.Remove(absTokenFile())
}
