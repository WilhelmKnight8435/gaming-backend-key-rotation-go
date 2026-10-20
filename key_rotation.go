package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

func (c *InfraiClient) do(method, path string, body any) (json.RawMessage, error) {
	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBytes = b
	}
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	for attempt := 0; attempt < 4; attempt++ {
		var payload io.Reader
		if bodyBytes != nil {
			payload = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequest(method, c.BaseURL+path, payload)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		var env envelope
		if err := json.Unmarshal(data, &env); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 100 * time.Millisecond
			if retry := res.Header.Get("Retry-After"); retry != "" {
				if seconds, err := strconv.Atoi(retry); err == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return nil, fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode >= 500 {
			return nil, fmt.Errorf("infrai transport status %d", res.StatusCode)
		}
		return env.Data, nil
	}
	return nil, fmt.Errorf("retry budget exhausted")
}

type createdKey struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

func (c *InfraiClient) CreateKey(projectID, name string) (createdKey, error) {
	data, err := c.do(http.MethodPost, "/v1/account/keys/create", map[string]any{
		"project_id": projectID, "name": name, "scopes": []string{"logs.search"}, "idempotency_key": "gaming-rotation-demo",
	})
	if err != nil {
		return createdKey{}, err
	}
	var result createdKey
	if err := json.Unmarshal(data, &result); err != nil {
		return createdKey{}, err
	}
	return result, nil
}

func (c *InfraiClient) RotateKey(id string, graceHours int) error {
	_, err := c.do(http.MethodPost, "/v1/account/keys/rotate/"+id, map[string]any{"grace_hours": graceHours, "idempotency_key": "gaming-rotation-demo"})
	return err
}

func (c *InfraiClient) SearchLogs() (json.RawMessage, error) {
	return c.do(http.MethodGet, "/v1/logs/search", nil)
}

func rotationPlan(graceHours int) string {
	if graceHours > 0 {
		return "overlap"
	}
	return "blocked"
}
