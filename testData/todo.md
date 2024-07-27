# Tests to Add

## Go Tests

- Test when a struct has a type param that isn't used in any field:
    `type A[T] struct { value int }; func (a A[T]) Foo(v A)`

- Test when a struct has a type param that isn't used in any method:
    `type A[T] struct { value A }; func (a A[T]) Foo() int`

- Test when a struct has a type param that is only used for another type param:
    `type A[T any, S []T] { list S }; func (a A[T, S]) Foo() string`

- Test when a type parameter has an approximate type:
    `type A[T any, S ~[]T] { list S }; func (a A[T, S]) Foo() T`

- Test when a type parameter has an approximate with function,
   see [generics spec](https://go.dev/ref/spec#General_interfaces):
    `interface { ~int; String() string }`

- Test multiple init functions in one package:
    `func init() {}`

- Test when there are identical identifiers for functions
  but all but one has unique receivers.

- Check the abstraction of an array:
    `var value [4]byte`
