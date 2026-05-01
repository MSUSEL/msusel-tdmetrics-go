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
| InterfaceDecl | `InterfaceDecl.java` | Partial (nested interfaces, instances) |
| InterfaceDesc | `InterfaceDesc.java` | Partial (inheritance, pinning) |
| InterfaceInst | `InterfaceInst.java` | Partial (only arrays use it currently) |
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
  `$get`, `$set` abstracts for array types, `Object` as empty interface).
- **Validator**: Checks that all refs are resolved and point to constructs in factories.
- **Diff**: Hirschberg & Wagner diff algorithms for test output comparison.
- **JSON**: Custom JSON parser, formatter, and serialization (`Jsonable` interface).
- **Comparison**: `Cmp`, `CmpOptions`, `CmpContext` for construct deduplication.
- **Ref<T>**: Lazy-resolution references supporting circular dependencies.

## Known TODOs and Gaps (from source code)

### Critical / Functional Gaps

1. **Package imports** (`Abstractor.java:104-151`): `getImports()` is stub code with
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

8. **Generic instances** (`ObjectInst`, `MethodInst`, `InterfaceInst`): Only
   `InterfaceInst` is used (for arrays via Baker). Real generic instantiations from
   source code are not tracked.

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

- **AppTests**: 5 tests total
  - `test0001`, `test0002`: Full Maven project tests in `testData/java/test000N/`.
    Runs the full abstractor pipeline and compares JSON output against
    `abstraction.yaml` expected files.
  - `test1001`, `test1002`, `test1003`: Single-class tests using `Tester.addClassFromFile()`.
    Parses a `.java` file directly (not Maven) and compares output.

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
