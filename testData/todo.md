# Tests to Add

## Go Tests

- Test creating multiple return struct with type param:
    `func (s S[T,U]) func() ([]T, *U)`

- Test when there are identical identifiers for functions
  but all but one has unique receivers.

- Check the abstraction of an array:
    `var value [4]byte`

- Multiple packages

- Multiple files for a package

- Test multiple init functions in one package:
    `func init() {}`

- A named type of a function via a definition with a signature
    and one typed by copying the data from another type:
    `type A[T any] func(value T); type B[X any] A[X]; type C A[int]`
