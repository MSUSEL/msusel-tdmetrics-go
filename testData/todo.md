# Tests to Add

## Go Tests

- Check the abstraction of an array:
    `var value [4]byte`

- Test multiple init functions in one package:
    `func init() {}`

- Add circular type param test:
    `type Set[T utils.Comparable[T]] interface {}`

- Test multiple assignments:
    `x, y := 1, 2` and
    `x, y := func()(int, int) { ⋯ }`

- Test literal cast and call:
    `type foo int; func(f foo) bar { ⋯ }; foo(6).bar()`

- Test returning a function literal as an argument.

- Test a method expression call:
    `type foo struct{ ... }; func(f foo) bar { ⋯ }; foo.bar(foo{})`

## Future

- Change access to info to return information for any node as an object
  with all the important features easily accessible. Simplify the duplicate
  checks of types and node/type information using an interface for info.

- Remove arguments and type params, i.e. remove names from function
  arguments and type params by defining generics and functions with
  only the key for the type.

- Add pointer reference to structs that need a pointer. And rethink
  how methods are added to pointers. Maybe create a top level pointer
  for objects that can have the pointer methods added to it. Double
  check when an object can call pointer and non-pointer methods.
