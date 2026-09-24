package ratelimit

import (
	"testing"
	"time"
)

func TestLimiter_BlocksAfter(t *testing.T) {
	lim := New(3, time.Minute)
	base := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	t0 := base
	lim.SetNow(func() time.Time { return t0 })

	key := "1,2,3,4"
	for i := 0; i < 3; i++ {
		if !lim.Allow(key) {
			t.Fatalf("attempt %d should be allowed", i+1)
		}
		lim.Fail(key)
		t0 = t0.Add(time.Second)
	}
	if lim.Allow(key) {
		t.Fatal("4th attempt should be blocked")
	}
}

func TestLimiter_WindowExpires(t *testing.T) {
	lim := New(2, 10*time.Second)
	base := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	t0 := base
	lim.SetNow(func() time.Time { return t0 })

	key := "ip"
	lim.Fail(key)
	lim.Fail(key)
	if lim.Allow(key) {
		t.Fatal("should be blocked inside window")
	}

	t0 = base.Add(11 * time.Second)
	if !lim.Allow(key) {
		t.Fatal("should allow after window expires")
	}
}

func TestLimiter_ResetClears(t *testing.T) {
	lim := New(1, time.Minute)
	base := time.Date(2026, 9, 25, 12, 0, 0, 0, time.UTC)
	lim.SetNow(func() time.Time { return base })

	key := "ip"
	lim.Fail(key)
	if lim.Allow(key) {
		t.Fatal("blocked after one fail when limit=1")
	}
	lim.Reset(key)
	if !lim.Allow(key) {
		t.Fatal("allow after reset")
	}
}

func TestLimiter_EmptyKeyAlwaysAllowed(t *testing.T) {
	lim := New(1, time.Minute)
	if !lim.Allow("") {
		t.Fatal("empty key should not block")
	}
	lim.Fail("")
	if !lim.Allow("") {
		t.Fatal("empty key still allowed")
	}
}
