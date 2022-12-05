package dj

import (
	"math/rand"
)

// Uniform returns a random number from a uniform distribution
// between the given min and max values.
func Uniform(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Normal returns a random number from a normal distribution
// with the given mean and standard deviation.
func Normal(mean, stddev float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

// Binomial returns a random number from a binomial distribution
// with the given number of trials and probability of success.
func Binomial(trials int, p float64) int {
	var successes int

	for i := 0; i < trials; i++ {
		if rand.Float64() < p {
			successes++
		}
	}

	return successes
}
