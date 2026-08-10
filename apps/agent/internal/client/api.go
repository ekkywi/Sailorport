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

type DeploymentJob struct {
	ID            string  `json:"id"`
	ServiceID     string  `json:"service_id"`
	WorkerID      *string `json:"worker_id"`
	Status        string  `json:"status"`
	ServiceName   string  `json:"service_name"`
	WorkspacePath string  `json:"workspace_path"`
}

type UpdateDeploymentRequest struct {
	Status       string `json:"status"`
	ImageTag     string `json:"image_tag"`
	ContainerID  string `json:"container_id"`
	Port         *int   `json:"port"`
	ErrorMessage string `json:"error_message"`
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

func (c *APIClient) ClaimNext(workerID string) (*DeploymentJob, error) {
	body, _ := json.Marshal(map[string]string{"worker_id": workerID})
	res, err := c.http.Post(c.baseURL+"/api/v1/agent/jobs/next", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	data, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("Claim %d: %s", res.StatusCode, data)
	}
	var job DeploymentJob
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (c *APIClient) UpdateDeployment(id string, req UpdateDeploymentRequest) error {
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequest(http.MethodPatch, c.baseURL+"/api/v1/agent/deployments/"+id, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return fmt.Errorf("Update deployment %d: %s", res.StatusCode, data)
	}
	return nil
}
