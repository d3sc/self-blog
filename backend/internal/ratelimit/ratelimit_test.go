package ratelimit

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	l := New(2, time.Minute)
	if !l.Allow("x") || !l.Allow("x") || l.Allow("x") {
		t.Fatal("expected third request to be limited")
	}
}
