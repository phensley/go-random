package boxmuller

import (
	"math"
	"math/rand"
	"time"
)

// BoxMuller generates random numbers in a gaussian distribution,
// drawn from an underlying source of uniform random numbers.
type BoxMuller struct {
	random  *rand.Rand
	mean    float64
	stddev  float64
	last    float64
	useLast bool
}

// Option for configuration
type Option func(*BoxMuller)

// WithRand sets the rand implementation
func WithRand(r *rand.Rand) Option {
	return func(b *BoxMuller) {
		b.random = r
	}
}

// NewBoxMuller converts a sequence of uniform random numbers 0..1 into
// a normally-distributed sequence with given mean, standard deviation,
// and options.
func NewBoxMuller(mean, stddev float64, opts ...Option) *BoxMuller {
	b := &BoxMuller{nil, mean, stddev, 0, false}
	for _, o := range opts {
		o(b)
	}
	if b.random == nil {
		seed := int64(time.Now().Nanosecond())
		b.random = rand.New(rand.NewSource(seed))
	}
	return b
}

// Next returns the next random number.
func (b *BoxMuller) Next() float64 {
	random := b.random
	var (
		w  float64
		x1 float64
		x2 float64
		y1 float64
		y2 float64
	)
	if b.useLast {
		y1 = b.last
		b.useLast = false
	} else {
		for w == 0.0 || w >= 1.0 {
			x1 = 2.0*random.Float64() - 1.0
			x2 = 2.0*random.Float64() - 1.0
			w = x1*x1 + x2*x2
		}

		w = math.Sqrt((-2.0 * math.Log(w)) / w)
		y1 = x1 * w
		y2 = x2 * w
		b.useLast = true
		b.last = y2
	}
	return b.mean + y1*b.stddev
}
