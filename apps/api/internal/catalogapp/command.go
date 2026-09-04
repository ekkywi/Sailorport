package catalogapp

import (
	"fmt"
	"strings"
)

func ResolveCommand(command []string, env map[string]string) ([]string, error) {
	if len(command) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(command))
	for i, raw := range command {
		arg := strings.TrimSpace(raw)
		resolved, err := resolveArg(arg, env)
		if err != nil {
			return nil, fmt.Errorf("command[%d]: %w", i, err)
		}
		out = append(out, resolved)
	}
	return out, nil
}

func resolveArg(arg string, env map[string]string) (string, error) {
	var b strings.Builder
	rest := arg
	for {
		start := strings.Index(rest, "${")
		if start < 0 {
			b.WriteString(rest)
			return b.String(), nil
		}
		b.WriteString(rest[:start])
		endRel := strings.Index(rest[start:], "}")
		if endRel < 0 {
			return "", fmt.Errorf("unclosed placeholder in %q", arg)
		}
		end := start + endRel
		name := rest[start+2 : end]
		val, ok := env[name]
		if !ok {
			return "", fmt.Errorf("missing env value for ${%s}", name)
		}
		b.WriteString(val)
		rest = rest[end+1:]
	}
}
