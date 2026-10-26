package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	Max    int
	Window time.Duration
}

func New(max int, window time.Duration) *Limiter {
	return &Limiter{hits: map[string][]time.Time{}, Max: max, Window: window}
}
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	a := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if now.Sub(t) < l.Window {
			a = append(a, t)
		}
	}
	if len(a) >= l.Max {
		l.hits[key] = a
		return false
	}
	l.hits[key] = append(a, now)
	return true
}
func Key(r *http.Request) string { return r.RemoteAddr }
