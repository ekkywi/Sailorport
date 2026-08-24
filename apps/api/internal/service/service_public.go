package service

import (
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func PublicService(svc model.Service) model.Service {
	svc.WebhookSecretSet = strings.TrimSpace(svc.WebhookSecret) != ""
	svc.WebhookSecret = ""
	return svc
}

func PublicServices(services []model.Service) []model.Service {
	out := make([]model.Service, len(services))
	for i := range services {
		out[i] = PublicService(services[i])
	}
	return out
}
