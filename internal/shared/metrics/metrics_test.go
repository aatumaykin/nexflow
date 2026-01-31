package metrics

import (
	"testing"
	"time"
)

func TestCounter_Inc(t *testing.T) {
	c := NewCounter()

	c.Inc()
	if c.Get() != 1 {
		t.Errorf("Expected 1, got %d", c.Get())
	}

	c.Inc()
	if c.Get() != 2 {
		t.Errorf("Expected 2, got %d", c.Get())
	}
}

func TestCounter_Add(t *testing.T) {
	c := NewCounter()

	c.Add(5)
	if c.Get() != 5 {
		t.Errorf("Expected 5, got %d", c.Get())
	}

	c.Add(10)
	if c.Get() != 15 {
		t.Errorf("Expected 15, got %d", c.Get())
	}
}

func TestCounter_Reset(t *testing.T) {
	c := NewCounter()
	c.Inc()
	c.Inc()

	c.Reset()

	if c.Get() != 0 {
		t.Errorf("Expected 0 after reset, got %d", c.Get())
	}
}

func TestHistogram_Observe(t *testing.T) {
	buckets := []float64{0.1, 0.5, 1.0}
	h := NewHistogram(buckets)

	h.Observe(0.05)
	h.Observe(0.3)
	h.Observe(0.7)
	h.Observe(1.5)

	if h.Count() != 4 {
		t.Errorf("Expected 4 observations, got %d", h.Count())
	}
}

func TestHistogram_Buckets(t *testing.T) {
	buckets := []float64{0.1, 0.5, 1.0}
	h := NewHistogram(buckets)

	h.Observe(0.05) // <= 0.1, so goes into all buckets
	h.Observe(0.3)  // <= 0.5, so goes into 0.5 and all below
	h.Observe(0.7)  // <= 1.0, so goes into 1.0 and all below
	h.Observe(1.5)  // > 1.0, so only goes into +Inf

	b := h.GetBuckets()

	t.Logf("Buckets: %+v", b)

	// Prometheus-style cumulative buckets:
	// - bucket 0.1: all observations <= 0.1 (only 0.05)
	if b["0.100000"] != 1 {
		t.Errorf("Expected bucket 0.100000 to have 1, got %d", b["0.100000"])
	}
	// - bucket 0.5: all observations <= 0.5 (0.05 and 0.3)
	if b["0.500000"] != 2 {
		t.Errorf("Expected bucket 0.500000 to have 2, got %d", b["0.500000"])
	}
	// - bucket 1.0: all observations <= 1.0 (0.05, 0.3, 0.7)
	if b["1.000000"] != 3 {
		t.Errorf("Expected bucket 1.000000 to have 3, got %d", b["1.000000"])
	}
	// - +Inf: all observations
	if b["+Inf"] != 4 {
		t.Errorf("Expected +Inf bucket to have 4, got %d", b["+Inf"])
	}
}

func TestHistogram_Sum(t *testing.T) {
	h := NewHistogram(DefaultBuckets())

	h.Observe(1.0)
	h.Observe(2.0)
	h.Observe(3.0)

	if h.Sum() != 6.0 {
		t.Errorf("Expected sum 6.0, got %v", h.Sum())
	}
}

func TestMetricsRegistry_GetCounter(t *testing.T) {
	r := NewMetricsRegistry()

	c1 := r.GetCounter("test_counter")
	c2 := r.GetCounter("test_counter")

	// Should return the same counter
	if c1 != c2 {
		t.Error("Expected same counter instance")
	}

	c1.Inc()
	if c2.Get() != 1 {
		t.Errorf("Expected both counters to have same value, got %d", c2.Get())
	}
}

func TestMetricsRegistry_GetHistogram(t *testing.T) {
	r := NewMetricsRegistry()

	buckets := []float64{0.1, 0.5}
	h1 := r.GetHistogram("test_histogram", buckets)
	h2 := r.GetHistogram("test_histogram", buckets)

	// Should return the same histogram
	if h1 != h2 {
		t.Error("Expected same histogram instance")
	}

	h1.Observe(0.2)
	if h2.Count() != 1 {
		t.Errorf("Expected both histograms to have same count, got %d", h2.Count())
	}
}

func TestMetricsRegistry_Snapshot(t *testing.T) {
	r := NewMetricsRegistry()

	c := r.GetCounter("test_counter")
	h := r.GetHistogram("test_histogram", DefaultBuckets())

	c.Inc()
	c.Inc()
	h.Observe(1.5)

	snapshot := r.Snapshot()

	if snapshot["test_counter"].(int64) != 2 {
		t.Errorf("Expected counter snapshot to be 2, got %v", snapshot["test_counter"])
	}

	histSnapshot := snapshot["test_histogram"].(map[string]any)
	if histSnapshot["count"].(int64) != 1 {
		t.Errorf("Expected histogram count to be 1, got %v", histSnapshot["count"])
	}
}

func TestRecordDuration(t *testing.T) {
	h := NewHistogram(DefaultBuckets())

	RecordDuration(h, func() {
		time.Sleep(10 * time.Millisecond)
	})

	if h.Count() != 1 {
		t.Errorf("Expected 1 observation, got %d", h.Count())
	}

	if h.Sum() == 0 {
		t.Error("Expected sum > 0")
	}
}

func TestRecordDurationWithError(t *testing.T) {
	h := NewHistogram(DefaultBuckets())

	err := RecordDurationWithError(h, func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if h.Count() != 1 {
		t.Errorf("Expected 1 observation, got %d", h.Count())
	}
}

func TestRecordDurationWithError_Error(t *testing.T) {
	h := NewHistogram(DefaultBuckets())
	expectedErr := assertError{}

	err := RecordDurationWithError(h, func() error {
		time.Sleep(10 * time.Millisecond)
		return expectedErr
	})

	if err != expectedErr {
		t.Errorf("Expected error, got %v", err)
	}

	// Should still record duration
	if h.Count() != 1 {
		t.Errorf("Expected 1 observation even with error, got %d", h.Count())
	}
}

type assertError struct{}

func (e assertError) Error() string {
	return "assertion error"
}

func TestDefaultBuckets(t *testing.T) {
	buckets := DefaultBuckets()

	if len(buckets) == 0 {
		t.Error("Expected default buckets to have values")
	}

	// Verify buckets are sorted
	for i := 1; i < len(buckets); i++ {
		if buckets[i-1] >= buckets[i] {
			t.Errorf("Buckets not sorted: %v", buckets)
		}
	}
}
