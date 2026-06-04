# Dependencies

This document lists external dependencies for each component, with notes on how they're used.

## goAbstractor (`goAbstractor/go.mod`)

Module: `github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor`. Go 1.25.

| Dependency | Usage |
| --- | --- |
| `github.com/Snow-Gremlin/goToolbox` | CLI argument parsing (`argers/args`), assorted utilities (`utils`, `terrors/terror`). Maintained by the same author as this project. |
| `golang.org/x/tools` | `go/packages` loader and Go type analysis. Core to the abstractor's input. |
| `gopkg.in/yaml.v3` | YAML serialization, used for test goldens and YAML output paths. |
| `golang.org/x/mod`, `golang.org/x/sync` | Indirect (pulled by `x/tools`). |

Linting: `golangci-lint` (latest, via `golangci/golangci-lint-action@v8` in CI).

## javaAbstractor (`javaAbstractor/pom.xml`)

Java 17 (`maven.compiler.release=17`). Maven build. Project artifact: `abstractor.app:abstractor:0.1`.

| Dependency | Version | Usage |
| --- | --- | --- |
| `fr.inria.gforge.spoon:spoon-core` | 11.2.0 | AST parsing via `MavenLauncher`. Provides `CtModel`, `CtType`, `CtTypeReference`, `CtMethod`, etc. The entire walk hangs off Spoon's `Ct*` types. |
| `commons-cli` (`commons-cli`) | 1.9.0 | CLI argument parsing in `abstractor.app.Config`. |
| `org.junit.jupiter:junit-jupiter-api` | from `junit-bom` 5.11.0 | Test framework. |
| `org.junit.jupiter:junit-jupiter-params` | from `junit-bom` 5.11.0 | Parameterized tests. |

Build plugins (under `pluginManagement`):

- `maven-jar-plugin` — sets `Main-Class: abstractor.app.App` in the manifest.
- `maven-assembly-plugin` — `jar-with-dependencies` descriptor produces the runnable fat jar `target/abstractor-0.1-jar-with-dependencies.jar`.

Spoon notes worth knowing (from `docs/spoonNotes.md` and `.agents/planning/.../research/spoon-api.md`):

- Use `tr.getTypeDeclaration()` rather than `getDeclaration()` so shadow types are returned.
- Spoon may pull large JDK graphs into the model when fixture code touches `System.out`, etc. Single-file Tester fixtures should avoid `System.out.println` where possible.

## techDebtMetrics (.NET 8)

Solution: `techDebtMetrics/techDebtMetrics.sln`. Each project has its own `.csproj`. The repository pins .NET 8 in CI.

| Project | External deps |
| --- | --- |
| `Commons` | .NET 8 BCL only (custom YAML in `Data/Yaml/`). |
| `Constructs` | References `Commons`. |
| `DesignRecovery` | References `Constructs`. |
| `TechDebt` | References `Constructs`. |
| `Runner` | References analysis projects (currently a stub). |
| `UnitTests` | xUnit-style; references `Commons` and `Constructs`. |

(Exact NuGet package versions, if any, are in the individual `.csproj` files.)

## Test Data

`testData/go/test0001`–`test0018` and `testData/java/test0001`–`test1005` carry their own self-contained source plus `abstraction.yaml` goldens. They have no external dependencies beyond what each fixture imports from the standard library / JDK.

## TDD Database

`javaAbstractor/tdd/td_V2.db` (SQLite, gitignored) — 31 Apache Java projects for plan Step 11 validation. Schema: `research/td-dataset.md`.

## Developer Tooling

- **Make** — `Makefile` orchestrates per-component test/clean targets.
- **Maven** — Apache Maven 3.9.9 (per `javaAbstractor/README.md`).
- **JDK** — OpenLogic OpenJDK 17.0.14 (the version the project was built against).
- **Go** — Go 1.25 (per `go.mod` and CI's `GO_VERSION: '^1.25'`).
- **dotnet** — .NET 8 (per CI `DOTNET_VERSION: '8'`).
- **golangci-lint** — for `goAbstractor` lint (CI; install locally if running lint).
- **SpotBugs** — referenced as a possible Maven lint in `javaAbstractor/README.md`; not currently wired into CI.

## Cursor / Agent Tooling

- `.cursor/rules/java-abstractor-handoff.mdc` — automatically loaded when relevant Java abstractor files are touched.
- `.cursor/commands/*.sop.md` — SOP commands invokable via `/<name>` in Cursor chat.
- `AGENTS.md` — researcher-binding rules for AI agents.
