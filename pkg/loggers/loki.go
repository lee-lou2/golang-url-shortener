package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

type LokiLoggerWriter struct {
	App string
	Env string
}

func (clw *LokiLoggerWriter) Write(p []byte) (n int, err error) {
	logMessage := string(p)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	payload := LokiPayload{
		Streams: []LokiStream{
			{
				Stream: map[string]string{
					"app": clw.App,
					"env": clw.Env,
				},
				Values: [][]string{
					{timestamp, logMessage},
				},
			},
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("error marshalling JSON: %v", err)
	}
	url := "http://loki:3100/loki/api/v1/push"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending log to Loki: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log to Loki, status code: %d", resp.StatusCode)
	}
	return len(p), nil
}
