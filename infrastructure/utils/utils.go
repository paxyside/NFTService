package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
)

func GenerateUniqueHash() (string, error) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New("failed to generate unique hash: " + err.Error())
	}
	return hex.EncodeToString(b), nil
}

func GenerateInfuraURL(networkName, apiKey string) (string, error) {
	baseURL := fmt.Sprintf("https://%s.infura.io/v3/", networkName)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	parsedURL.Path += apiKey
	return parsedURL.String(), nil
}

func LoadABIFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read ABI file: %w", err)
	}
	return string(data), nil
}
