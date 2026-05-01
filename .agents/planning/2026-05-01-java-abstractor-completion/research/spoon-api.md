# Research: Spoon API Patterns

## Overview

Spoon (v11.2.0) is a Java AST library designed for program analysis and
transformation. Its meta-model mirrors Java's structure with three parts:
**structural** (declarations), **code** (executable statements/expressions),
and **references** (pointers to elements). All types are prefixed with "Ct"
(compile-time).

## Type Hierarchy (Structural Elements)

```
CtElement
├── CtType<T>                     (base for all type declarations)
│   ├── CtClass<T>                (classes, also base for enums/records)
│   │   ├── CtEnum<T>            (enumerations)
│   │   └── CtRecord             (Java 16+ records)
│   ├── CtInterface<T>           (interfaces)
│   ├── CtAnnotationType<A>      (annotation type declarations)
│   └── CtTypeParameter          (generic type parameters like <T>)
├── CtPackage
├── CtMethod<T>
├── CtConstructor<T>
├── CtField<T>
├── CtParameter<T>
├── CtAnonymousExecutable        (static/instance initializer blocks)
└── CtCompilationUnit
```

### Key relationships:
- `CtClass` is supertype of `CtEnum` and `CtRecord`.
- `CtInterface` is separate from `CtClass`.
- `CtAnnotationType` is a `CtType` but NOT a `CtClass` or `CtInterface`.
- `CtTypeParameter` is a `CtType` (represents `<T extends Foo>`).

## References vs Declarations

| Method | Returns | On external types |
|--------|---------|-------------------|
| `tr.getDeclaration()` | `CtType` or `null` | Returns `null` |
| `tr.getTypeDeclaration()` | `CtType` (may be shadow) | Builds from reflection |

**Critical for the abstractor**: The current code uses `tr.getDeclaration()`
which returns `null` for external types (java.lang, java.util, etc.). This
causes the abstractor to fall through to `this.proj.baker.objectDesc()`.
Consider using `tr.getTypeDeclaration()` when more information is needed,
but be aware that shadow types have limited information.

### Shadow Types

A "shadow" element is built via Java reflection when the source isn't available.
Check with `type.isShadow()`. Shadow types have:
- Declaration structure (fields, methods, superclass, interfaces)
- No source positions
- No method bodies
- Limited annotation information

## CtTypeReference Query Methods

These boolean methods on `CtTypeReference` classify a type:

| Method | What it matches |
|--------|----------------|
| `isPrimitive()` | `int`, `boolean`, `void`, etc. |
| `isClass()` | classes (including enums, records) |
| `isInterface()` | interfaces |
| `isEnum()` | enumerations |
| `isGenerics()` | type parameters (`T`, `E`, etc.) |
| `isArray()` | array types (`int[]`, `String[]`) |
| `isAnnotationType()` | annotation types (`@Override`) |
| `isAnonymous()` | anonymous classes |
| `isLocalType()` | local classes defined inside methods |
| `isParameterized()` | types with actual type arguments (`List<String>`) |
| `isShadow()` | shadow types (from reflection, not source) |

### Important: these are NOT mutually exclusive
- An enum is also `isClass() == true` (because `CtEnum extends CtClass`).
- A parameterized type is also `isClass() == true`.
- Check order matters in the abstractor's `addTypeDesc()`.

## Wildcard Types

`CtWildcardReference` represents `?` in generics:
- `?` — unbounded wildcard
- `? extends Foo` — upper-bounded (`isUpper() == true`)
- `? super Bar` — lower-bounded (`isUpper() == false`)
- `getBoundingType()` returns the bound type (may be `CtIntersectionTypeReference`)

Currently not handled by the abstractor. Would map to an interface description
or type parameter in the generalized feature definition.

## Anonymous Classes

`CtNewClass` represents anonymous class creation:
```java
Runnable r = new Runnable() {
    @Override public void run() { ... }
};
```
- `getAnonymousClass()` returns the `CtClass` for the anonymous type.
- The type reference will have `isAnonymous() == true`.
- `getSimpleName()` returns a numeric name like `1`, `2`, etc.

## Lambda Expressions

`CtLambda` represents lambda expressions:
```java
list.stream().map(x -> x.toString());
```
- Extends both `CtExpression` and `CtExecutable`.
- `getBody()` returns the block body (if block lambda).
- `getExpression()` returns the expression body (if expression lambda).
- Has parameters but generates a synthetic method.
- The type of a lambda is a functional interface.

## Code Elements Relevant to Metrics

### Statements
| Spoon Type | Java Construct | Complexity? |
|------------|---------------|-------------|
| `CtIf` | if/else | +1 |
| `CtFor` | for loop | +1 |
| `CtForEach` | for-each loop | +1 |
| `CtWhile` | while loop | +1 |
| `CtDo` | do-while loop | +1 |
| `CtSwitch` | switch statement | per case |
| `CtCase` | case in switch | +1 (non-default) |
| `CtTry` | try/catch/finally | per catch |
| `CtThrow` | throw statement | 0 |
| `CtReturn` | return statement | 0 |
| `CtAssignment` | x = value | 0 |
| `CtOperatorAssignment` | x += value | 0 |

### Expressions
| Spoon Type | Java Construct | Notes |
|------------|---------------|-------|
| `CtInvocation` | method call | Track for invokes metric |
| `CtConstructorCall` | `new Foo()` | Track for invokes |
| `CtFieldRead` | `obj.field` (read) | Track for reads metric |
| `CtFieldWrite` | `obj.field = x` | Track for writes metric |
| `CtVariableRead` | local var read | |
| `CtVariableWrite` | local var write | |
| `CtTypeAccess` | `ClassName.method()` | Static access |
| `CtLiteral` | `42`, `"foo"` | |
| `CtBinaryOperator` | `a + b`, `a && b` | `&&`/`||` add complexity |
| `CtUnaryOperator` | `!x`, `++i` | |
| `CtConditional` | `a ? b : c` | Ternary |
| `CtSuperAccess` | `super.method()` | |
| `CtThisAccess` | `this.field` | |
| `CtExecutableReferenceExpression` | `Foo::bar` | Method reference |
| `CtLambda` | `x -> x+1` | |
| `CtNewClass` | anonymous class | |

## Import Handling

`CtCompilationUnit` has `getImports()` returning `List<CtImport>`.
Each `CtImport` has a kind (`CtImportKind`):
- `TYPE` — `import java.util.List;`
- `ALL_TYPES` — `import java.util.*;`
- `METHOD` — `import static Math.abs;`
- `ALL_STATIC_MEMBERS` — `import static Math.*;`
- `FIELD` — `import static System.out;`

For package-level import tracking, the abstractor should look at what types
are actually *used* from other packages rather than relying solely on import
statements (since imports can be unused or wildcard).

## Generics / Parameterized Types

`CtTypeReference` has:
- `getActualTypeArguments()` — returns `List<CtTypeReference<?>>` for the
  type arguments (e.g., for `List<String>` this returns `[String]`).
- `isParameterized()` — true if has actual type arguments.

`CtFormalTypeDeclarer` (implemented by `CtType`, `CtMethod`, `CtConstructor`):
- `getFormalCtTypeParameters()` — returns `List<CtTypeParameter>` for the
  formal type parameters (e.g., for `class Foo<T>` this returns `[T]`).

`CtTypeParameter`:
- `getSuperclass()` — the bound type (e.g., `T extends Comparable` → `Comparable`).
- `getTypeErasure()` — the erasure type (e.g., `T extends Comparable` → `Comparable`).
- `getSuperInterfaces()` — additional interface bounds.

## MavenLauncher

`MavenLauncher` is the entry point for Maven projects:
- `MavenLauncher(path, SOURCE_TYPE.APP_SOURCE)` — parses only main source.
- `MavenLauncher(path, SOURCE_TYPE.TEST_SOURCE)` — parses only test source.
- `MavenLauncher(path, SOURCE_TYPE.ALL_SOURCE)` — parses both.
- After construction, call `buildModel()` to get a `CtModel`.
- The model can be queried: `model.getAllTypes()`, `model.getAllPackages()`.
- `addInputResource(path)` adds additional source directories.

The current abstractor uses `APP_SOURCE` and falls back to adding the project
path as an input resource if the model comes back empty.

## Common Pitfalls

1. **Null from `getDeclaration()`**: Always check for null. External types
   (java.lang.String, etc.) return null. Use `getTypeDeclaration()` for shadow.

2. **`CtEnum extends CtClass`**: An enum passes `instanceof CtClass`. Check
   `instanceof CtEnum` BEFORE `instanceof CtClass` in dispatch chains.

3. **`CtAnnotationType`**: Not a `CtClass` or `CtInterface`. The abstractor's
   `addTypeDesc()` doesn't handle it — will throw "Unhandled Type".

4. **Anonymous/local classes**: `tr.isAnonymous()` or `tr.isLocalType()` may be
   true. Their simple names are numeric. Need special handling.

5. **Wildcard references**: `CtWildcardReference` is a `CtTypeReference` subtype.
   Not handled by the abstractor.

6. **Intersection types**: `CtIntersectionTypeReference` represents `T extends A & B`.
   May appear in type parameter bounds and wildcard bounds.

7. **Method references**: `CtExecutableReferenceExpression` (`Foo::bar`) needs to
   be tracked as invocations in metrics.

8. **Records**: `CtRecord extends CtClass`. Java 16+. May appear in newer projects.

9. **Sealed classes**: Java 17 feature. `getPermittedTypes()` on `CtSealable`.

10. **Switch expressions**: `CtSwitchExpression` (Java 14+) vs `CtSwitch` (statement).
    Both have cases that add complexity.
