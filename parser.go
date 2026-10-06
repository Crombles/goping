package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func fetchSubscription(link string) ([]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var rawData string

	decodedBytes, err := decodeBase64(strings.TrimSpace(string(body)))
	if err != nil {
		rawData = string(body)
	} else {
		rawData = string(decodedBytes)
	}

	lines := strings.Split(strings.ReplaceAll(rawData, "\r\n", "\n"), "\n")
	return lines, nil
}

func decodeBase64(input string) ([]byte, error) {
	cleaned := strings.TrimSpace(input)

	if cleaned == "" {
		return nil, fmt.Errorf("empty base64 string")
	}

	decoders := []*base64.Encoding{
		base64.URLEncoding,
		base64.RawURLEncoding,
		base64.StdEncoding,
		base64.RawStdEncoding,
	}

	for _, dec := range decoders {
		if data, err := dec.DecodeString(cleaned); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("failed to decode base64: invalid format or corrupted data")
}
