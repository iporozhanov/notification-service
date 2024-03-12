package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	ips         map[string]*limiter
	timeout     time.Duration
	mutex       sync.Mutex
	maxRequests int64
}

type limiter struct {
	lastAccess time.Time
	count      int64
}

// NewRateLimiter creates a new instance of the RateLimiter.
func NewRateLimiter(maxRequests int64, timeout time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:         make(map[string]*limiter),
		timeout:     timeout,
		maxRequests: maxRequests,
	}
}

// Allow checks if the given key is allowed to make a request.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	l, ok := rl.ips[key]
	if !ok {
		l = &limiter{}
		rl.ips[key] = l
	}

	now := time.Now()
	if now.Sub(l.lastAccess) > rl.timeout {
		l.count = 0
	}

	l.lastAccess = now
	l.count++

	return l.count <= rl.maxRequests
}

// ClearExpired removes expired entries from the RateLimiter to prevent memory leaks.
func (rl *RateLimiter) ClearExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		for ip, l := range rl.ips {
			if time.Since(l.lastAccess) > rl.timeout {
				delete(rl.ips, ip)
			}
		}
		rl.mutex.Unlock()
	}
}
