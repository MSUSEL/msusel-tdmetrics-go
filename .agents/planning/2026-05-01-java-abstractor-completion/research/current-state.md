# Research: Current State of the Java Abstractor

## Architecture Overview

The Java Abstractor uses **Spoon** (v11.2.0) to parse Java source via Maven
(`MavenLauncher`) or inline source (`prepareClassesFromSource` for tests). Output
conforms to `docs/genFeatureDef.md`.

### Entry Point & Flow

1. `abstractor.app.App.main` → `Config.FromArgs` → `App.run`.
2. `Abstractor.prepareMavenProject(path)` builds a `CtModel` (may retry with
   `addInputResource` for small `testData/java` fixtures).
3. `Abstractor.performAbstraction()`:
   - Drains `pendingPackages` → `processPackage` (types per package).
   - After each package pass, `processPendingMetrics()` (batch copy/clear/process).
   - Then `consolidateCons()` → `crossConnectConstructs()` → `Validator.validate()`.
4. `Project.toJson(JsonHelper)` — CLI toggles `writeKinds`, `writeIndices`, `writeRefs`.

There is **no** `Resolver` class yet; post-walk work lives on `Abstractor`.

### Key classes

| Area | Class | Role |
| --- | --- | --- |
| Walk | `Abstractor` | Declarations, types, enums, objects, interfaces, metrics queue |
| Bodies | `Analyzer` | Cyclomatic complexity, getter/setter, invokes/reads/writes |
| Well-known | `Baker` | `anyDesc`, `$Array` + `InterfaceInst`, boxed → `Basic` |
| Spoon | `SpoonUtils` | `describeElem`, package name/path, `isObject` / `isNull` helpers |
| Validate | `Validator` | Resolved refs, construct graph sanity |
| Model | `Project` + `Factory<T>` / `Ref<T>` | All construct kinds |

### Constructs (population status)

| Construct | Status |
| --- | --- |
| Abstract, Argument, Basic, Field, Signature, StructDesc, Selection | Populated |
| MethodDecl | Populated (incl. constructors, `constructor` flag, `isStatic`) |
| Metrics | Populated; writes / some invoke paths incomplete |
| ObjectDecl, InterfaceDecl, InterfaceDesc | Populated; `inherits` wired |
| Value | Enum constants only; package static fields not yet |
| InterfaceInst | Baker arrays; user/external parameterized types partial |
| ObjectInst | Created in skeleton form; finisher incomplete |
| MethodInst | Class exists; not populated by walk |
| PackageCon | Populated; **`imports` empty** (TODO in `performAbstraction`) |

## Type dispatch (`addTypeDesc`)

- Primitives → `Basic`; arrays → Baker `arrayInst`.
- Wildcards → bound type or `anyDesc` when unbounded (`Object` bound treated as unbounded).
- Boxed / `String` → `basicForBoxedOrString`.
- **`tr.isShadow()`** → `addShadowTypeDesc` → **`anyDesc`** (named JDK stubs planned — see plan Step 5).
- Anonymous / local types → `null` (notice); they must not become separate declarations.
- Annotation types → notice, `null`.
- User `CtEnum` / `CtClass` / `CtInterface` / type parameters → respective adders.

`addDeclaration` skips annotation types; orders **`CtEnum` before `CtClass`**.

## Enums (`addEnum`)

- `ObjectDecl` with struct containing only **`$value`** (enum superclass type).
- Each `CtEnumValue` → package **`Value`** (`const: true`); **`Value.type` still null** in finisher.
- User enum methods and super-interface wiring **not** in finisher yet.

## Objects & interfaces

- `addObjectDecl`: constructors (non-implicit), instance methods, synthesized
  `InterfaceDesc` (non-static instance abstracts), **`inherits`** from
  `getSuperInterfaces()`, nested type references, generic **`addObjectInstances`** hook.
- `addInterfaceDecl`: abstracts, **`inherits`**, **`pin`** for nested interfaces.
- Object synthesized interface uses **`pin`** = object ref.

## Metrics (`Analyzer`)

- **Implemented:** complexity, line/code counts, indents, getter/setter detection,
  `CtInvocation` → invokes, field read/reference → reads.
- **TODO:** `addAssignmentUsage` (writes), `addExecutableReferenceUsage`;
  `logElementTree` / `logUsage` hardcoded `true`.

## Cross-connect & validation

- `crossConnectConstructs` adds method decls, object decls, interface decls, and
  values to their `PackageCon`.
- Validation throws if `log.errorCount() > 0` after `Validator.validate`.

## Tests

| Test class | Coverage |
| --- | --- |
| `AppTests` | `test0001`, `test0002` (Maven); `test1001`–`test1006` (single-file) |
| `RobustnessTests` | Wildcards, annotations, boxing, etc. (no YAML) |
| `MetricsTests` | Metric fragments via `Tester.checkConstruct` |
| `JsonTests`, `DiffTests`, `IterTests` | Infrastructure |

**Note:** `test1006` source is nested generics, not an enum — plan Step 1 expects
an enum fixture there.

## Build

- Java 17, Maven; `mvn clean compile assembly:single` → fat jar.
- `mvn test` — JUnit 5.

## Active plan

Remaining steps: `.agents/planning/2026-05-01-java-abstractor-completion/implementation/plan.md` (11 steps). **Next: Step 1 — enum completion.**
