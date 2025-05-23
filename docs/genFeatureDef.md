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
    - [Extra Information](#extra-information)
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
    - [Method](#method)
    - [Method Instance](#method-instance)
    - [Metrics](#metrics)
    - [Object](#object)
    - [Object Instance](#object-instance)
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
[Method](#method),
[Method Instance](#method-instance),
[Metrics](#metrics),
[Object](#object),
[Object Instance](#object-instance),
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
[Object](#object),
[Object Instance](#object-instance),
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

### Extra Information

The JSON may contain additional information about scoping or how information
was defined. This additional information can be ignored for most cases.

For example `vis` may exist on most named constructs
(declarations, abstracts, and fields) to indicate the scope of that construct
when it was defined in Go.
If "exported" a Go declaration can be used anywhere inside a project with the
exception of constructs made inside an "internal" package.
For specifics see Go's documentation on exported constructs and
internal packages.
If `vis` is not defined or set to empty then the declaration is local
to the current package.
For Java `vis` may be set `private`, `public`, `internal`, etc.

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

| Name             | Optional | Extra | Description |
|:-----------------|:--------:|:-----:|:------------|
| `abstracts`      | ⬤ | ◯ | List of [abstracts](#abstract) |
| `arguments`      | ⬤ | ◯ | List of [arguments](#argument) |
| `basics`         | ⬤ | ◯ | List of [basics](#basic) |
| `fields`         | ⬤ | ◯ | List of [Fields](#field) |
| `interfaceDecls` | ⬤ | ◯ | List of [interface declarations](#interface-declaration) |
| `interfaceDescs` | ⬤ | ◯ | List of [interface descriptions](#interface-description) |
| `interfaceInsts` | ⬤ | ◯ | List of [interface instances](#interface-instance) |
| `language`       | ◯ | ◯ | A string for the source code language, e.g. `go`, `java`. |
| `locs`           | ◯ | ◯ | [locations](#locations) |
| `methods`        | ⬤ | ◯ | List of [methods](#method) |
| `methodInsts`    | ⬤ | ◯ | List of [method instances](#method-instance) |
| `metrics`        | ⬤ | ◯ | List of [metrics](#metrics) |
| `objects`        | ⬤ | ◯ | List of [objects](#object) |
| `objectInsts`    | ⬤ | ◯ | List of [object instances](#object-instance) |
| `packages`       | ⬤ | ◯ | List of [packages](#package) |
| `selections`     | ⬤ | ◯ | List of [selections](#selection) |
| `signatures`     | ⬤ | ◯ | List of [signatures](#signature) |
| `structDescs`    | ⬤ | ◯ | List of [structure descriptions](#structure-description) |
| `typeParams`     | ⬤ | ◯ | List of [type parameters](#type-parameter) |
| `values`         | ⬤ | ◯ | List of [values](#value) |

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

```Java
public interface Named {
  void Foo();
  int Bar(String s);
}
```

| Name        | Optional | Extra | Description |
|:------------|:--------:|:-----:|:------------|
| `index`     | ◯ | ⬤ | The [index](#indices) of this abstract in the project's `abstracts` list. |
| `kind`      | ◯ | ⬤ | `abstract` |
| `name`      | ◯ | ◯ | The string name for the abstract. |
| `signature` | ◯ | ◯ | [Index](#indices) for the [signature](#signature). |
| `vis`       | ◯ | ⬤ | A string of the scope modifiers, like "public", "exported", or "private". |

### Argument

An argument (`argument`) is an optionally named parameter or result.

For example, the `Foo` method below has four arguments: `x string`,
`y int`, `ok bool`, `err error`.

```Go
func Foo(x string, y int)(ok bool, err error) { }
```

The next method `Bar` has a matching signature to `Foo`
but defined with nameless arguments.
In many cases the name of the argument is ignored since interface abstracts
may have different argument names from the method that it matches with.

```Go
func Bar(string, int)(bool, error) { }
```

Any repeat type argument in Go is expanded.
For example, `Baz` below has two arguments, `x float64` and `y float64`.

```Go
func Baz(x, y float64) { }
```

Java is simpler since it only allows one return value that may not be named.
In the following, there are 3 arguments, an unnamed `boolean`, `String x`, and `int y`.

```Java
boolean Foo(String x, int y) { }
```

| Name    | Optional | Extra | Description |
|:--------|:--------:|:-----:|:------------|
| `index` | ◯ | ⬤ | The [index](#indices) of this argument in the project's `arguments` list. |
| `kind`  | ◯ | ⬤ | `argument` |
| `name`  | ⬤ | ◯ | The optional string name for the argument. |
| `type`  | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |

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

The Java types like `Integer`, `Character`, etc will use the basic tpe
`int`, `char` and so on. `String` will be treated as a basic.

If the `index` and `kind` are not being outputted, then the basic type may
be a string containing the name instead of an object.

| Name    | Optional | Extra | Description |
|:--------|:--------:|:-----:|:------------|
| `index` | ◯ | ⬤ | The [index](#indices) of this basic in the project's `basics` list. |
| `kind`  | ◯ | ⬤ | `basic` |
| `name`  | ◯ | ◯ | The string name for the basic type, e.g. `int`, `bool`, `string`. |

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

```Java
public class Named {
  public int x;
  private String y;
}
```

Any repeat type fields in Go is expanded, e.g. `struct { x, y float64 }` will
have two fields, `x float64` and `y float64`.

Any embedded fields will be named using the name of the base type,
e.g. `struct { *Foo[int] }` will be the same as `struct { Foo *Foo[int] }`.
Unnamed types are not allowed to be embedded into a struct.

| Name       | Optional | Extra | Description |
|:-----------|:--------:|:-----:|:------------|
| `embedded` | ◯ | ⬤ | True if the field is from an embedded struct. |
| `index`    | ◯ | ⬤ | The [index](#indices) of this field in the project's `fields` list. |
| `kind`     | ◯ | ⬤ | `field` |
| `name`     | ◯ | ◯ | The string name for the field. |
| `type`     | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |
| `vis`      | ◯ | ⬤ | A string of the scope modifiers, like "public", "exported", or "private". |

### Interface Declaration

An interface declaration (`interfaceDecl`) is a named definition of an
interface. For example, the following is the declaration of `Foo` with
a [type parameter](#type-parameter) `T` and an [abstract](#abstract) `Bar`.

```Go
type Foo[T any] interface {
  Bar(value T) string
}

type Baz interface {
  Error() string
}
```

```Java
public interface Foo<T> {
  String bar(T value);
}

public interface Baz {
  String error();
}
```

| Name         | Optional | Extra | Description |
|:-------------|:--------:|:-----:|:------------|
| `index`      | ◯ | ⬤ | The [index](#indices) of this interface declaration in the project's `interfaceDecls` list. |
| `instances`  | ⬤ | ◯ | List of [indices](#indices) to [interface instances](#interface-instance). |
| `interface`  | ◯ | ◯ | The [index](#indices) to the declared [interface](#interface-description) type. |
| `kind`       | ◯ | ⬤ | `interfaceDecl` |
| `loc`        | ⬤ | ◯ | The [location](#locations) offset. |
| `name`       | ◯ | ◯ | The name of the declared interface. |
| `package`    | ◯ | ◯ | The [index](#indices) to the [package](#package) this declaration is declared in. |
| `typeParams` | ⬤ | ◯ | List of [indices](#indices) to [type parameters](#type-parameter) if this interface is generic. |
| `vis`        | ◯ | ⬤ | A string of the scope modifiers, like "public", "exported", or "private". |

### Interface Description

An interface description (`interfaceDesc`) describes the type of an interface.
This interface type may be the type for an interface declaration, interface
instances, and interface literal.

```Go
interface { String() string }

interface { int | ~string }
```

```Java
public interface Named {
  Strings toString();
}
```

| Name        | Optional | Extra | Description |
|:------------|:--------:|:-----:|:------------|
| `abstracts` | ⬤ | ◯ | List of [indices](#indices) to [abstracts](#abstract). |
| `approx`    | ⬤ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for approximate constraints. |
| `exact`     | ⬤ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for exact constraints. |
| `hint`      | ◯ | ⬤ | A string indicating if the interface is a stand-in for a type, e.g. `pointer`, `chan`, `list` |
| `index`     | ◯ | ⬤ | The [index](#indices) of this interface in the projects' `interfaceDescs` list. |
| `inherits`  | ⬤ | ◯ | List of [indices](#indices) to inherited [interfaces](#interface-description). |
| `kind`      | ◯ | ⬤ | `interfaceDesc` |
| `pin`       | ⬤ | ◯ | The [key](#keys) to the [object](#object), [interface](#interface-declaration) declaration, or [package](#package) this interface is pinned to. |

### Interface Instance

An interface instance (`interfaceInst`) is an instantiation of a generic
interface declaration.
The instance types are the type arguments used in the type parameters.
The instance types may be type parameters as well as a fully realized type.
For example, `Foo[T any]` with instance type `bool` will create `Foo[bool]`.

```Go
type Foo[T any] interface { Value() T }

type Foo[bool] interface { Value() bool }
```

Some instances are still generic instead of a fully realized type.
For example, `Foo[T any]` with type parameter `S` will create `Foo[S]`.
Any type parameter being used in an instance will always be the same
as the original or a subtype. If they are the same, then at least one of the
other type parameters in the instance must be a subtype.

```Go
func Bar[S int | string]() Foo[S] { }
```

| Name            | Optional | Extra | Description |
|:----------------|:--------:|:-----:|:------------|
| `generic`       | ◯ | ◯ | The [index](#indices) of the generic [interface declaration](#interface-declaration) this is an instance of. |
| `index`         | ◯ | ⬤ | The [index](#indices) of this interface in the projects' `interfaceInsts` list. |
| `instanceTypes` | ◯ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for type arguments. |
| `kind`          | ◯ | ⬤ | `interfaceInst` |
| `resolved`      | ◯ | ◯ | The [index](#indices) to the resolved [interface description](#interface-description) this instance defines. |

### Locations

The locations (`locations`) is a map from offsets to file names.
The offsets start at 1 and accumulate with the line count of each file.
For example, given 3 files `Apple.go`, `Banana.go`, and `Carrot.go` that
have 42, 55, and 38 lines respectively, the following locations map
could be created:

```JSON
{
  "1": "Apple.go",
  "43": "Banana.go",
  "98": "Carrot.go",
}
```

The `loc` item in some constructs are the offset into this map.
The `loc` offset may be between two file offsets because the file offsets
indicate line 1 in the file.

To determine the file and lines, find the highest number offset
that is less than or equal to the location. That will give you the file
and the offset to that file's first line. Then $loc - firstLine + 1 = line$.

For the above example, given $104$ we can determine it is in "Carrot.go".
"Carrot.go"'s first line is at the offset $98$, so we calculate the line with
$104 - 98 + 1 = 7$, to determine the `loc` is on the 7th line of "Carrot.go".

| loc |    file   | lineNum |
|----:|----------:|--------:|
|   1 |  Apple.go |    1    |
|   2 |  Apple.go |    2    |
|  41 |  Apple.go |   41    |
|  42 |  Apple.go |   42    |
|  43 | Banana.go |    1    |
|  44 | Banana.go |    2    |
|  96 | Banana.go |   53    |
|  97 | Banana.go |   54    |
|  98 | Carrot.go |    1    |
|  99 | Carrot.go |    2    |
| 104 | Carrot.go |    7    |

### Method

A method declaration (`method`) us a named definition of a function not
attached to an object or a method with a receiver object.
The method may be generic and have instances attached to it.

```Go
func Foo[T any](value T) { }

func (b *Bar[T]) Foo(value T) { }

func Baz(value int)(ok bool, err error) { }
```

```Java
public class Bar<T> {

  public void Foo(T value);
  
  public <T2> void Foo(T2 value);
}
```

| Name         | Optional | Extra | Description |
|:-------------|:--------:|:-----:|:------------|
| `index`      | ◯ | ⬤ | The [index](#indices) of this method in the projects' `methods` list. |
| `instances`  | ⬤ | ◯ | List of [indices](#indices) to [method instances](#method-instance). |
| `kind`       | ◯ | ⬤ | `method` |
| `loc`        | ⬤ | ◯ | The [location](#locations) offset. |
| `name`       | ◯ | ◯ | The name of the declared method. |
| `metrics`    | ⬤ | ◯ | The [index](#indices) of the [metrics](#metrics) for this method. |
| `package`    | ◯ | ◯ | The [index](#indices) of the [package](#package) this method is declared in. |
| `ptrRecv`    | ◯ | ⬤ | A boolean indicating if the method had a Go's pointer receiver. |
| `receiver`   | ⬤ | ◯ | The [index](#indices) of the [object](#object) that is the receiver if there is one. |
| `recvName`   | ◯ | ⬤ | The name given to the receiver. |
| `signature`  | ◯ | ◯ | The [index](#indices) of the [signature](#signature) for this method. |
| `typeParams` | ⬤ | ◯ | List of [indices](#indices) to [type parameters](#type-parameter) if this method is generic. |
| `vis`        | ◯ | ⬤ | A string of the scope modifiers, like "public", "exported", or "private". |

### Method Instance

A method instance (`methodInst`) is an instantiation of a generic
[method declaration](#method).
The instance types are the type arguments used in the type parameters.
The instance types may be type parameters as well as a fully realized type.

For example, `Foo[T any](value T)` with instance type `bool`
will create, `Foo[bool](value bool)`. The instance was created by
calling `Foo` inside a method with the type argument given.
The type argument may be another tye parameter, meaning the method is
still generic, or concrete to fully realize the method.

```Go
func Foo[T any](value T) { }

func Foo[bool]() bool { }
```

Instances of a receiver cascade down to create instances of the methods with
that receiver. Below `Bar[T]` gets fully realized with `int`. The generic
method `Bar[T].Foo` then has an instance for `Bar[int].Foo`.

```Go
type Bar[T any] struct { }

func (b *Bar[T]) Foo(value T) { }

func (b *Bar[int]) Foo(value int) { }
```

| Name            | Optional | Extra | Description |
|:----------------|:--------:|:-----:|:------------|
| `generic`       | ◯ | ◯ | The [index](#indices) of the generic [method](#method) this is an instance of. |
| `index`         | ◯ | ⬤ | The [index](#indices) of this method instance in the projects' `methodInsts` list. |
| `instanceTypes` | ◯ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for type arguments. |
| `kind`          | ◯ | ⬤ | `methodInst` |
| `receiver`      | ⬤ | ◯ | The [index](#indices) of the [object instance](#object-instance) for the receiver of this method, if it has one. |
| `resolved`      | ◯ | ◯ | The [index](#indices) to the resolved [signature](#signature) this instance defines. |

### Metrics

A metrics (`metrics`) is measurements done to a set of expressions.
The body of a [method](#method) and the initializer for a [value](#value)
contain expressions that are measured. These measurements are used in
technical debt analysis.

| Name         | Optional | Extra | Description |
|:-------------|:--------:|:-----:|:------------|
| `codeCount`  | ⬤ | ◯ | The number of lines in the method that are not comments or empty. |
| `complexity` | ⬤ | ◯ | The cyclomatic complexity of the method. |
| `getter`     | ⬤ | ◯ | True indicates the method is a getter pattern. |
| `indents`    | ⬤ | ◯ | The indent complexity of the method. |
| `index`      | ◯ | ⬤ | The [index](#indices) of this metrics in the projects' `metrics` list. |
| `invokes`    | ⬤ | ◯ | List of [keys](#keys) to methods (declaration or instance) that were invoked in the method. |
| `kind`       | ◯ | ⬤ | `metrics` |
| `lineCount`  | ⬤ | ◯ | The number of lines in the method. |
| `loc`        | ◯ | ◯ | The [location](#locations) offset. |
| `reads`      | ⬤ | ◯ | List of [keys](#keys) to types that were read from in the method. |
| `setter`     | ⬤ | ◯ | True indicates the method is a setter pattern. |
| `sideEffect` | ◯ | ⬤ | True indicates this method directly has side effects without checking invoked method. |
| `writes`     | ⬤ | ◯ | List of [keys](#keys) to types that were written to in the method. |

### Object

An object declaration (`object`) is a collection of data via a
[structure](#structure-description) with zero or more [methods](#method),
like a "class" in Java.

The following code defines the object `Foo` with the structure
`struct { x, y int }`  and a method `Bar`.

```Go
type Foo struct { x, y int }

func (f Foo) Bar() { }
```

```Java
public class Foo {
  private int x, y;

  public Bar() { }
}
```

The object may be generic if it has type parameters.
Objects

```Go
type Bar[T any] struct { x, y T }

func (b Bar[T]) Bar(x, y T) { }
```

```Java
public class Foo<T> {
  private T x, y;
  
  public Bar(T x, T y) { }
}
```

If the named type in Go is defined without a struct, e.g. `Baz`,
the abstractor will have to pack the type into struct.

```Go
type Baz int

type Baz struct { $data int }
```

| Name         | Optional | Extra | Description |
|:-------------|:--------:|:-----:|:------------|
| `data`       | ◯ | ◯ | The [index](#indices) of the [structure description](#structure-description). |
| `index`      | ◯ | ⬤ | The [index](#indices) of this object in the projects' `objects` list. |
| `instances`  | ⬤ | ◯ | List of [indices](#indices) to [object instances](#object-instance). |
| `kind`       | ◯ | ⬤ | `object` |
| `loc`        | ⬤ | ◯ | The [location](#locations) offset. |
| `methods`    | ⬤ | ◯ | List of [indices](#indices) to [methods](#method) that have this object as a receiver. |
| `name`       | ◯ | ◯ | The name of the declared object. |
| `package`    | ◯ | ◯ | The [index](#indices) of the [package](#package) this object is declared in. |
| `typeParams` | ⬤ | ◯ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
| `interface`  | ◯ | ◯ | The [index](#indices) to the [interface description](#interface-description) that this object matches with. |
| `vis`        | ◯ | ⬤ | A string of the scope modifiers, like "public", "exported", or "private". |

### Object Instance

A object instance (`objectInst`) is an instantiation of a generic
[object declaration](#object).
The instance types are the type arguments used in the type parameters.
The instance types may be type parameters as well as a fully realized type.
For example, `Foo[T any]` with instance type `bool` will create `Foo[bool]`:

```Go
type Foo[T any] struct { value T }

type Foo[bool] struct { value bool }
```

| Name            | Optional | Extra | Description |
|:----------------|:--------:|:-----:|:------------|
| `generic`       | ◯ | ◯ | The [index](#indices) of the generic [object](#object) this is an instance of. |
| `index`         | ◯ | ⬤ | The [index](#indices) of this object instance in the projects' `objectInsts` list. |
| `instanceTypes` | ◯ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for type arguments. |
| `kind`          | ◯ | ⬤ | `objectInst` |
| `methods`       | ◯ | ◯ | List of [indices](#indices) of the [method instances](#method) for this instance. |
| `resData`       | ◯ | ◯ | The [index](#indices) to the resolved [structure description](#structure-description) this instance defines. |
| `resInterface`  | ◯ | ◯ | The [index](#indices) to the resolved [interface description](#interface-description) this instance defines. |

### Package

A package (`package`) is a collection of code usually in several files that
typically are all part of a related library.

| Name         | Optional | Extra | Description |
|:-------------|:--------:|:-----:|:------------|
| `imports`    | ⬤ | ◯ | List of [indices](#indices) of [packages](#package) that this package depends on. |
| `index`      | ◯ | ⬤ | The [index](#indices) of this package in the projects' `packages` list. |
| `interfaces` | ⬤ | ◯ | List of [indices](#indices) of [interfaces](#interface-declaration) declared in this package. |
| `kind`       | ◯ | ⬤ | `package` |
| `methods`    | ⬤ | ◯ | List of [indices](#indices) of [methods](#method) declared in this package. |
| `name`       | ◯ | ◯ | The name of the package. |
| `objects`    | ⬤ | ◯ | List of [indices](#indices) of [object](#object) declared in this package. |
| `path`       | ◯ | ◯ | The path to this package. |
| `values`     | ⬤ | ◯ | List of [indices](#indices) of [values](#value) declared in this package. |

### Selection

A selection (`selection`) represents a field, method, or abstract being
accessed. A selection is typically caused by a `dot` in both Java and Go,
e.g. `f.value` is `value` selected from `f` of type `Foo`.

```Go
type Foo { value int }

func Bar(f Foo) int { return f.value }
```

Selections are used in [metrics](#metrics) to indicate higher detailed
information than simply specifying the type of the selected field, method,
or abstract.

| Name     | Optional | Extra | Description |
|:---------|:--------:|:-----:|:------------|
| `index`  | ◯ | ⬤ | The [index](#indices) of this selection in the projects' `selections` list. |
| `kind`   | ◯ | ⬤ | `selection` |
| `name`   | ◯ | ◯ | The name of the field, method, or abstract that is selected. The `f` in `x.f`. |
| `origin` | ◯ | ◯ | The [key](#keys) to the [construct](#constructs) that is selected from. The `x` in `x.f`. |

### Signature

A signature (`signature`) represents the shape of a method's input and output
[arguments](#argument). It can be used in interface abstracts, as function pointers,
delegate types, methods of an object, or a function.
For example `func(x int) string`. The names of arguments are ignored in many
cases since the signature types determine if two signatures are the same even
if the names are different.

In the following there are three abstracts in an interface. `Foo` and `Bar` have
the same signature `func(int) string` and `Baz` has the signature
`func() (int, bool)`.

```Go
interface {
  Foo(x int) string
  Bar(y int) (name string)
  Baz() (value int, okay bool)
}
```

| Name       | Optional | Extra | Description |
|:-----------|:--------:|:-----:|:------------|
| `index`    | ◯ | ⬤ | The [index](#indices) of this signature in the projects' `signatures` list. |
| `kind`     | ◯ | ⬤ | `signature` |
| `variadic` | ⬤ | ◯ | True indicates the last parameter is a variable length parameter. |
| `params`   | ⬤ | ◯ | List of [indices](#indices) of [arguments](#argument) for input parameters. |
| `results`  | ⬤ | ◯ | List of [indices](#indices) of [arguments](#argument) for output results. |

### Structure Description

A structure description (`structDesc`) describes a collection of values,
called [fields](#field), like a record, tuple, or class properties.
The following struct contains three fields:

```Go
struct {
  name string
  age  int
  id   uint64
}
```

```Java
public class Named {
  private String name;
  private int age;
  private long id;
}
```

| Name        | Optional | Extra | Description |
|:------------|:--------:|:-----:|:------------|
| `fields`    | ⬤ | ◯ | List of [indices](#indices) of [fields](#field) in this structure. |
| `index`     | ◯ | ⬤ | The [index](#indices) of this structure in the projects' `structDescs` list. |
| `kind`      | ◯ | ⬤ | `structDesc` |
| `synthetic` | ◯ | ⬤ | A boolean indicating if the abstractor had to create this structure for non-struct types, e.g. `type foo int` |

### Type Parameter

A type parameter (`typeParam`) is a type defined with a generic
[object](#object) or [method](#method). These are named parameters that define
custom types for the declaration. Different instances define the types,
instance types or type arguments, that realize these parameters.

For example `T any` is a type parameter in the following code. The `value`
is of type `T` meaning it will become the type used as an argument into `T`.

```Go
type Foo[T any] struct { value T }
```

```Java
public class Foo<T> {
  private T value;
}
```

| Name    | Optional | Extra | Description |
|:--------|:--------:|:-----:|:------------|
| `index` | ◯ | ⬤ | The [index](#indices) of this type parameter in the projects' `typeParams` list. |
| `kind`  | ◯ | ⬤ | `typeParam` |
| `name`  | ◯ | ◯ | The name of the type parameter. |
| `type`  | ◯ | ◯ | The [key](#keys) for the [type](#type-descriptions) of this type parameter. |

### Value

A value declaration (`value`) is a package level variable outside of any
declared object. A value may be constant and may be initialized by an
expression.

```Go
var X = 10

var Y, Z = func() (int, int) {
  return X*2, X+3
}()
```

| Name      | Optional | Extra | Description |
|:----------|:--------:|:-----:|:------------|
| `const`   | ⬤ | ◯ | True indicates the value is constant after initialization |
| `index`   | ◯ | ⬤ | The [index](#indices) of this value in the projects' `values` list. |
| `kind`    | ◯ | ⬤ | `value` |
| `loc`     | ⬤ | ◯ | The [location](#locations) offset. |
| `metrics` | ⬤ | ◯ | The [index](#indices) of the [metrics](#metrics) for the initializer expression. |
| `name`    | ◯ | ◯ | The name of the value. |
| `package` | ◯ | ◯ | The [index](#indices) of the [package](#package) this value is described in. |
| `type`    | ◯ | ◯ | The [key](#keys) for the [type](#type-descriptions) of this value. |
| `vis`     | ◯ | ⬤ | A string of the scope modifiers, like "public", "exported", or "private". |
