# Genetic algorithm implementation

## Basic rules

- A location may only appear once in a given slot
- A worker may only appear once in a given slot

## Chromosome representation

Each position in the chromosome should indicate a task in the schedule. At this position
we store a tuple (worker, location, slot).

A way of eliminating the same location being used two times in slot?

## Fitness

Fitness is calculated based on the number of constraints fulfilled.
