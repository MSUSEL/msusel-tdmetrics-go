# Go Abstractor

- [Background](#background)
- [Running Abstractor](#running-abstractor)

## Background

In OO languages we use three metrics to detect a form of technical debt
called a God [God object](https://en.wikipedia.org/wiki/God_object).
God objects are also class bloaters or junk drawers are classes that are
too big and complex, so should be broken up into smaller classes.

1. **Weighted Method Count (WMC)** is the sum of the cyclomatic
   complexity for all the methods in the class.
2. **Tight Class Cohesion (TCC)** is the ratio of methods directly
   connected to each other through instance variables in the class
   and the total possible connections that could have existed.
3. **Access to Foreign Data (ATFD)** is the number of connection from
   this class to another class’s instance variables either directly,
   via accessor methods, or via attributes.

All of those three metrics require knowledge of which methods
are members of which classes. For example, the WMC is the sum of all
the methods that are members of a class. For our uses, a method
and function are considered the same, they are a set of instructions
to modify data. Also a class and structure are considered the same, they
are a collection of data.

In a procedural languages, like Go, methods are not associated directly
with a structure. In Go, the receiver is not much more than a specialized
parameter which allows the parser to handle a selector invocation
(`cat.meow()` vs `meow(cat)`) and work like a closure to match method
signatures. The receiver can be a `nil` object (`nil` is a null value with
a specific type) or an `int`, function pointer, or other non-structure type
making it not quiet the same as a method on a class.

To simply considering receivers as indication of membership in technical
debt analysis would bias developer towards not using receivers and just
making the first parameter a receiver or using a closure. Instead we
should determine how structures participate in a functions purpose.
The participation is similar to membership except that participation
is a fuzzy estimate between zero and one for all structures participating.

The participation score of methods to structures can be used to fractionally
weight the cyclomatic complexity of methods in the sum for WMC and the
other technical debt metrics.

But, how do we determine the participation of methods to structures?
We abstract relevant structural information from the static source code
that can be used in a Design Recovery algorithm. Design Recovery was
designed to associate data and functions into class estimates.
This code will not perform Design Recovery, only abstract the relevant
information that can be used to later perform the design recovery.

## Running Abstractor

Run the abstractor with:

```bash
cd ./goAbstractor
go run ./main.go -i <package path> -o <result json>
```

The package path is the path to the directory containing the go mod file
for the project or path to package to abstract. The result json is the path
to the file to write the abstraction json out to. Add `-m` to minimize the
json output file.

For more information about arguments run:

```bash
go run ./main -h
```
