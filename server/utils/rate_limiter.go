package utils

import (
	"fmt"
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
	reqPerMin int
	lock      *sync.Mutex
	clients   map[string]*client
	useIp     bool
}

// Adapted from https://blog.logrocket.com/rate-limiting-go-application/#token-bucket-algorithm
func NewRateLimiter(
	reqPerMin int, // requests per minute
	clientTtl time.Duration,
	clearInterval time.Duration,
	useIp bool,
) RateLimiter {
	rateLimiter := RateLimiter{
		reqPerMin: reqPerMin,
		lock:      &sync.Mutex{},
		clients:   make(map[string]*client),
		useIp:     useIp,
	}
	go rateLimiter.SetClearInterval(clientTtl, clearInterval)
	return rateLimiter
}

func (l *RateLimiter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var rateKey string
		if l.useIp {
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				rateKey = ip
			}
		} else {
			rateKey = r.Header.Get("X-Client-Token")
			if !CheckValidToken(rateKey) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		fmt.Printf("-> request from %s\n", rateKey)

		l.lock.Lock()
		defer l.lock.Unlock()

		if _, found := l.clients[rateKey]; !found {
			limit := rate.Limit(float64(l.reqPerMin) / 60)
			l.clients[rateKey] = &client{limiter: rate.NewLimiter(limit, l.reqPerMin)}
		}
		l.clients[rateKey].lastSeen = time.Now()

		if !l.clients[rateKey].limiter.Allow() {
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

			for rateKey, client := range l.clients {
				if time.Since(client.lastSeen) > ttl {
					delete(l.clients, rateKey)
				}
			}
		}()
	}
}
