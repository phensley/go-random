package stats

import (
	"math"
	"sync/atomic"
)

// Stats records statistics for a metric, e.g. latency in milliseconds
type Stats struct {
	limit     uint64
	count     uint64
	populated uint64
	data      []uint64
	min       uint64
	max       uint64
}

// NewStats ...
func NewStats(limit uint64) *Stats {
	limit++
	data := make([]uint64, limit)
	return &Stats{
		limit:     limit,
		count:     0,
		populated: 0,
		data:      data,
		min:       math.MaxUint64,
		max:       0,
	}
}

// Reset ..
func (s *Stats) Reset() {
	for i := 0; i < len(s.data); i++ {
		s.data[i] = 0
	}
	s.count = 0
	s.populated = 0
	s.min = math.MaxUint64
	s.max = 0
}

// Record ...
func (s *Stats) Record(n uint64) int {
	if n > s.limit {
		return 0
	}

	// Increment metric counter and number of populated metrics
	pop := atomic.AddUint64(&s.data[n], 1) == 1
	atomic.AddUint64(&s.count, 1)
	if pop {
		atomic.AddUint64(&s.populated, 1)
	}

	// Update min and max
	for n < s.min && !atomic.CompareAndSwapUint64(&s.min, s.min, n) {
	}
	for n > s.max && !atomic.CompareAndSwapUint64(&s.max, s.max, n) {
	}
	return 1
}

// Mean ...
func (s *Stats) Mean() float64 {
	if s.count == 0 {
		return 0.0
	}
	sum := uint64(0)
	for i := uint64(s.min); i <= s.max; i++ {
		sum += s.data[i] * i
	}
	return float64(sum) / float64(s.count)
}

// Stddev ...
func (s *Stats) Stddev(mean float64) float64 {
	if s.count < 2 {
		return 0.0
	}
	sum := 0.0
	for i := s.min; i <= s.max; i++ {
		if s.data[i] != 0 {
			sum += math.Pow(float64(i)-mean, 2) * float64(s.data[i])
		}
	}
	return math.Sqrt(sum / float64(s.count-1))
}

// Percentile ...
func (s *Stats) Percentile(p float64) uint64 {
	rank := uint64(Round((p/100.0)*float64(s.count) + 0.5))
	total := uint64(0)
	for i := s.min; i <= s.max; i++ {
		total += s.data[i]
		if total >= rank {
			return i
		}
	}
	return 0
}

// Round n to the nearest integer
func Round(n float64) float64 {
	i, frac := math.Modf(n)
	if frac >= 0.5 {
		return i + 1
	}
	return i
}
