package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Counter is a simple counter for tracking counts
type Counter struct {
	value atomic.Int64
}

// NewCounter creates a new counter
func NewCounter() *Counter {
	return &Counter{}
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	c.value.Add(1)
}

// Add adds the given value to the counter
func (c *Counter) Add(v int64) {
	c.value.Add(v)
}

// Get returns the current value of the counter
func (c *Counter) Get() int64 {
	return c.value.Load()
}

// Reset resets the counter to 0
func (c *Counter) Reset() {
	c.value.Store(0)
}

// Histogram tracks a distribution of values (typically for timing)
type Histogram struct {
	count   atomic.Int64
	sum     float64
	buckets []bucket
	mu      sync.RWMutex
}

type bucket struct {
	value float64
	count atomic.Int64
}

// NewHistogram creates a new histogram with default buckets
func NewHistogram(buckets []float64) *Histogram {
	h := &Histogram{
		buckets: make([]bucket, 0, len(buckets)+1),
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Initialize buckets
	for _, b := range buckets {
		h.buckets = append(h.buckets, bucket{
			value: b,
			count: atomic.Int64{},
		})
	}

	// Add +Inf bucket
	h.buckets = append(h.buckets, bucket{
		value: -1, // -1 represents +Inf
		count: atomic.Int64{},
	})

	return h
}

// Observe records a value
func (h *Histogram) Observe(value float64) {
	h.count.Add(1)

	h.mu.Lock()
	h.sum += value
	h.mu.Unlock()

	h.mu.RLock()
	defer h.mu.RUnlock()

	// Increment buckets
	for i := range h.buckets {
		if value <= h.buckets[i].value || h.buckets[i].value == -1 {
			h.buckets[i].count.Add(1)
		}
	}
}

// Count returns the total count of observations
func (h *Histogram) Count() int64 {
	return h.count.Load()
}

// Sum returns the sum of all observations
func (h *Histogram) Sum() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sum
}

// GetBuckets returns the bucket counts
func (h *Histogram) GetBuckets() map[string]int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]int64)
	for _, b := range h.buckets {
		key := "+Inf"
		if b.value != -1 {
			key = formatFloat64(b.value)
		}
		result[key] = b.count.Load()
	}
	return result
}

func formatFloat64(f float64) string {
	// Use %f to preserve decimal places
	return fmt.Sprintf("%f", f)
}

// DefaultBuckets returns default histogram buckets for timing (in seconds)
func DefaultBuckets() []float64 {
	return []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
}

// MetricsRegistry holds all metrics
type MetricsRegistry struct {
	mu      sync.RWMutex
	metrics map[string]any
}

// NewMetricsRegistry creates a new metrics registry
func NewMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{
		metrics: make(map[string]any),
	}
}

// Register registers a metric with a name
func (r *MetricsRegistry) Register(name string, metric any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[name] = metric
}

// GetCounter retrieves or creates a counter
func (r *MetricsRegistry) GetCounter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m, exists := r.metrics[name]; exists {
		if c, ok := m.(*Counter); ok {
			return c
		}
	}

	c := NewCounter()
	r.metrics[name] = c
	return c
}

// GetHistogram retrieves or creates a histogram
func (r *MetricsRegistry) GetHistogram(name string, buckets []float64) *Histogram {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m, exists := r.metrics[name]; exists {
		if h, ok := m.(*Histogram); ok {
			return h
		}
	}

	h := NewHistogram(buckets)
	r.metrics[name] = h
	return h
}

// Snapshot returns a snapshot of all metrics
func (r *MetricsRegistry) Snapshot() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := make(map[string]any)
	for k, v := range r.metrics {
		switch m := v.(type) {
		case *Counter:
			snapshot[k] = m.Get()
		case *Histogram:
			snapshot[k] = map[string]any{
				"count":   m.Count(),
				"sum":     m.Sum(),
				"buckets": m.GetBuckets(),
			}
		}
	}

	return snapshot
}

// RecordDuration measures and records the duration of a function
func RecordDuration(histogram *Histogram, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start).Seconds()
	histogram.Observe(duration)
}

// RecordDurationWithError measures and records the duration of a function that may return an error
func RecordDurationWithError(histogram *Histogram, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start).Seconds()
	histogram.Observe(duration)
	return err
}
