package handler

import (
	"net/http"
	"testing"
)

func TestClientIP_StripsPort(t *testing.T) {
	r := &http.Request{RemoteAddr: "127.0.0.1:54321"}
	if got := clientIP(r); got != "127.0.0.1" {
		t.Fatalf("got %q", got)
	}
}

func TestClientIP_IPv6(t *testing.T) {
	r := &http.Request{RemoteAddr: "[::1]:8080"}
	if got := clientIP(r);got != "::1" {
		t.Fatalf("got %q", got)
	}
}