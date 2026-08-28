package model

type CatalogEnv struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

type CatalogEnvPublic map[string]any
