package utils

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	rate    int
	lock    *sync.Mutex
	clients map[string]*client
}

// Adapted from https://blog.logrocket.com/rate-limiting-go-application/#token-bucket-algorithm
func NewRateLimiter(
	rate int, // requests per minute
	clientTtl time.Duration,
	clearInterval time.Duration,
) RateLimiter {
	rateLimiter := RateLimiter{
		rate:    rate,
		lock:    &sync.Mutex{},
		clients: make(map[string]*client),
	}
	go rateLimiter.SetClearInterval(clientTtl, clearInterval)
	return rateLimiter
}

func (l *RateLimiter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		l.lock.Lock()
		defer l.lock.Unlock()

		if _, found := l.clients[ip]; !found {
			limit := rate.Limit(float64(l.rate) / 60)
			l.clients[ip] = &client{limiter: rate.NewLimiter(limit, l.rate)}
		}
		l.clients[ip].lastSeen = time.Now()

		if !l.clients[ip].limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (l *RateLimiter) SetClearInterval(ttl, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			l.lock.Lock()
			defer l.lock.Unlock()

			for ip, client := range l.clients {
				if time.Since(client.lastSeen) > ttl {
					delete(l.clients, ip)
				}
			}
		}()
	}
}
