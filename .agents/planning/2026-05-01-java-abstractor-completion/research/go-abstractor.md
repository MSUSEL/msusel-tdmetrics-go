# Research: Go Abstractor Patterns

## Overview

The Go abstractor follows a two-phase architecture: **abstraction** (AST walk)
then **resolution** (post-processing pipeline). The Java abstractor should
follow the same pattern.

## 1. Main Flow

```
Abstract(cfg) -> project
  1. Create querier (wraps go/packages data)
  2. Create locs (source location tracking)
  3. Create project (output model)
  4. Create baker (factory for built-in/synthetic types)
  5. Iterate all packages -> abstractPackage()
     -> Iterate all files -> abstractFile()
        -> Iterate all declarations (types, values, functions)
  6. Pass to resolver.Resolve() for post-processing
```

State carried during abstraction:
- `curPkg` — current package being processed
- `curNest` — current nesting context (for types inside functions)
- `implicitTypes` — type parameters from enclosing function
- `tpReplacer` — maps receiver type params back to the original type's params
- `typeCache` — deduplication cache for converted types

## 2. External/Library Types (Baker + Innate)

### Innate Names

Constants prefixed with `$` represent synthetic operations:
- `$builtin` — synthetic builtin package
- `$data` — synthetic field wrapping underlying type in named types
- `$deref`, `$get`, `$set`, `$len` — operations on pointers, slices, maps
- `$recv`, `$send` — channel operations
- `$equal` — comparable types
- `$real`, `$imag` — complex numbers

### Baker

A lazy, memoized factory creating synthetic type declarations:

| Go concept | Baker creates | Shape |
|---|---|---|
| `any` | Empty `InterfaceDesc` | `interface{}` |
| `*T` | `Pointer[T any]` | `interface{ $deref() T }` |
| `[]T` / `[n]T` | `List[T any]` | `interface{ $len() int; $get(int) T; $set(int, T) }` |
| `map[K]V` | `Map[K comparable, V any]` | `interface{ $len() int; $get(K) (V,bool); $set(K, V) }` |
| `chan T` | `Chan[T any]` | `interface{ $len() int; $recv() (T,bool); $send(T) }` |
| `complex64/128` | Non-generic interface | `interface{ $real() float32; $imag() float32 }` |
| `error` | Non-generic interface | `interface{ Error() string }` |
| `comparable` | Non-generic interface | `interface{ $equal(other any) bool }` |

### Java Equivalent

The Java Baker already handles arrays with `$Array[T]` using `$len`, `$get`,
`$set`. Additional baking needed for:
- `Object` as empty interface (already done as `baker.objectDesc()`)
- Boxing types — map `Integer`→`int`, `Character`→`char`, etc.
- `String` as a basic type
- Potentially `Iterable`, `Comparable` as named interfaces if encountered

## 3. Generics / Instantiations

Three parallel instantiation constructs:

### ObjectInst
- `generic` — back-pointer to `Object`
- `resolvedData` — struct with concrete types
- `resolvedInterface` — interface with concrete types
- `implicitTypes` — from enclosing function (for nested types)
- `instanceTypes` — the actual type arguments
- Registered on generic via `Generic.AddInstance(inst)`

### MethodInst
- `generic` — back-pointer to `Method`
- `resolved` — concrete `Signature`
- `instanceTypes` — type arguments
- `metrics` — re-analyzed with concrete types
- `receiver` — optional `ObjectInst`

### InterfaceInst
- `generic` — back-pointer to `InterfaceDecl`
- `resolved` — concrete `InterfaceDesc`
- `implicitTypes` — from enclosing function
- `instanceTypes` — type arguments

### Expansion Loop (Fixed-Point)

The resolver runs until stable:
1. `fillOutAllMetrics()` — re-analyze AST for each new `MethodInst`
2. `expandAllInstantiations()` — for each object, ensure every method has a
   `MethodInst` for every `ObjectInst`
3. `fillOutAllPointerReceivers()` — pointer interface declarations
4. `expandAllNestedTypes()` — instances of nested types per `MethodInst`

**Java parallel:** Same expansion needed. When `List<String>` is instantiated,
all its methods need instantiated signatures with `String` replacing `T`.

## 4. Nested Types

### Detection
After processing a function declaration, `abstractNestedTypes` inspects the
body for `*ast.TypeSpec` declarations.

### Nesting Context
1. A `TempDeclRef` placeholder is created for the method
2. `curNest` is set to this placeholder
3. `implicitTypes` set to the function's type parameters
4. Types found inside body are created with `Nest: ab.curNest`
5. After method is created, placeholder is resolved

### Object Tracking
- `nest` field stores the enclosing method
- `ImplicitTypeParams()` returns the nesting method's type params
- Comparisons include nest in sort key (`pkg.foo:Bar` != `pkg.Bar`)

### Nested Expansion
When a generic method has instances, nested types also need instances with
the implicit types filled in.

**Java parallel:** Java inner classes carry the outer class's type params.
Named nested classes → separate objects with `nest` set to the enclosing
method or object. Anonymous classes and lambdas → part of enclosing method
(per requirements).

## 5. Metrics Computation

Three independent analyzers:

### Complexity
- Starts at 1
- +1 for: `if`, `for`, `range`, `go`, non-default `case`, `&&`, `||`
- Tracks: `lineCount`, `codeCount`, `indents` (column offset sum)

### Accessor (Getter/Setter Detection)
- **Getter**: No params, one result, body is single return of "simple fetch"
- **Setter**: 0-1 params, no results, body is single assignment of simple fetches

### Usages (Reads/Writes/Invokes)
- Maintains a `pending` construct flushed as read/write/invoke based on context
- **Reads**: Constructs evaluated but not assigned to or called
- **Writes**: LHS of assignments, `++/--`, definition targets
- **Invokes**: Call expression targets (excluding type conversions and builtins)
- **SideEffect**: Writing to globals or calling print/println
- Uses `Selection` constructs for `obj.Field` / `obj.Method` tracking
- Uses `TempDeclRef`/`TempReference` for forward references (resolved later)

### Re-analysis for Instances
`MethodInst` stores the AST node + `tpReplacer` so metrics can be re-analyzed
with concrete types for each instantiation.

**Java parallel:** All three translate directly. The re-analysis pattern is
important for generic method instantiations.

## 6. Resolver Post-Processing Pipeline

```
1. Imports           — resolve import paths to Package constructs
2. Receivers         — attach methods to receiver objects (Go-specific)
3. Fixed-point loop:
   a. References     — resolve TempReference/TempDeclRef to constructs
   b. Instantiations — expand all generic instances
   (loop until no changes)
4. References (final) — second pass for remaining
5. RemoveDuplicates  — dedup constructs from reference resolution
6. GenerateInterfaces — create InterfaceDesc for each object from its methods
7. Inheritance       — compute inheritance forest (duck typing)
8. DCE              — dead code elimination from entry points
9. Locations         — flag used source locations
10. Indices          — assign output indices
```

### GenerateInterfaces
- For each `Object`, creates `InterfaceDesc` from its non-pointer methods
- For each `ObjectInst`, creates resolved interface from instance's methods
- Extends pointer interfaces with pointee methods

### Inheritance
- Generic forest algorithm using structural subtyping
- `Implements(other)` tests if interface A is subtype of B
- Transitively prunes redundant edges

### DCE
- Work-list algorithm from: main, init, test functions, exported decls
- Propagates liveness through all construct references

**Java parallel:** Steps 1, 3-5, 9-10 map directly. Step 2 is unnecessary
(Java methods are inside their class). Step 6 needs adaptation — Java classes
have explicit interfaces but we still synthesize the effective interface. Steps
7-8 apply with explicit `extends`/`implements` as the base.

## 7. Package Imports

- Each package stores: `path`, `name`, `importPaths` (string list)
- Contains sorted sets of: `imports`, `interfaces`, `methods`, `objects`, `values`
- `EntryPoint()` returns true for the main package
- Import resolution: iterate `importPaths`, call `proj.FindPackageByPath()`

**Java parallel:** Java packages map 1:1. Package dependencies should be
derived from resolved type references rather than `import` statements.

## Key Alignment Patterns for Java

| Pattern | Go | Java Equivalent |
|---|---|---|
| Two-phase (abstract + resolve) | ✅ | Same needed |
| Baker / synthetic types | Pointer, List, Map, Chan | Array, Object, boxing |
| `$`-prefixed innate names | ✅ | Same convention |
| Generic instantiation tracking | ObjectInst, MethodInst, InterfaceInst | Same three |
| Fixed-point expansion loop | References + Instantiations | Same needed |
| Nested types via NestType | `curNest` + `implicitTypes` | Inner/local classes |
| Metrics: complexity, accessor, usages | Three analyzers | Same three |
| Metrics re-analysis for instances | Store AST, re-analyze | Same needed |
| Structural interface generation | From object methods | Synthesize effective interface |
| TempReference / TempDeclRef | Lazy forward references | Java's `Ref<T>` already similar |
| DCE from entry points | Work-list algorithm | Same algorithm |
