package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsPercentile(t *testing.T) {
	s := NewStats(100)
	for i := uint64(1); i < 100; i++ {
		s.Record(i)
	}

	assert.Equal(t, 50.0, s.Mean())

	assert.Equal(t, uint64(10), s.Percentile(10.0))
	assert.Equal(t, uint64(25), s.Percentile(25.0))
	assert.Equal(t, uint64(50), s.Percentile(50.0))
	assert.Equal(t, uint64(75), s.Percentile(75.0))
	assert.Equal(t, uint64(90), s.Percentile(90.0))
	assert.Equal(t, uint64(95), s.Percentile(95.0))
}

func TestStatsMeanStddev(t *testing.T) {
	s := NewStats(100)
	assert.Equal(t, 0.0, s.Mean())
	assert.Equal(t, 0.0, s.Stddev(s.Mean()))

	s.Record(10)
	assert.Equal(t, 10.0, s.Mean())
	assert.Equal(t, 0.0, s.Stddev(s.Mean()))

	s.Record(20)
	s.Record(30)
	assert.Equal(t, 20.0, s.Mean())
	assert.Equal(t, 10.0, s.Stddev(s.Mean()))

	s.Reset()
	s.Record(25)
	s.Record(50)
	s.Record(75)
	assert.Equal(t, uint64(50), s.Percentile(50.0))
	assert.Equal(t, 50.0, s.Mean())
	assert.Equal(t, 25.0, s.Stddev(s.Mean()))

	s.Reset()
	s.Record(20)
	s.Record(30)
	s.Record(40)
	assert.Equal(t, uint64(30), s.Percentile(50.0))
	assert.Equal(t, 30.0, s.Mean())
	assert.Equal(t, 10.0, s.Stddev(s.Mean()))
}
