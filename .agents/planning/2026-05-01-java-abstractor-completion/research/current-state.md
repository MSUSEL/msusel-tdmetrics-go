# Research: Current State of the Java Abstractor

## Architecture Overview

The Java Abstractor uses **Spoon** (v11.2.0) to parse Java source code via Maven
projects (`pom.xml`). It produces JSON/YAML output conforming to the Generalized
Feature Definition (`docs/genFeatureDef.md`).

### Entry Point & Flow

1. `abstractor.app.App` is the main class.
2. `abstractor.app.Config` parses CLI args: `-i` (input project path with pom.xml),
   `-o` (output JSON/YAML), `-v` (verbose), `-m` (minimize), `-e` (extra debug info).
3. `Abstractor` orchestrates the parsing:
   - `addMavenProject(path)` uses Spoon's `MavenLauncher` to build a `CtModel`.
   - Iterates over all packages, types, and declarations.
   - Type references: **`addTypeDesc`** resolves primitives, arrays, wildcards,
     user types, and (Steps 1–2) shadow/JDK types via **`addExternalStub`** with
     boxing for wrappers and `String`.
   - `finish()` runs: `processPendingMetrics()` → `consolidateCons()` →
     `crossConnectConstructs()` → `validate()`.
4. `Analyzer` computes metrics (complexity, line counts, indents, getter/setter
   detection, invocations, reads, writes).
5. `Project` holds all `Factory<T>` instances for each construct type.
6. `Factory<T>` manages deduplication via `CtElement`-keyed maps and `TreeSet`-based
   comparison/consolidation.

### Constructs Implemented (Data Model)

All 18 construct types from `genFeatureDef.md` have corresponding Java classes:

| Construct | Class | Status |
|-----------|-------|--------|
| Abstract | `Abstract.java` | Implemented |
| Argument | `Argument.java` | Implemented |
| Basic | `Basic.java` | Implemented |
| Field | `Field.java` | Implemented |
| InterfaceDecl | `InterfaceDecl.java` | Partial (nested interfaces; **Step 2:** JDK stubs) |
| InterfaceDesc | `InterfaceDesc.java` | Partial (inheritance, pinning) |
| InterfaceInst | `InterfaceInst.java` | Partial (**Baker** arrays; **Step 2:** external `List<>` / `Map<>` etc.) |
| Locations | `Location.java` / `Locations` | Implemented |
| MethodDecl | `MethodDecl.java` | Implemented (TODO: constructor flag) |
| MethodInst | `MethodInst.java` | Exists but not populated by Abstractor |
| Metrics | `Metrics.java` | Mostly implemented |
| ObjectDecl | `ObjectDecl.java` | Implemented |
| ObjectInst | `ObjectInst.java` | Exists but not populated by Abstractor |
| Package | `PackageCon.java` | Partial (imports unfinished) |
| Selection | `Selection.java` | Implemented |
| Signature | `Signature.java` | Implemented |
| StructDesc | `StructDesc.java` | Implemented |
| TypeParam | `TypeParam.java` | Partial |
| Value | `Value.java` | Exists but not populated by Abstractor |

### Supporting Infrastructure

- **Baker**: Pre-builds synthetic constructs (e.g., `$Array` interface with `$len`,
  `$get`, `$set` abstracts for array types); **`basicForBoxedOrString`** (Step 2)
  maps boxed JDK types and `java.lang.String` to shared **`Basic`** nodes.
- **Validator**: Checks that all refs are resolved and point to constructs in factories.
- **Diff**: Hirschberg & Wagner diff algorithms for test output comparison.
- **JSON**: Custom JSON parser, formatter, and serialization (`Jsonable` interface).
- **Comparison**: `Cmp`, `CmpOptions`, `CmpContext` for construct deduplication.
- **Ref<T>**: Lazy-resolution references supporting circular dependencies.

## Known TODOs and Gaps (from source code)

### Critical / Functional Gaps

*(Steps 1–2: robust `addTypeDesc` / `addDeclaration`, wildcard and anonymous/local
handling, **`addExternalStub`** + boxing. Remaining gaps below.)*

1. **Package imports** (`Abstractor` — package imports TODO near `addPackage`): `getImports()` is stub code with
   debug `println` statements. Returns `null`. This means package dependency tracking
   is completely missing.

2. **Interface inheritance** (`Abstractor.java:566`): `InterfaceDesc` finisher has
   `// TODO: Implement Inheritance`. Super-interfaces are not connected.

3. **Interface pinning** (`Abstractor.java:562`): `InterfaceDesc` doesn't set the `pin`
   field linking the description back to its declaration.

4. **Super-interface connections on classes** (`Abstractor.java:262`): `addObjectDecl`
   has `// TODO: Finish implementing` with commented-out code for
   `c.getSuperInterfaces()`.

5. **Nested interface handling** (`Abstractor.java:461`): Nested interfaces log an error
   but aren't properly differentiated.

6. **Enum handling** (`Abstractor.java:393, 474-496`): `addEnum()` exists but is
   incomplete. Enum values are not added as constants to packages.

7. **Values (package-level variables/constants)** are not extracted from source code at
   all.

8. **Generic instances** (`ObjectInst`, `MethodInst`, `InterfaceInst`): **`InterfaceInst`**
   is used for Baker arrays and **external** parameterized types (Step 2).
   **`ObjectInst` / `MethodInst`** and rich user-side generic tracking are still open
   (later steps).

9. **Metrics usage tracking** (`Analyzer.java:285, 308`): `addAssignmentUsage` and
   `addExecutableReferenceUsage` are TODO stubs.

10. **Cross-connection** (`Abstractor.java:671`): `crossConnectConstructs()` has
    `// TODO: Add more to packages` — interface declarations and values are not added
    to packages.

### Minor / Cleanup

- Debug `println` statements throughout `getImports()`.
- `logElementTree` and `logUsage` flags hardcoded to `true` in `Analyzer.java`.
- Test file `test0001/abstraction.yaml` has `youShallNotPass: true` intentionally
  preventing the test from passing.

## Existing Tests

### Test Structure

- **AppTests**: Maven integration (`test0001`, `test0002`) plus single-file fixtures
  **`test1001`–`test1005`** under `testData/java/test100N/` (`Foo.java` + `abstraction.yaml`).
  `test1004` / `test1005` align with Steps 1–2. **`test1002`** avoids `System.out.println`
  in the fixture (Spoon otherwise pulls a large JDK graph for `System`). Goldens may
  still need updates per environment; **`MetricsTests`** may be out of sync until refreshed.

- **RobustnessTests**: Smoke tests for Spoon edge cases (wildcards, annotations, boxing fields, etc.).
- **MetricsTests**: Tests for metrics computation (complexity, indents, getter/setter detection).
- **JsonTests**: Tests for JSON parser/formatter.
- **IterTests**: Tests for iterator utilities.
- **DiffTests**: Tests for diff algorithms.

### Test Infrastructure

- `Tester` class provides helpers: `addClassFromSource()`, `addClassFromFile()`,
  `checkProjectWithFile()`, `checkProject()`, `checkConstruct()`.
- Uses `Diff.PlusMinusByLine()` for readable failure output.
- JSON comparison uses relaxed format (comments allowed in YAML/JSON files).

## Build

- Maven project, Java 17 (OpenJDK 17.0.14).
- Spoon 11.2.0 for AST parsing.
- `mvn clean compile assembly:single` to build.
- `mvn test` to run tests.
- Produces `abstractor-0.1-jar-with-dependencies.jar`.
