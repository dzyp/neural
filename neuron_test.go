package neural 

import (
	"testing"
	"math/rand"
	"log"
)

func TestInputOutput(t *testing.T) {
	numWeights := 10

	weights := make([]*weight, numWeights)
	for i := 0; i < numWeights; i++ {
		ch := make(chan *message)
		weight := makeWeight(ch, rand.Float64())
		weights[i] = weight
	}

	output := make(chan *message)
	makeNeuron(output, weights...)
	for i := 0; i < numWeights; i++ {
		weights[i].channel <- makeMessage(`test`, rand.Float64(), -1)
	}

	result := <- output

	log.Println(result)
}