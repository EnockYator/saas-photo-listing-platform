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

// shardCount creates 64 independent maps and locks for low memory and good concurrency
const shardCount = 64

// bucket represents one client's token bucket.
type bucket struct {
	mu         sync.Mutex // protects bucket data
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
	trustProxy bool // if true, trust X-Forwarded-For and X-Real-IP headers for client IP headers
}

// RateLimiterOption configures a RateLimiter at construction time.
type RateLimiterOption func(*RateLimiter)

// WithTrustProxy enables trusting X-Forwarded-For and X-Real-IP headers for client IP.
func WithTrustProxy(trust bool) RateLimiterOption {
	return func(rl *RateLimiter) {
		rl.trustProxy = trust
	}
}

// NewRateLimiter creates a new rate limiter with the specified requests per second, burst size, and rate limiting options.
func NewRateLimiter(rps float64, burst float64, opts ...RateLimiterOption) *RateLimiter {
	rl := &RateLimiter{
		rps:   rps,
		burst: float64(burst),
	}

	for i := 0; i < shardCount; i++ {
		rl.shards[i] = &shard{
			clients: make(map[string]*bucket),
		}
	}

	for _, opt := range opts {
		opt(rl)
	}

	return rl
}

// getShard selects a shard based on hashed client key.
func (rl *RateLimiter) getShard(key string) *shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return rl.shards[h.Sum32()%shardCount]
}

// clientIP extracts the client IP address from the request, considering proxy headers if trustProxy is enabled.
//
// X-Forwarded-For / X-Real-IP are only trusted when rl.trustProxy is true,
// i.e. when this service sits behind a proxy/load balancer that the operator
// controls and that overwrites (rather than appends to) these headers. If
// trustProxy is false, these headers are ignored entirely, since any client
// can set them to an arbitrary value and spoof its rate-limit identity.
//
// When trusted, the rightmost address in X-Forwarded-For is used, since that
// is the hop closest to (and set by) the trusted proxy; the leftmost entries
// are client-supplied and can still be spoofed even through a legitimate
// proxy chain that appends rather than overwrites.
func (rl *RateLimiter) clientIP(r *http.Request) string {
	if rl.trustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			return strings.TrimSpace(parts[len(parts)-1])
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			return strings.TrimSpace(xri)
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}

// allow checks whether a request is allowed and returns a retry-after duration if not. It also updates the bucket state.
func (rl *RateLimiter) allow(key string) (bool, time.Duration) {
	now := time.Now()

	sh := rl.getShard(key)

	// Get or create the bucket (minimizing lock scope)
	sh.mu.Lock()

	b, ok := sh.clients[key]
	if !ok {
		b = &bucket{
			tokens:     rl.burst,
			lastRefill: now,
			lastSeen:   now,
		}
		sh.clients[key] = b
	}
	sh.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// Refill tokens based on elapsed time
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * rl.rps
	
	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}

	b.lastRefill = now
	b.lastSeen = now

	if b.tokens < 1 {
		deficit := 1 - b.tokens
		retryAfter := time.Duration(deficit/rl.rps*float64(time.Second))
		return false, retryAfter
	}

	b.tokens-- // Consume a token
	return true, 0
}

// StartCleanup periodically removes inactive clients from the rate limiter to free memory. It runs in a separate goroutine.
//
// Called once after constructing the RateLimiter, e.g.:
//
//		rl := NewRateLimiter(10, 20)
//		rl.StartCleanup(5 * time.Minute)
func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()

			for _, sh := range rl.shards {
				sh.mu.Lock()
				for key, b := range sh.clients {
					b.mu.Lock()
					stale := now.Sub(b.lastSeen) > 2*interval

					b.mu.Unlock()

					if stale {
						delete(sh.clients, key)
					}
				}
				sh.mu.Unlock()
			}
		}
	}()
}

// RateLimitMiddleware returns the HTTP middleware for this limiter instance.
//
// Lifecycle: 
// 		identify client -> find bucket -> refill tokens ->
//		consume token or reject -> emit observability data -> continue
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.clientIP(r)

		allowed, retryAfter := rl.allow(key)

		// OpenTelemetry enrichment
		span := trace.SpanFromContext(r.Context())

		if !allowed {
			if span.SpanContext().IsValid() {
				span.SetAttributes(
					attribute.String("rate_limit.status", "exceeded"),
					attribute.String("rate_limit.client", key),
					attribute.Float64("rate_limit.retry_after", retryAfter.Seconds()),
				)
			}

			if retryAfter > 0 {
				w.Header().Set("Retry-After", retryAfter.Truncate(time.Second).String())
			}

			http.Error(
				w,
				"rate limit exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}