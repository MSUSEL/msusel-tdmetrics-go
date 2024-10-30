# Generalized Feature Definition

The generalized feature definition is a
[JSON](https://www.json.org/json-en.html) file with specific fields to
define the methods, data, and other constructs of an application.
This definition needs to be flexible enough to handle procedural and
object-oriented languages with some adjustment specific to the language.
For example, Go uses duck-typing but the definition requires a list of
specific implementation so the Go abstractor must perform steps to
predetermine which types would duck-type and define that via implementations.

- [Generalized Feature Definition](#generalized-feature-definition)
  - [Components](#components)
    - [Constructs](#constructs)
    - [Type Descriptions](#type-descriptions)
    - [Declarations](#declarations)
    - [Indices](#indices)
    - [Keys](#keys)
    - [Additional Information](#additional-information)
  - [Schema](#schema)
    - [Project](#project)
    - [Abstract](#abstract)
    - [Argument](#argument)
    - [Basic](#basic)
    - [Field](#field)
    - [Interface Declaration](#interface-declaration)
    - [Interface Description](#interface-description)
    - [Interface Instance](#interface-instance)
    - [Locations](#locations)
    - [Method Instance](#method-instance)
    - [Method](#method)
    - [Metrics](#metrics)
    - [Object Instance](#object-instance)
    - [Object](#object)
    - [Package](#package)
    - [Selection](#selection)
    - [Signature](#signature)
    - [Structure Description](#structure-description)
    - [Type Parameter](#type-parameter)
    - [Value](#value)

## Components

The JSON file has the following components and concepts needed for the schema.

### Constructs

The constructs are any structural or definitional part of the application.
The [Locations](#locations) are not considered constructs since they describe
the organization of the code in files which has no direct effect
on the application.
The following is the list of all the constructs:

[Abstract](#abstract),
[Argument](#argument),
[Basic](#basic),
[Field](#field),
[Interface Declaration](#interface-declaration),
[Interface Description](#interface-description),
[Interface Instance](#interface-instance),
[Method Instance](#method-instance),
[Method](#method),
[Metrics](#metrics),
[Object Instance](#object-instance),
[Object](#object),
[Package](#package),
[Project](#project),
[Selection](#selection),
[Signature](#signature),
[Structure Description](#structure-description),
[Type Parameter](#type-parameter),
[Value](#value)

### Type Descriptions

Type descriptions are any construct that describes a value type.
Think of these as anything that could be used in a parameter, e.g. `int`,
`struct{ x string }`, `Foo[T]`.
The following is the list of all type descriptions:

[Basic](#basic),
[Interface Declaration](#interface-declaration),
[Interface Description](#interface-description),
[Interface Instance](#interface-instance),
[Object Instance](#object-instance),
[Object](#object),
[Signature](#signature),
[Structure Description](#structure-description),
[Type Parameter](#type-parameter)

### Declarations

The declarations are the named constructs defined by the source code,
e.g. `type Foo[T any] struct{ x T }`, `func Bar() string`, `var baz int`.
This does not include the [Abstracts](#abstract), [Fields](#field),
[Type Parameters](#type-parameter), nor named [Arguments](#argument) since
those are only parts of other types.
The following is the list of all declarations.

[Interface Declaration](#interface-declaration),
[Method](#method),
[Object](#object),
[Value](#value)

### Indices

The [Project](#project) contains lists of constructs. When referencing a
construct, where the specific construct kind is known, the index into
the list of that construct in the project can be used.
The indices are all one based to differentiate from default zero values
used by many JSON implementations.
If an item in an object is optional it may use a zero index to indicate
the item is not set to reference anything.

For example, if we know we are referencing a [Basic](#basic) and the
project has `basics: [ "string", "int", "uint" ]`, then an index of 1
will mean `string`, 2 will mean `int`, 3 will mean `uint`.

### Keys

The keys are similar to [Indices](#indices) except for when the construct
kind is not specific. The key is made up of the kind of construct followed
by the index without any space, e.g. `basic2`, `field5`.

Since no construct kind contains a number, a key can be broken up using
a regular expression similar to `^([a-zA-Z]+)(\d+)$`. The the first capture
group will be the construct kind. The kinds are used in the project's items'
names by adding an "s" to any kind that doesn't already end in an "s"
(i.e. `metrics` is still `metrics`) as basic pluralization of the kind name.
The second capture is the index into the list for the project's item.

### Additional Information

The JSON may contain additional information about scoping or how information
was defined. This additional information can be ignored for most cases.

For example `exported` may exist on most named constructs
(declarations, abstracts, and fields) to indicate the scope of that construct
when it was defined in Go.
If exported a Go declaration can be used anywhere inside a project with the
exception of constructs made inside an "internal" package.
For specifics see Go's documentation on exported constructs and
internal packages.
If `exported` is not defined or set to `false` then the declaration is local
to the current package.

Another example is `scope` that may be set `private`, `public`, `internal`, etc.
This is similar to `exported` except for Java instead of Go.

Some additional information can be added on request during the abstraction.
For example the construct `kind` and `index` in the project's list
can be added to the JSON for debugging.

## Schema

The following is the object definitions for the different parts that are
found in the JSON files. When a item in an object is marked as optional
it may not exist, exist but set to the default value,
e.g. `0`, `[]`, `false`, `""`, to indicate it isn't set, or it may exist.
Typically optional things will not be outputted when empty or not set to
reduce the size of the JSON file.

### Project

The project is the top most JSON object. It contains lists of all the
constructs and additional information about the project.

| Name             | Optional | Description |
|:-----------------|:--------:|:------------|
| `abstracts`      | ⬤ | List of [Abstracts](#abstract) |
| `arguments`      | ⬤ | List of [Arguments](#argument) |
| `basics`         | ⬤ | List of [Basics](#basic) |
| `fields`         | ⬤ | List of [Fields](#field) |
| `interfaceDecls` | ⬤ | List of [Interface Declarations](#interface-declaration) |
| `interfaceDescs` | ⬤ | List of [Interface Descriptions](#interface-description) |
| `interfaceInsts` | ⬤ | List of [Interface Instances](#interface-instance) |
| `language`       | ◯ | A string for the source code language, e.g. `go`, `java`. |
| `locs`           | ◯ | [Locations](#locations) |
| `methodInsts`    | ⬤ | List of [Method Instances](#method-instance) |
| `methods`        | ⬤ | List of [Methods](#method) |
| `metrics`        | ⬤ | List of [Metrics](#metrics) |
| `objectInsts`    | ⬤ | List of [Object Instances](#object-instance) |
| `objects`        | ⬤ | List of [Objects](#object) |
| `packages`       | ⬤ | List of [Packages](#package) |
| `selections`     | ⬤ | List of [Selections](#selection) |
| `signatures`     | ⬤ | List of [Signatures](#signature) |
| `structDescs`    | ⬤ | List of [Structure Descriptions](#structure-description) |
| `typeParams`     | ⬤ | List of [Type Parameters](#type-parameter) |
| `values`         | ⬤ | List of [Values](#value) |

### Abstract

An abstract (`abstract`) is a single named signature, i.e a bodiless function,
that exists in an interface.

For example, the following has the abstracts `Foo()` and `Bar(s string) int`.
The first is named `Foo` with the signature `func()` and
the second is named `Bar` with the signature `func(s string) int`.

```Go
interface {
  Foo()
  Bar(s string) int
}
```

| Name        | Optional | Description |
|:------------|:--------:|:------------|
| `exported`  | ⬤ | Boolean defaulted to false. True if the scope is "exported". |
| `index`     | ⬤ | The [Index](#indices) of this abstract in the project's `abstracts` list. |
| `kind`      | ⬤ | `abstract` |
| `name`      | ◯ | The string name for the abstract. |
| `signature` | ◯ | [Index](#indices) for the [Signature](#signature). |

### Argument

An argument (`argument`) is an optionally named parameter or result. For example
`Foo(x string, y int)(ok bool, err error)` has four arguments: `x string`,
`y int`, `ok bool`, `err error`. The same function could be defined without
names `Foo(string, int)(bool, error)` and still have four arguments.
In many cases the name of the argument is ignored since interface abstracts
may have different argument names from the method that it matches with.
Any repeat type argument in Go is expanded, e.g. `Foo(x, y float64)` will
have two arguments, `x float64` and `y float64`.

| Name    | Optional | Description |
|:--------|:--------:|:------------|
| `index` | ⬤ | The [Index](#indices) of this argument in the project's `arguments` list. |
| `kind`  | ⬤ | `argument` |
| `name`  | ⬤ | The optional string name for the argument. |
| `type`  | ◯ | [Key](#keys) for any [Type Description](#type-descriptions). |

### Basic

A basic type (`basic`) defines a built-in type being used by the application,
e.g. `int`, `float64`, `float`, `double`. The basics between Go and Java
don't overlap perfectly. Things that are aliased, such as `rune` is an alias
for `uint32`, are un-aliased. No methods are attached to the basic types.
If the language specifies methods for the type then an interface can be
defined, such as `Integer`, to home those methods.

Things like `any` and `Object` are empty [interfaces](#interface-declaration)
not basic types even though the empty interface is used as a subtype
to all other types.
Compound built-in types, such as `error` and `complex64`, are replaced by
interfaces with the components of the compound type accessible through
the [abstracts](#abstract) of that interface.

| Name    | Optional | Description |
|:--------|:--------:|:------------|
| `index` | ⬤ | The [Index](#indices) of this basic in the project's `basics` list. |
| `kind`  | ⬤ | `basic` |
| `name`  | ◯ | The string name for the basic type, e.g. `int`, `bool`, `string`. |

### Field

A field (`field`) is a single named type inside of
a [structure](#structure-description).
For example, the following has the fields `x int` and `y string`.

```Go
struct {
  x int
  y string
}
```

Any repeat type fields in Go is expanded, e.g. `struct { x, y float64 }` will
have two fields, `x float64` and `y float64`.

| Name       | Optional | Description |
|:-----------|:--------:|:------------|
| `embedded` | ⬤ | Boolean defaulted to false. True if the field is from an embedded struct. |
| `exported` | ⬤ | Boolean defaulted to false. True if the scope is "exported". |
| `index`    | ⬤ | The [Index](#indices) of this field in the project's `fields` list. |
| `kind`     | ⬤ | `field` |
| `name`     | ◯ | The string name for the field. |
| `type`     | ◯ | [Key](#keys) for any [Type Description](#type-descriptions). |

### Interface Declaration

An interface declaration (`interfaceDecl`) is a named definition of an
interface. For example, the following is the declaration of `Foo` with
a [type parameter](#type-parameter) `T` and an [abstract](#abstract) `Bar`.

```Go
type Foo[T any] interface {
  Bar(value T) string
}
```

| Name         | Optional | Description |
|:-------------|:--------:|:------------|
| `exported`   | ⬤ | Boolean defaulted to false. True if the scope is "exported". |
| `index`      | ⬤ | The [Index](#indices) of this interface declaration in the project's `interfaceDecls` list. |
| `instances`  | ⬤ | List of [Indices](#indices) to [Interface Instances](#interface-instance). |
| `interface`  | ◯ | The [Index](#indices) to the declared [Interface](#interface-description) type. |
| `kind`       | ⬤ | `interfaceDecl` |
| `loc`        | ⬤ | The [Location](#locations) offset. |
| `name`       | ◯ | The name of the declared interface. |
| `package`    | ◯ | The [Index](#indices) to the [Package](#package) this declaration is declared in. |
| `typeParams` | ⬤ | List of [Indices](#indices) to [Type Parameters](#type-parameter) if this interface is generic. |

### Interface Description

An interface description (`interfaceDesc`) describes the type of an interface.
This interface type may be the type for an interface declaration, interface
instances, and interface literal, e.g. `interface { String() string }`.

| Name         | Optional | Description |
|:-------------|:--------:|:------------|
| `abstracts`  | ⬤ | List of [Indices](#indices) to [Abstracts](#abstract). |
| `approx`     | ⬤ | List of [Keys](#keys) to any [Type Description](#type-descriptions) for approximate constraints. |
| `exact`      | ⬤ | List of [Keys](#keys) to any [Type Description](#type-descriptions) for exact constraints. |
| `index`      | ⬤ | The [Index](#indices) of this interface in the projects' `interfaceDescs` list. |
| `inherits`   | ⬤ | List of [Indices](#indices) to inherited [Interfaces](#interface-description). |
| `kind`       | ⬤ | `interfaceDesc` |
| `package`    | ⬤ | The [Index](#indices) to the [Package](#package) this interface is pinned to. |

### Interface Instance

An interface instance (`interfaceInst`) is an instantiation of a generic
interface declaration.
The instance types are the type arguments used in the type parameters.
The instance types may be type parameters as well as a fully realized type.
For example, `type Foo[T any] interface { Value() T }` with instance type `bool`
will create `type Foo[bool] interface { Value() bool }`.

| Name            | Optional | Description |
|:----------------|:--------:|:------------|
| `generic`       | ◯ | The [Index](#indices) of the generic [Interface Declaration](#interface-declaration) this is an instance of. |
| `index`         | ⬤ | The [Index](#indices) of this interface in the projects' `interfaceInsts` list. |
| `instanceTypes` | ◯ | List of [Keys](#keys) to any [Type Description](#type-descriptions) for exact constraints. |
| `kind`          | ⬤ | `interfaceInst` |
| `resolved`      | ◯ | The [Index](#indices) to the [Interface Type Description](#interface-description) this instance defines. |

### Locations

TODO: Add comment

### Method Instance

TODO: Add comment

### Method

TODO: Add comment

### Metrics

TODO: Add comment

### Object Instance

TODO: Add comment

### Object

TODO: Add comment

### Package

TODO: Add comment

### Selection

TODO: Add comment

### Signature

TODO: Add comment

### Structure Description

TODO: Add comment

### Type Parameter

TODO: Add comment

### Value

TODO: Add comment
