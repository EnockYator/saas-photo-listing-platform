package middleware

import (
	"hash/fnv"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Create 64 independent maps and locks for low memory and good concurrency
const shardCount = 64

// bucket represents one client's token bucket.
type bucket struct {
	mu         sync.Mutex // protect bucket data
	tokens     float64 // current token count
	lastRefill time.Time
	lastSeen   time.Time // tracks last request - used for cleanup
}

// shard holds a subset of clients to reduce lock contention.
type shard struct {
	mu      sync.Mutex
	clients map[string]*bucket
}

// RateLimiter is a sharded token-bucket rate limiter.
type RateLimiter struct {
	rps   float64 // request/sec eg. 10rps = 10 tokens/sec
	burst float64 // maximum bucket capacity eg. 20 = Can instantly spend 20 requests
	shards [shardCount]*shard // array of 64 shard pointers
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		rps:   rps,
		burst: float64(burst),
	}

	// Create all 64 shards
	for i := 0; i < shardCount; i++ {
		rl.shards[i] = &shard{
			clients: make(map[string]*bucket),
		}
	}

	return rl
}

// getShard selects a shard based on hashed client key.
func (rl *RateLimiter) getShard(key string) *shard {
	h := fnv.New32a() // create hash function
	_, _ = h.Write([]byte(key)) // convert string eg. "192.168.1.5" to bytes
	return rl.shards[uint(h.Sum32())%shardCount]
}

// clientIP extracts a safe client IP from request headers or remote address.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}

// allow checks whether a request is allowed and returns retry-after duration.
func (rl *RateLimiter) allow(ip string) (bool, time.Duration) {
	now := time.Now()

	s := rl.getShard(ip)

	// Get or create bucket (minimize lock scope)
	s.mu.Lock()
	b, ok := s.clients[ip]
	if !ok {
		b = &bucket{
			tokens:     rl.burst,
			lastRefill: now,
			lastSeen:   now,
		}
		s.clients[ip] = b
	}
	s.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// refill tokens
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * rl.rps

	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}

	b.lastRefill = now
	b.lastSeen = now

	// reject if no tokens
	if b.tokens < 1 {
		deficit := 1 - b.tokens
		retryAfter := time.Duration(deficit/rl.rps) * time.Second
		return false, retryAfter
	}

	b.tokens--
	return true, 0
}

// StartCleanup periodically removes inactive clients to prevent memory leaks.
func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()

			for _, s := range rl.shards {
				s.mu.Lock()

				for ip, b := range s.clients {
					b.mu.Lock()
					if now.Sub(b.lastSeen) > 2*interval {
						delete(s.clients, ip)
					}
					b.mu.Unlock()
				}

				s.mu.Unlock()
			}
		}
	}()
}

// Middleware returns the HTTP middleware.
// Lifecycle of the middleware: 
// identify client → find bucket → refill tokens → consume token or reject → emit observability data → continue request processing
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := clientIP(r)

		allowed, retryAfter := rl.allow(ip)

		span := trace.SpanFromContext(r.Context())

		if !allowed {

			if span != nil {
				span.SetAttributes(
					attribute.String("rate_limit", "exceeded"),
					attribute.String("client_ip", ip),
				)
			}

			if retryAfter > 0 {
				w.Header().Set("Retry-After", retryAfter.String())
			}

			http.Error(
				w,
				"rate limit exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		if span != nil {
			span.SetAttributes(
				attribute.String("rate_limit", "allowed"),
				attribute.String("client_ip", ip),
			)
		}

		next.ServeHTTP(w, r)
	})
}