# Package `genetic`

Package `genetic` implements genetic algorithms using Golang's Generics support.

## What are Genetic Algorithms?

Genetic algorithms are a form of machine learning software inspired by the process of natural selection. GA's can quickly approximate solutions for problems with very large search spaces. They are often used to tackle NP problems whose solutions can vary in levels of quality. In these kinds of problems, measuring a solution's validity and quality are fast, but yet there is no easy way to determine the _optimal_ solution without guessing and checking. Genetic Algorithms provide a means to get a _pretty good_ solution to any such problem. Examples include:

- Transportation routing & logistics
- Storage Optimization (Knapsack Problem)
- 'Travelling Salesman' problems
- Machine Learning meta-parameter tuning
- Hardware and Engineering Design optimization
- Scheduling

GAs work by iteratively evolving a finite pool of possible solutions to your problem. With each iteration, the best solutions are mixed and matched with one another to produce novel solutions, and then mutated to introduce variety. This process repeats until some termination condition is reached, at which point the algorithm terminates and returns the best found solution.

## Word Thingies

Terminology that's important to know before using Genetic Algorithms:

- Genome: One possible evolutionary candidate. AKA "Solution".
- Population: A pool of genomes.
- Fitness: A scoring metric which objectively ranks the quality of any genome.
- Crossover: The process of combining information from parent genomes to produce novel child genomes.
- Mutations: Random variations in genome structure or content occurring during crossover.
- Selection: The process of picking which genomes will crossover to produce the next generation of genomes.
- Elitism: Preferential treatment to the highest-scoring genomes, allowing them to survive into the next generation. This prevents regression and loss of quality solutions.
- Allele: A candidate gene for a specific location on the Genome. If two genes are alleles of each other, they are evolutionary competitors.

## Usage

Find a full example in [`example_test.go`](./example_test.go).

Let's say we want to guess a 33-character string, `"how did you ever guess my secret!"`. Our goal is to evolve a population of random strings into the target string.

Thanks to Generics, `genetic` can be unopinionated about the data type you use to represent your genomes. Instead, you define several functions acting on those genomes, which in concert will define the evolutionary behavior of a population. In our example, we will store our genomes as byte slices.

We need a `GenesisFunc[T]` which initializes and returns random genomes of type `T`, where in our case `T` is `[]byte`.

```go
func genesis() []byte {
  genome := make([]byte, 33)
  rand.Read(genome)
  return genome
}
```

We need a `FitnessFunc[T]` which scores every genome on how many of its bytes match the target string. We'll use a minimum fitness of 1, which means our maximum fitness is 34.

```go
func fitnessFn(guesses [][]byte, fitnesses []int) {
  for i, guess := range guesses {
    fitnesses[i] = 1
    for j, b := range guess {
      if b == bytesToGuess[j] {
        fitnesses[i]++
      }
    }
  }
}
```

_Since this is a static fitness function (its output is not dependent on the other genomes in the population), we can optimize this slightly with the `StaticFitnessFunc` utility._

```go
fitnessFn := genetic.StaticFitnessFunc(func(guess []byte) int {
  fitness := 1
  for j, b := range guess {
    if b == bytesToGuess[j] {
      fitness++
    }
  }
  return fitness
})
```

We need a `CrossoverFunc[T]` which combines two genomes together to produce offspring. `genetic` provides several simple crossover functions which work on any generic slice type genome. Let's use `UniformCrossover`, which produces offspring whose genomes are a random mix of their parents' genomes.

```go
crossover := genetic.UniformCrossover[[]byte]
```

We need a `SelectionFunc[T]` which will pair genomes off into mating pairs, usually based on their fitnesses. A good selection function should give higher-fitness genomes more opportunities to mate than lower-fitness genomes, but should also ensure that the elite (highest-fitness) genomes do not completely dominate all mating pairs, as genetic homogeneity can lead to evolutionary stagnation.

`genetic` also provides a couple of classic selection functions which work with any genome type. Let's use `TournamentSelection`, with a tournament size of 3 - randomly selected genomes will be repeatedly subjected to 'tournaments' where the highest-fitness entrants will become mating candidates.

```go
selection := genetic.TournamentSelection[[]byte](3)
```

Finally, we need a `MutationFunc[T]` which should randomly (usually with some small probability) mutate the genomes of each new generation. This injects some diversity, helping to explore more optimal solutions.

```go
func mutate(guess []byte) {
  r := make([]byte, len(guess))
  rand.Read(r)
  for i := range guess {
    if rand.Float64() < 0.05 { // 5% mutation rate
      guess[i] ^= r[i]
    }
  }
}
```

Now we can construct a `Population[[]byte]` instance:

```go
population := genetic.NewPopulation(
  100, // population size
  genesis,
  crossover,
  fitness,
  selection,
  mutation,
)
```

This population now contains 100 randomized 33-byte strings. Let's try to evolve them into our target string based solely on how many bytes they guessed correctly.

```go
population.Evolve(
  34,    // minimum fitness. Evolution terminates upon finding a solution with this fitness or higher.
  2000,  // max generations. Evolution terminates after iterating through this many generations.
  2,     // elitism. Carries over this many of the highest-fitness genomes from each previous generation into the next.
)
```

We can now pull out our best candidate to see how we did.

```go
bestSolution, bestFitness := population.Best()
fmt.Printf("evolved string: %q\n", bestSolution)
fmt.Println("best fitness:", bestFitness)
```

If we did everything right, we should see our population evolved into our target string and reached maximum fitness 34:

```
evolved string: "how did you ever guess my secret!"
best fitness: 34
```
