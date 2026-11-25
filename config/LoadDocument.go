package config

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func LoadDocument() (string, error) {

	documentLink := os.Getenv("TNC_DOCUMENT_LINK")
	if documentLink == "" {
		log.Println("TNC_DOCUMENT_LINK is not set")
		return "", errors.New("TNC_DOCUMENT_LINK not set")
	}

	resp, err := http.Get(documentLink)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
