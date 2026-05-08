# Knowledge Base Index — msusel-tdmetrics-go

> **For AI assistants**: this file is the primary entry point into the generated documentation under `.agents/summary/`. Add **only this file** to context for routine questions; pull additional files (listed below) into context when the question targets their specific domain. Each entry includes a one-paragraph summary so you can decide whether to fetch it.

## How to Use This Knowledge Base

1. Start here. Skim the table of contents to identify likely-relevant files.
2. Pull `codebase_info.md` for a fast orientation if you don't know the repo at all.
3. Pull `architecture.md` for "how does it fit together" / "what runs in what order" questions.
4. Pull `components.md` for "where is X" / "which file/package owns Y" questions.
5. Pull `interfaces.md` for CLI flags, file-format contract, or in-process API questions.
6. Pull `data_models.md` for construct/schema/type-system questions.
7. Pull `workflows.md` for end-to-end pipeline, developer loop, or Java-abstractor flow questions.
8. Pull `dependencies.md` for external library / tooling / version questions.
9. Pull `review_notes.md` for known gaps, inconsistencies, and recommended improvements to this documentation set.
10. For active research direction, also consult `.agents/planning/2026-05-01-java-abstractor-completion/` (especially `implementation/plan.md` and `summary.md`).

`AGENTS.md` (repo root) holds **binding** rules from the researcher (git restrictions, plan-first interaction model, file-modification permissions). Always honor it.

## Quick Repo Snapshot

- **What it is**: PhD research pipeline for technical-debt analysis of procedural and OO languages (Go + Java today).
- **Three components**: `goAbstractor/` (Go 1.25), `javaAbstractor/` (Java 17 + Spoon), `techDebtMetrics/` (.NET 8).
- **Shared contract**: `docs/genFeatureDef.md` — the JSON schema both abstractors emit and the .NET side consumes.
- **Active workstream**: completing the Java abstractor against the 31-project TDD; Steps 1–2 of a 15-step plan are done; Step 3 (enums) is next.
- **Researcher controls all changes**: agents plan first, write only when asked, and never run destructive git commands. See `AGENTS.md`.

## Documentation Files

### `codebase_info.md`
**Use for**: first-look orientation, repo layout, languages/tooling table, current research status.
**Summary**: Top-level identity of the project; high-level pipeline diagram; per-directory role; per-component language/build/test/dependency table; cross-component conventions (Factory/Ref, Cmp, Jsonable, Logger); pointer to the active 15-step plan.

### `architecture.md`
**Use for**: how the components fit together, internal package layouts, two-phase abstraction story, cross-cutting concerns.
**Summary**: System architecture (multi-component pipeline communicating via JSON). Architectural style notes (schema-first, mirror-but-don't-share). Per-component flow diagrams: goAbstractor (querier → baker → converter → analyzer → resolver), javaAbstractor (Spoon `MavenLauncher` → walk → factories/baker/analyzer → JSON), techDebtMetrics (Constructs → DesignRecovery + TechDebt → Runner). Cross-cutting concerns: error handling philosophy, logging, locations, Cmp, external-type stubbing.

### `components.md`
**Use for**: "where is X implemented?" — finds packages, classes, files for a given concern.
**Summary**: Component map plus per-component sub-package tables. For goAbstractor: every `internal/*` package and what it does. For javaAbstractor: every `abstractor.*` package, plus notable state (`pendingMetrics`, `externalInterfaceStubByErasure`) and the test class catalog. For techDebtMetrics: each .NET project (Commons, Constructs, DesignRecovery, TechDebt, Runner, UnitTests) and its responsibilities, including the fact that `Runner/Program.cs` is currently a `NotImplementedException` stub. Also covers `testData/`, `docs/`, planning, and Cursor agent rules.

### `interfaces.md`
**Use for**: CLI flag references, the file-format contract, in-process API entry points, CI integration.
**Summary**: Diagram of integration surfaces. CLI flag tables for `goAbstractor` (`-i`, `-o`, `-v`, `-m`, `-h`) and `javaAbstractor` (Apache Commons CLI via `Config`). Note that the .NET runner is a stub. Recap of the schema's construct/typedesc/declaration vocabulary. In-process API entry points per component (`abstractor.Abstract`, `Abstractor.addMavenProject`/`finish`/`Project.toJson`, `Constructs.Project`/`IConstruct`/`IDeclaration`/…). External integration points (Spoon, `golang.org/x/tools`, etc.). CI workflow summary and `Makefile` mapping.

### `data_models.md`
**Use for**: schema/construct/type-system questions, generic instantiation model, external-type stubbing model.
**Summary**: Class diagram of the construct hierarchy. Construct categories (declarations, type descriptions, members, metrics). Per-language materialization table mapping each construct to its Go/Java/.NET counterpart. Notes on the Go-only temp/reference constructs used during resolver passes. Generic instantiation model (`InterfaceInst`, `ObjectInst`, `MethodInst`). External/stub type model in Java (boxed/String → `Basic`, others → cached stub `InterfaceDecl` keyed by erasure-qualified name). TDD database role. Output formats (JSON canonical, YAML for goldens, minimization flag).

### `workflows.md`
**Use for**: end-to-end runtime sequences, the researcher-agent interaction loop, test commands, schema-change procedure.
**Summary**: Sequence diagrams for end-to-end pipeline, goAbstractor internal flow, and javaAbstractor internal flow (with notes on `getTypeDeclaration` vs `getDeclaration`, try/catch dispatch, `<nulltype>` handling, batched `processPendingMetrics`, and `addExternalStub` routing). Developer workflow diagram of the plan-first loop from `AGENTS.md`. Test command table. Step-by-step procedure for adding a new construct or schema change. Pointer to the 15-step Java-abstractor completion plan.

### `dependencies.md`
**Use for**: dependency versions, build tooling, external-library usage notes.
**Summary**: Per-component dependency tables. For Go: `Snow-Gremlin/goToolbox`, `golang.org/x/tools`, `yaml.v3`. For Java: Spoon 11.2.0, Commons CLI 1.9.0, JUnit 5 (junit-bom 5.11.0), plus Maven plugin pinning. For .NET: per-project structure (BCL only). Developer tooling versions (JDK 17.0.14 OpenLogic, Maven 3.9.9, Go 1.25, .NET 8). Cursor/agent tooling.

### `review_notes.md`
**Use for**: identifying documentation gaps, known inconsistencies, suggested improvements.
**Summary**: Findings from the consistency/completeness review of this documentation set. Highlights areas where source-of-truth docs (e.g. `docs/genFeatureDef.md`, `.agents/planning/.../implementation/plan.md`) should be consulted, and where this doc set is intentionally lighter than the source.

## Source-of-Truth Files Outside `.agents/summary/`

When this knowledge base is insufficient, the canonical sources are:

| Topic | File |
| --- | --- |
| Schema (the contract) | `docs/genFeatureDef.md` |
| Pipeline overview | top-level `README.md` and `docs/pipelineDiagram.svg` |
| Researcher rules for agents | `AGENTS.md` |
| Active plan (Java abstractor) | `.agents/planning/2026-05-01-java-abstractor-completion/implementation/plan.md` |
| Step handoff context | `.cursor/rules/java-abstractor-handoff.mdc` |
| goAbstractor entry | `goAbstractor/main.go`, `goAbstractor/internal/abstractor/abstractor.go` |
| javaAbstractor entry | `javaAbstractor/src/main/java/abstractor/app/App.java`, `…/core/Abstractor.java` |
| techDebtMetrics entry | `techDebtMetrics/Runner/Program.cs` (stub) |

## Example Queries

- *"Where is cyclomatic complexity computed for Go methods?"* → `components.md` (analyzer table) → `goAbstractor/internal/abstractor/analyzer/complexity/`.
- *"How are JDK types like `java.lang.String` represented in the JSON?"* → `data_models.md` (External/Stub Types) → `Baker.basicForBoxedOrString` and `Abstractor.addExternalStub`.
- *"What CLI flags does the Go abstractor accept?"* → `interfaces.md` (goAbstractor table).
- *"Which component owns the runner today?"* → `components.md` (techDebtMetrics row showing `Runner/Program.cs` is stubbed).
- *"What's the next step in the Java abstractor plan?"* → `codebase_info.md` (Current Research Status) and `.agents/planning/.../implementation/plan.md` Step 3.
- *"Can the agent commit changes for me?"* → `AGENTS.md` (Git Restrictions: never).
