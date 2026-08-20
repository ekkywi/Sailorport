package service

import (
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func labelString(labels map[string]any, key string) string {
	if labels == nil {
		return ""
	}
	v, ok := labels[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return strings.TrimSpace(fmt.Sprint(v))
	}
	return strings.TrimSpace(s)
}

func workerEnvironments(w model.Worker) []string {
	raw := labelString(w.Labels, "environments")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, p := range parts {
		slug := strings.ToLower(strings.TrimSpace(p))
		if slug == "" {
			continue
		}
		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}
		out = append(out, slug)
	}
	return out
}

func workerAllowsEnvironment(w model.Worker, envSlug string) bool {
	envSlug = strings.ToLower(strings.TrimSpace(envSlug))
	if envSlug == "" {
		envSlug = "dev"
	}
	allowed := workerEnvironments(w)
	if len(allowed) == 0 {
		return true
	}
	for _, slug := range allowed {
		if slug == envSlug {
			return true
		}
	}
	return false
}

func workerEnvConflict(w model.Worker, envSlug string) error {
	if workerAllowsEnvironment(w, envSlug) {
		return nil
	}
	allowed := strings.Join(workerEnvironments(w), ", ")
	return fmt.Errorf(
		"%w: worker %q does not allow environment %q (allowed: %s)",
		ErrConflict, w.Name, envSlug, allowed,
	)
}
