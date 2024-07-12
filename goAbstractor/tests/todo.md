# Tests to Add

- Test when a struct has a type param that isn't used in any field:
    `type A[T] struct { value int }; func (a A[T]) Foo(v A)`

- Test when a struct has a type param that isn't used in any method:
    `type A[T] struct { value A }; func (a A[T]) Foo() int`

- Test when a struct has a type param that is only used for another type param:
    `type A[T any, S []T] { list S }; func (a A[T, S]) Foo() string`
