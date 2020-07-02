package poisson

import (
	"math"
	"math/rand"
	"time"
)

// Poisson is ".. a discrete probability distribution expressing
// the probability that a given number of events occur in
// a fixed interval of time ..".
type Poisson struct {
	random    *rand.Rand
	rate      float64
	generator func() int
}

// Option ..
type Option func(p *Poisson)

// WithRand sets the rand implementation
func WithRand(r *rand.Rand) Option {
	return func(p *Poisson) {
		p.random = r
	}
}

// NewPoisson with the given rate and options.
func NewPoisson(rate float64, opts ...Option) *Poisson {
	p := &Poisson{nil, rate, nil}
	if rate < 10 {
		p.generator = p.poisson
	} else {
		p.generator = p.poissonApprox
	}
	for _, o := range opts {
		o(p)
	}
	if p.random == nil {
		seed := int64(time.Now().Nanosecond())
		p.random = rand.New(rand.NewSource(seed))
	}
	return p
}

// Get the next value
func (p *Poisson) Get() int {
	return p.generator()
}

func (p *Poisson) poisson() int {
	l := math.Pow(math.E, -p.rate)
	k := 0
	n := 1.0
	for n > l {
		k++
		u := rand.Float64()
		n = n * u
	}
	return k - 1
}

var (
	logSqrt2Pi = math.Log(math.Sqrt(math.Pi * 2))
)

// For large λ we use Atkinson approximation.
func (p *Poisson) poissonApprox() int {
	beta := math.Pi * (1.0 / math.Sqrt(3.0*p.rate))
	alpha := beta * p.rate
	k := math.Log(0.8065) - p.rate - math.Log(beta)
	x := 0.0
	n := 0
	for {
		for {
			u1 := p.random.Float64()
			x = (alpha - math.Log((1.0-u1)/u1)) / beta
			if x > -0.5 {
				break
			}
		}

		n = int(x + 0.5)
		u2 := p.random.Float64()

		a := math.Exp(alpha - beta*x)
		b := math.Pow(1.0+a, 2.0)
		c := math.Log(u2 / b)
		if alpha-(beta*x)+c <= k+float64(n)*math.Log(p.rate)-factorialLogApprox(n) {
			break
		}
	}
	return n
}

// factorialLogApprox generates Stirling approximation of the
// natural log of the factorial of N
func factorialLogApprox(num int) float64 {
	n := float64(num)
	n3 := n * n * n
	n5 := n3 * n * n
	// Up to n^5 gives adequate precision
	// n7 := n5 * n * n
	a := n*math.Log(n) - n + (0.5 * math.Log(2*math.Pi*n))
	return a + (1.0 / (12.0 * n)) - (1.0 / (360.0 * n3)) + (1.0 / (1260.0 * n5)) //- (1.0 / (1680.0 * n7))
}
