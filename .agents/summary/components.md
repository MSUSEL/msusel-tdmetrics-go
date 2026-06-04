# Components

This document enumerates the major components, their entry points, and their responsibilities. For higher-level relationships, see `architecture.md`.

## Component Map

```mermaid
graph TB
    subgraph Repo[msusel-tdmetrics-go]
        GA[goAbstractor<br/>Go module]
        JA[javaAbstractor<br/>Maven project]
        TD[techDebtMetrics<br/>.NET solution]
        Docs[docs<br/>schema & notes]
        Test[testData<br/>integration fixtures]
        Plan[.agents/planning<br/>step plans]
    end

    GA -. emits .-> Schema[(Generalized Feature Definition)]
    JA -. emits .-> Schema
    Schema -. consumed by .-> TD
    Test -. exercises .-> GA
    Test -. exercises .-> JA
    Docs -. defines .-> Schema
```

## goAbstractor

- **Entry point**: `goAbstractor/main.go` (CLI flags `-i`, `-o`, `-v`, `-m`, `-h`).
- **Library entry**: `internal/abstractor.Abstract(Config) constructs.Project`.
- **Output**: JSON tree built by `internal/jsonify`, optionally minimized.

### Sub-packages (selected)

| Package | Files of interest | Responsibility |
| --- | --- | --- |
| `internal/abstractor` | `abstractor.go` | Phase-1 walk over `golang.org/x/tools/go/packages` results. |
| `internal/abstractor/analyzer` | `analyzer.go`, `complexity/`, `usages/`, `accessor/` | Per-method analysis (cyclomatic complexity, reads/writes/invokes, accessor detection). |
| `internal/abstractor/baker` | `baker.go` | Pre-bakes well-known/basic types so they're shared across constructs. |
| `internal/abstractor/converter` | `converter.go` | Maps `go/types` to construct references. |
| `internal/abstractor/instantiator` | `instantiator.go` | Generic instantiation creation. |
| `internal/abstractor/querier` | `querier.go` | Wrapper around the Go `packages` loader. |
| `internal/abstractor/resolver` | `resolver.go`, `dce/`, `genInterfaces/`, `inheritance/`, `instantiations/`, `references/` | Phase-2 resolver passes. |
| `internal/constructs` | one folder + `.go` per construct kind, plus `factory.go`, `project/`, `packageCon/` | Construct types and factories. |
| `internal/jsonify` | – | JSON tree, minimize/format. |
| `internal/logger` | – | Indented logger. |
| `internal/locs` | – | File/line location set built on `go/token.FileSet`. |
| `internal/reader` | – | High-level "read packages" helper used by `main.go`. |
| `internal/assert`, `internal/debug`, `internal/stringer` | – | Internal helpers. |

### Tests

- `goAbstractor/tests/tests_test.go`, `tool_test.go` — integration runner over `testData/go/test*` fixtures.
- Per-package `*_test.go` files (e.g. `analyzer_test.go`).

## javaAbstractor

- **Entry point**: `abstractor.app.App.main` (fat jar via `mvn clean compile assembly:single`).
- **Library entry**: `prepareMavenProject(path)` → `performAbstraction()` → `proj.toJson(JsonHelper)`.
- **Parser**: Spoon 11.2.0 via `MavenLauncher` (or `Launcher` + `VirtualFile` in tests).

### Sub-packages

| Package | Key classes | Responsibility |
| --- | --- | --- |
| `abstractor.app` | `App`, `Config` | CLI and Apache Commons CLI parsing. |
| `abstractor.core` | `Abstractor`, `Analyzer` | Spoon walk, metrics, finish pipeline. |
| `abstractor.core.spoonUtils` | `SpoonUtils` | Element descriptions, package name/path, type guards. |
| `abstractor.core.constructs` | `Project`, `PackageCon`, `Factory`, `Ref`, `Baker`, one class per construct kind | Schema mirrors `genFeatureDef.md`. |
| `abstractor.core.json` | `JsonHelper`, `JsonNode`, `JsonFormat`, `parser/` | JSON build/parse/format. |
| `abstractor.core.cmp` | `Cmp`, `CmpOptions` | Deduplication and golden comparison ordering. |
| `abstractor.core.iter` | – | Iterator helpers. |
| `abstractor.core.diff` | Hirschberg, Wagner, comparators | YAML golden diffs in tests. |
| `abstractor.core.log` | `Logger` | Indented logging (mirrors Go). |
| `abstractor.core.validator` | `Validator` | Post-walk validation; errors fail `performAbstraction`. |

### Notable state on `Abstractor`

- `pendingPackages` — packages discovered during the walk; drained in `performAbstraction`.
- `pendingMetrics` — executables awaiting `Analyzer`; batch drain avoids `ConcurrentModificationException`.
- Shadow JDK types currently map to `Baker.anyDesc()` (`addShadowTypeDesc`); named stub cache is planned (plan Step 5).

### Tests (`javaAbstractor/src/test/java/abstractor/`)

| Class | Purpose |
| --- | --- |
| `AppTests` | `test0001`–`test0002` (Maven), `test1001`–`test1006` (single-file) vs `abstraction.yaml`. |
| `core.Tester` | `prepareClassesFromSource`, golden diff helpers. |
| `core.RobustnessTests` | Smoke tests (wildcards, annotations, boxing, etc.). |
| `core.MetricsTests` | Focused metric YAML fragments. |
| `core.JsonTests`, `core.DiffTests`, `core.IterTests` | Infrastructure. |

## techDebtMetrics

.NET 8 solution at `techDebtMetrics/techDebtMetrics.sln`.

| Project | Path | Role |
| --- | --- | --- |
| `Commons` | `Commons/` | `Data/Yaml/`, `Data/Locations/`, `Extensions/GeneralExt.cs`, `Tooling/TestHook.cs`. Shared utilities. |
| `Constructs` | `Constructs/` | One `.cs` file per schema construct (`Project.cs`, `Package.cs`, `ObjectDecl.cs`, `InterfaceDecl.cs`, `Method*.cs`, …). Loads/holds the abstractor output. Includes `IConstruct`, `IDeclaration`, `IInterface`, `IMethod`, `IObject`, `ITypeDesc` interfaces. |
| `DesignRecovery` | `DesignRecovery/` | `DesignRecovery.cs`, `IAlgorithm.cs`, `Manager.cs`, `Membership.cs`, `TestHook.cs`. Currently the participation/membership computation is mostly commented out (TODO: synthesised object for basics and projects). |
| `TechDebt` | `TechDebt/` | `Class.cs`, `Method.cs`, `Math.cs`, `Participation.cs`, `Project.cs`, `Source.cs`, `Validator.cs`. Computes WMC/TCC/ATFD-style metrics. |
| `Runner` | `Runner/Program.cs` | CLI driver. **Currently stubbed** (`throw new NotImplementedException`). |
| `UnitTests` | `UnitTests/` | `CommonsTests/{ExtensionsTests,LocationsTests}.cs`, `ConstructTests/ConstructTests.cs`. |

## Test Fixtures (`testData/`)

- `testData/go/test0001` … `test0018` — Go projects, each with `main.go` and `abstraction.yaml` (and sometimes `expStub.txt`).
- `testData/java/test0001`, `test0002`, `test1001`–`test1006` — Java fixtures; `test10NN` are single-file; `test000N` are full Maven projects.
- `testData/todo.md` — running list of fixtures-to-add.

## Documentation (`docs/`)

| File | Content |
| --- | --- |
| `genFeatureDef.md` | The canonical schema for the JSON interchange format. |
| `pipelineDiagram.svg` | The pipeline figure used in the top-level README. |
| `ducktype.md` | Notes on representing Go duck typing. |
| `extendingPointers.md` | Pointer extension rationale (Go). |
| `participationMatrix.md` | Participation-matrix definition. |
| `spoonNotes.md` | Spoon API gotchas relevant to the Java abstractor. |
| `tdResults.md` | Notes on TD analysis results. |
| `CompPaper.pdf` | Reference paper. |

## Planning and Agent Context

- `.agents/planning/2026-05-01-java-abstractor-completion/` — Java completion plan (11 remaining steps in `implementation/plan.md`).
- `.cursor/rules/java-abstractor-handoff.mdc` — concise handoff rule loaded by Cursor when editing Java abstractor files.
- `.cursor/commands/*.sop.md` — workflow SOP commands available in Cursor (`code-assist`, `code-task-generator`, `codebase-summary`, `eval`, `pdd`).
- `AGENTS.md` — researcher's binding rules for AI agents (git restrictions, plan-first workflow, code quality).
