package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/phensley/go-random/boxmuller"
	"github.com/phensley/go-random/stats"
)

func main() {
	seed := int64(time.Now().Nanosecond())
	r := rand.New(rand.NewSource(seed))

	iters := 100
	mean := 1.0
	stddev := 0.005

	gen := boxmuller.NewBoxMuller(mean, stddev, boxmuller.WithRand(r))
	for i := 0; i < iters; i++ {
		generate(gen)
	}
}

func generate(gen *boxmuller.BoxMuller) {

	// Generate a bunch of random numbers with a gaussian distribution
	lim := 50000
	nums := make([]float64, lim)

	min := 0.0
	max := 0.0
	for i := 0; i < lim; i++ {
		r := gen.Next()
		nums[i] = r
		if min == 0 || r < min {
			min = r
		}
		if max == 0 || r > max {
			max = r
		}
	}

	// Display
	fmt.Printf("min=%f  max=%f  avg=%f", min, max, avg(nums))

	// Buckets
	b := 50
	scale := (max - min) / float64(b)

	buckets := make([]int, b+1)
	maxb := 0
	for _, n := range nums {
		i := int((n - min) / scale)
		buckets[i]++
		if buckets[i] > maxb {
			maxb = buckets[i]
		}
	}

	width := 100
	nscale := float64(maxb) / float64(width)

	for i, bucket := range buckets {
		fmt.Printf("%02d  %.2f  ", i+1, (float64(i)*scale)+min)
		for j := 0; j < int(float64(bucket)/nscale); j++ {
			fmt.Print("*")
		}
		fmt.Println()
	}

	// Feed into a stats structure and confirm the mean is close
	// to the expected value
	stats := stats.NewStats(uint64(max * 10000))
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for _, n := range nums {
				stats.Record(uint64(n * 10000))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("statistics-computed mean: ", stats.Mean()/float64(10000))
	fmt.Println("statistics-computed std-dev: ", stats.Stddev(stats.Mean())/float64(10000))
	fmt.Println("---------------------------------------------------")
}

func avg(nums []float64) float64 {
	s := 0.0
	for _, n := range nums {
		s += n
	}
	return s / float64(len(nums))
}
