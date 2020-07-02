package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/phensley/go-random/poisson"
)

func main() {
	seed := int64(time.Now().Nanosecond())
	random := rand.New(rand.NewSource(seed))

	iters := 10000
	args := os.Args[1:]
	if len(args) == 1 {
		n, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			panic(err)
		}
		iters = int(n)
	}
	nums := []float64{0.001, 0.01, 0.1, 1, 10, 100, 1000, 10000, 100000, 1000000}
	for {
		for _, n := range nums {
			run(random, iters, n)
		}
		fmt.Println()
	}
}

func run(random *rand.Rand, iters int, rate float64) {
	g := poisson.NewPoisson(rate, poisson.WithRand(random))
	sum := 0
	for i := 0; i < iters; i++ {
		sum += g.Get()
	}
	avg := float64(sum) / float64(iters)
	fmt.Printf("%15.5f ", avg)
}
