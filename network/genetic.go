package network

import (
	"fmt"
	"math/rand"
)

type FeedForward func(input []float64) []float64
type evaluateIndividual func(FeedForward) float64

type GeneticAlgorithm struct {
	populationSize int
	sizes          []int
	population     []*individual

	generationNumber int

	mutationChance float64
	mutationRate   float64

	evaluate evaluateIndividual
}

func NewGeneticAlgoritm(populationSize int, sizes []int, mutationChance, mutationRate float64, evaluate evaluateIndividual) *GeneticAlgorithm {
	ga := &GeneticAlgorithm{
		populationSize: populationSize,
		sizes:          sizes,
		mutationChance: mutationChance,
		mutationRate:   mutationRate,
		evaluate:       evaluate,
	}

	ga.generatePopulation()

	return ga
}

func (ga *GeneticAlgorithm) generatePopulation() {
	ga.population = make([]*individual, ga.populationSize)
	for i := 0; i < ga.populationSize; i++ {
		ga.population[i] = newIndividual(ga.sizes)
	}
}

func (ga *GeneticAlgorithm) EvaluateGeneration() {
	fitnessTrack := make([]float64, 5)

	for _, individual := range ga.population {
		for i := 0; i < len(fitnessTrack); i++ {
			fitnessTrack[i] = ga.evaluate(individual.feedForward)
		}

		// median fitness
		individual.fitness = fitnessTrack[len(fitnessTrack)/2]
	}
}

func (ga *GeneticAlgorithm) GetBestIndividual() FeedForward {
	best := ga.population[0]

	for _, individual := range ga.population {
		if individual.fitness > best.fitness {
			best = individual
		}
	}

	// print generation and best fitness
	fmt.Printf("Generation: %d, Fitness: %f", ga.generationNumber, best.fitness)

	return best.feedForward
}

func (ga *GeneticAlgorithm) tournamentSelection(tournamentSize int) *individual {
	best := ga.population[rand.Intn(len(ga.population))]

	for i := 0; i < tournamentSize-1; i++ {
		individual := ga.population[rand.Intn(len(ga.population))]
		if individual.fitness > best.fitness {
			best = individual
		}
	}

	return best
}

func (ga *GeneticAlgorithm) crossoverParents(parent1, parent2 *individual) *individual {
	// TODO: newIndividual will generate random weights, performance boost possible if we fix that
	child := newIndividual(ga.sizes)

	for synapse := 0; synapse < len(parent1.weights); synapse++ {
		for weightIndex := 0; weightIndex < len(parent1.weights[synapse]); weightIndex++ {
			if rand.Float64() < 0.5 {
				child.weights[synapse][weightIndex] = parent1.weights[synapse][weightIndex]
			} else {
				child.weights[synapse][weightIndex] = parent2.weights[synapse][weightIndex]
			}
		}

		for biasIndex := 0; biasIndex < len(parent1.biases[synapse]); biasIndex++ {
			if rand.Float64() < 0.5 {
				child.biases[synapse][biasIndex] = parent1.biases[synapse][biasIndex]
			} else {
				child.biases[synapse][biasIndex] = parent2.biases[synapse][biasIndex]
			}
		}
	}

	return child
}

func (ga *GeneticAlgorithm) mutateIndividual(individual *individual) {
	for synapse := 0; synapse < len(individual.weights); synapse++ {
		for weightIndex := 0; weightIndex < len(individual.weights[synapse]); weightIndex++ {
			if rand.Float64() < ga.mutationChance {
				individual.weights[synapse][weightIndex] += (rand.Float64()*2 - 1) * ga.mutationRate
			}
		}

		for biasIndex := 0; biasIndex < len(individual.biases[synapse]); biasIndex++ {
			if rand.Float64() < ga.mutationChance {
				individual.biases[synapse][biasIndex] += (rand.Float64()*2 - 1) * ga.mutationRate
			}
		}
	}
}

func (ga *GeneticAlgorithm) EvolveGeneration() {
	newPopulation := make([]*individual, ga.populationSize)
	for i := 0; i < ga.populationSize; i++ {
		parent1 := ga.tournamentSelection(10)
		parent2 := ga.tournamentSelection(10)
		child := ga.crossoverParents(parent1, parent2)
		ga.mutateIndividual(child)
		newPopulation[i] = child
	}
	ga.population = newPopulation
	ga.generationNumber++
}
