# Extending Pointer

Go allows functions to be called from pointers and methods to be defined
with a pointer receiver. When implementing inheritance to replace duck-typing,
the extending of a pointer is needed. In the following code snippet the
expressed how underlying types and pointers determine what may be done
to an object.

```Go
package main

import "fmt"

type A map[string]int

func (a A) Foo() { fmt.Println(`Foo`) }

func (a *A) Bar() { fmt.Println(`Bar`) }

type B interface {
  Foo()
  Bar()
}

func main() {
  var a A = A{}
  a[`a`] = 1
  a.Foo()
  a.Bar()

  var pa *A = &a
  //pa[`pa`] = 2 // invalid operation: cannot index pa (variable of type *A)
  pa.Foo()
  pa.Bar()

  var ppa **A = &pa
  //ppa[`ppa`] = 3 // invalid operation: cannot index ppa (variable of type **A)
  //ppa.Foo() // ppa.Foo undefined (type **A has no field or method Foo)
  //ppa.Bar() // ppa.Bar undefined (type **A has no field or method Bar)
  _ = ppa

  //var ba B = a // cannot use a (variable of type A) as B value in variable declaration: A does not implement B (method Bar has pointer receiver)

  var bpa B = pa
  //bpa[`bpa`] = 4 // invalid operation: cannot index bpa (variable of type B)
  bpa.Foo()
  bpa.Bar()

  var pbpa *B = &bpa
  //pbpa[`pbpa`] = 6 // invalid operation: cannot index pbpa (variable of type *B)
  //pbpa.Foo() // pbpa.Foo undefined (type *B is pointer to interface, not interface)
  //pbpa.Bar() // pbpa.Bar undefined (type *B is pointer to interface, not interface)
  _ = pbpa
}
```
