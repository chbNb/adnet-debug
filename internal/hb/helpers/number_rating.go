package helpers

import "math/rand"

func NumberRating(numberRating int) int {
	if numberRating > 10000 {
		return numberRating
	}
	return rand.Intn(40001) + 10000
}

func RandFloat(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}

func RandFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}
