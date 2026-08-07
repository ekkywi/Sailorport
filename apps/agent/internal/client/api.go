package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Worker struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}

type APIClient struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *APIClient) Register(name, hostname string, labels map[string]any) (Worker, error) {
	body, _ := json.Marshal(map[string]any{
		"name":     name,
		"hostname": hostname,
		"labels":   labels,
	})
	res, err := c.http.Post(c.baseURL+"/api/v1/workers/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return Worker{}, err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return Worker{}, fmt.Errorf("Register %d: %s", res.StatusCode, data)
	}
	var w Worker
	if err := json.Unmarshal(data, &w); err != nil {
		return Worker{}, err
	}
	return w, nil
}

func (c *APIClient) Heartbeat(id, status string) (Worker, error) {
	body, _ := json.Marshal(map[string]string{"status": status})
	url := fmt.Sprintf("%s/api/v1/workers/%s/heartbeat", c.baseURL, id)
	res, err := c.http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return Worker{}, err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return Worker{}, fmt.Errorf("Heartbeat %d: %s", res.StatusCode, data)
	}
	var w Worker
	if err := json.Unmarshal(data, &w); err != nil {
		return Worker{}, err
	}
	return w, nil
}
