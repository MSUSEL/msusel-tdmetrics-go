# Experiment 001 - Participation

:seedling: *This is a work in progress*

- [Experiment Background](#experiment-background)
- [Running Experiment](#running-experiment)

This is an experiment for parsing Go files and determining all the times
that a type is used as either the receiver or parameter for a function.
This participation metric could be useful in producing similar results to
determining if a [God object](https://en.wikipedia.org/wiki/God_object).

## Experiment Background

In OO languages we use three metrics to detect a God object:

1. **Weighted Method Count (WMC)** is the sum of the cyclomatic
   complexity for all the methods in the class.
2. **Tight Class Cohesion (TCC)** is the ratio of methods directly
   connected to each other through instance variables in the class
   and the total possible connections that could have existed.
3. **Access to Foreign Data (ATFD)** is the number of connection from
   this class to another class’s instance variables either directly,
   via accessor methods, or via attributes.

In a procedural languages, like Go, methods are not associated directly with a structure.
In Go the receiver is not much more than a specialize parameter
which allows the parser to handle a selector invocation (`cat.meow()` vs `meow(cat)`).
The receiver can be `nil` object or an `int`, function pointer, or other non-structure type
making it not the same as a method on a class.

Receivers are a desired pattern in most cases, so simply considering the
receiver would bias development towards not using receivers and just
making the first parameter the "main" object.
If the participation of an object type anywhere in the receiver or parameters
was counted equally then we could consider that participation as
the same as partial ownership of that method.
That participation can then be weighted by the cyclomatic complexity of the function
to take the place of the WMC. The weighting will contribute to the sum of
all the types used in the function's receiver or parameters.

The Weighted Participation Counts (WPC) that are over some threshold
along with a procedural form of TCC and ATFD would indicate the probability
that a structure is being used like a God object.

## Running Experiment

Run the experiment with

`go run main.go <base path>`

where the base path is the path to the root folder of the Go files to evaluate.

Test and vendor files will be ignored by default.
You may want to configure the main method as needed to check specific files.
