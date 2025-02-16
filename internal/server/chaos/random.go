package chaos

import "math/rand"

func selectWeightedRandom(elements []FailureType, probabilities []float64) FailureType {
	cumulativeProbabilities := make([]float64, len(probabilities))
	cumulativeProbabilities[0] = probabilities[0]
	for i := 1; i < len(probabilities); i++ {
		cumulativeProbabilities[i] = cumulativeProbabilities[i-1] + probabilities[i]
	}

	randomNumber := rand.Float64()

	for i := 0; i < len(cumulativeProbabilities); i++ {
		if randomNumber < cumulativeProbabilities[i] {
			return elements[i]
		}
	}

	return elements[len(elements)-1]
}
