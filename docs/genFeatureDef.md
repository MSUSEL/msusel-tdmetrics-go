# Generalized Feature Definition

The generalized feature definition is a
[JSON](https://www.json.org/json-en.html) file with specific fields to
define the methods and data of an application. This definition needs to be
flexible enough to handle procedural and object-oriented languages with some
adjustment specific to the language. For example, Go uses duck-typing
but the definition requires a list of specific implementation so the Go
abstractor must perform steps to predetermine which types would duck-type
and define that via implementations.

The JSON file contains a top level object with the following fields:
[`language`](#language),
[`interfaces`](#interfaces),
[`signatures`](#signatures),
[`structs`](#structs),
[`typeParams`](#typeparams),
and [`packages`](#packages)

## language

The language field's value is a string defining the language that was abstracted,
e.g. `"language": "go"`.

## interfaces

The interfaces field's value is an array of [interface objects](#interface-object).
The interfaces are unnamed and reduced to have no repeats.

When referencing a type by index these indices are between
$[\;0\;..\;|\text{interfaces}|\;)$
where type index of $x$ is the interface object $\text{interfaces}[\;x\;]$.

Typically the first interface is an empty interface, `{}`, that represents
the base object, i.e. `any` in Go, `Object` in Java.

### interface object

Each interface object has two arrays:

- `inherits`: This is an array of indices of types. Each type should be
    another interface in the top level `interfaces`.
- `methods`: This is an array of method objects. Each method
    object has a `name` and `signature`.
  - `name` is a string for the identifier of the method defined in the interface.
  - `signature` is an index of a signature type.

Example:

```json
{
    "inherits": [ 1, 4 ],
    "methods": [
        { "name": "String", "signature": 612 },
        { "name": "Count",  "signature": 755 }
    ]
}
```

```go
interface {
    String() string
    Count() int
}
```

## signatures

The signatures field's value is an array of [signature objects](#signature-object).
When referencing a type by index these indices are between
$[\;\text{offset}\;..\;\text{offset}+|\text{signatures}|\;)$
where $\text{offset} = |\text{interfaces}|$ and
type index of $x$ is the interface object $\text{signatures}[\;x-\text{offset}\;]$.

### signature object

The signature objects represents a function signature without a name or body.
The signature has the following fields:

- `variadic`: Is boolean where true indicates the last parameter will be a
  list and is a variadic parameter.
- `params`: This is a list of names and types for each parameter in the order
   the parameters appear in the signature.
  - `name`: Each name should be unique for a parameter.
  - `type`: The [type](#type) of the parameter.
- `return`: The [type](#type) to return. Must be one type only, meaning for Go
  when multiple types are returned the return value needs to be converted
  to a struct type.
- `typeParam`: The type parameters for the signature.

Example:

```json
 {
    "params": [
        { "name": "name", "type": "string" },
        { "name": "age",  "type": "int" }
    ],
    "return": {
        "elem": "Person",
        "kind": "pointer"
    }
},
```

```go
func(name string, age int) *Person
```

## structs

TODO: Finish

## typeParams

TBD

## packages

TODO: Finish

## type

TODO: Finish
