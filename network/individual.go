package network

import (
	"math"
	"math/rand"
)

type individual struct {
	weights [][]float64
	biases  [][]float64
	fitness float64
}

// sizes is the number of neurons in each layer, including the input and output layers
func newIndividual(sizes []int) *individual {
	weights := make([][]float64, len(sizes)-1)
	biases := make([][]float64, len(sizes)-1)
	for i := 0; i < len(sizes)-1; i++ {
		weights[i] = make([]float64, sizes[i]*sizes[i+1])
		biases[i] = make([]float64, sizes[i+1])

		for j := 0; j < len(weights[i]); j++ {
			weights[i][j] = rand.Float64()*2 - 1
		}

		for j := 0; j < len(biases[i]); j++ {
			biases[i][j] = rand.Float64()*2 - 1
		}
	}

	return &individual{weights: weights, biases: biases, fitness: 0}
}

func (i *individual) feedForward(input []float64) []float64 {
	// confirm that the input is the correct size
	if len(input) != (len(i.weights[0]) / len(i.biases[0])) {
		panic("Input size does not match network input size")
	}

	// loop through each synapse (between the layers)
	for synapseIndex := 0; synapseIndex < len(i.weights); synapseIndex++ {
		output := make([]float64, len(i.biases[synapseIndex]))

		// loop through each neuron in the INPUT layer
		for inputNeuronIndex := 0; inputNeuronIndex < len(input); inputNeuronIndex++ {
			// loop through each neuron in the NEXT layer
			for outputNeuronIndex := 0; outputNeuronIndex < len(output); outputNeuronIndex++ {
				// add up the activations
				output[outputNeuronIndex] += input[inputNeuronIndex] * i.weights[synapseIndex][inputNeuronIndex*len(output)+outputNeuronIndex]
			}
		}

		// add the biases
		for outputNeuronIndex := 0; outputNeuronIndex < len(output); outputNeuronIndex++ {
			output[outputNeuronIndex] += i.biases[synapseIndex][outputNeuronIndex]
		}

		// apply the activation function
		for outputNeuronIndex := 0; outputNeuronIndex < len(output); outputNeuronIndex++ {
			if synapseIndex == len(i.weights)-1 {
				// apply the sigmoid function to output layers
				output[outputNeuronIndex] = 1 / (1 + math.Exp(-output[outputNeuronIndex]))
			} else {
				// apply the leaky relu function to hidden layers
				if output[outputNeuronIndex] < 0 {
					output[outputNeuronIndex] = 0.01 * output[outputNeuronIndex]
				}
			}
		}

		// apply the activation function

		// set input to output so that the next layer can use it
		input = output
	}

	return input
}
