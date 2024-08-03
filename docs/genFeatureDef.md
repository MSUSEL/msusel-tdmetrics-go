# Generalized Feature Definition

The generalized feature definition is a
[JSON](https://www.json.org/json-en.html) file with specific fields to
define the methods and data of an application. This definition needs to be
flexible enough to handle procedural and object-oriented languages with some
adjustment specific to the language. For example, Go uses duck-typing
but the definition requires a list of specific implementation so the Go
abstractor must perform steps to predetermine which types would duck-type
and define that via implementations.

:notice: Out-of-date

```mermaid
classDiagram
  direction LR

  class Abstraction {
    language: string
    packages: []Package
    types:    Types
  }
  Abstraction --> "*" Package
  Abstraction --> "1" Types

  class Package {
    path:    string
    name:    string
    imports: []Package
    types:   []TypeDef
    values:  []ValueDef
    methods: []Method
  }
  Package --> "*" Package
  Package --> "*" TypeDef
  Package --> "*" ValueDef
  Package --> "*" Method

  class Types {
    basics:     []Basic
    interfaces: []Interface
    named:      []Named
    signatures: []Signature
    solids:     []Solid
    structs:    []Struct
    unions:     []Union
  }
  Types --> "*" Basic
  Types --> "*" Interface
  Types --> "*" Named
  Types --> "*" Signature
  Types --> "*" Solid
  Types --> "*" Struct
  Types --> "*" Union

  class TypeDef {
    name:       string
    type:       TypeDesc
    methods:    []Method
    typeParams: []Named
    interface:  Interface
  }
  TypeDef --> "1" TypeDesc
  TypeDef --> "*" Method
  TypeDef --> "*" Named
  TypeDef --> "1" Interface

  class ValueDef {
    name:  string
    const: bool
    type:  TypeDesc
  }
  ValueDef --> "1" TypeDesc

  class Method {
    name:      string
    signature: TypeDesc
    metrics:   Metrics
  }
  Method --> "1" TypeDesc
  Method --> "1" Metrics

  class TypeDesc {
    <<interface>>
  }

  class Basic {
    string
  }
  Basic ..o TypeDesc

  class Interface {
    typeParams: []Named
    inherits:   []Interface
    union:      Union
    methods:    map[string, TypeDesc]
  }
  Interface ..o TypeDesc
  Interface --> "*" Named
  Interface --> "*" Interface
  Interface --> "1" Union
  Interface --> "*" TypeDesc

  class Named {
    name: string
    type: TypeDesc
  }
  Named ..o TypeDesc
  Named --> "1" TypeDesc
  
  class Signature {
    variadic:   bool
    params:     []Named
    typeParams: []Named
    return:     TypeDesc
  }
  Signature ..o TypeDesc
  Signature --> TypeDesc
  Signature --> Named
  
  class Solid {
    target:     TypeDesc
    typeParams: []TypeDesc
  }
  Solid ..o TypeDesc
  Solid --> TypeDesc
  
  class Struct {
    fields: []Named
  }
  Struct ..o TypeDesc
  Struct --> Named
  
  class Union {
    exact:  []TypeDesc
    approx: []TypeDesc
  }
  Union ..o TypeDesc
  Union --> TypeDesc
```
