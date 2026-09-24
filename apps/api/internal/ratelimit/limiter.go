package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	now      func() time.Time
	attempts map[string][]time.Time
}

func New(limit int, window time.Duration) *Limiter {
	if limit < 1 {
		limit = 1
	}
	if window < time.Second {
		window = time.Second
	}
	return &Limiter{
		limit:    limit,
		window:   window,
		now:      time.Now,
		attempts: make(map[string][]time.Time),
	}
}

func (l *Limiter) SetNow(now func() time.Time) {
	if now == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.now = now
}

func (l *Limiter) pruneLocked(key string, now time.Time) {
	cut := now.Add(-l.window)
	arr := l.attempts[key]
	i := 0
	for i < len(arr) && arr[i].Before(cut) {
		i++
	}
	if i > 0 {
		arr = arr[i:]
	}
	if len(arr) == 0 {
		delete(l.attempts, key)
		return
	}
	l.attempts[key] = arr
}

func (l *Limiter) Allow(key string) bool {
	if key == "" {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.pruneLocked(key, now)
	return len(l.attempts[key]) < l.limit
}

func (l *Limiter) Fail(key string) {
	if key == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.pruneLocked(key, now)
	l.attempts[key] = append(l.attempts[key], now)
}

func (l *Limiter) Reset(key string) {
	if key == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}
