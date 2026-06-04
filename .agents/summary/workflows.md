# Workflows

This document captures the key runtime and developer workflows.

## End-to-End Pipeline

```mermaid
sequenceDiagram
    actor Researcher
    participant GoAb as goAbstractor
    participant JavaAb as javaAbstractor
    participant File as out.json
    participant TDM as techDebtMetrics

    Researcher->>GoAb: go run main.go -i <pkg> -o out.json
    GoAb->>GoAb: load packages, walk, resolve
    GoAb->>File: write generalized JSON

    Researcher->>JavaAb: java -jar abstractor.jar -i <maven proj> -o out.json
    JavaAb->>JavaAb: MavenLauncher → CtModel → walk → JSON
    JavaAb->>File: write generalized JSON

    Researcher->>TDM: dotnet run --project Runner (planned)
    TDM->>File: load JSON via Constructs
    TDM->>TDM: DesignRecovery + TechDebt analysis
    TDM-->>Researcher: participation matrix / TD metrics
```

The .NET runner is currently stubbed; the analysis pipeline is exercised via `UnitTests` in the meantime.

## goAbstractor — Internal Flow

```mermaid
sequenceDiagram
    participant Main as main.go
    participant Reader as reader
    participant Ab as abstractor.Abstract
    participant Walk as Phase 1 walk
    participant Res as resolver
    participant J as jsonify

    Main->>Reader: load packages from -i path
    Reader-->>Main: []*packages.Package
    Main->>Ab: Abstract(Config{Packages, Log})
    Ab->>Walk: querier + baker + converter + analyzer
    Walk-->>Ab: project populated with constructs
    Ab->>Res: Resolve(proj, querier, log, skipDead)
    Res->>Res: dce, genInterfaces, inheritance, instantiations, references
    Res-->>Ab: resolved project
    Ab-->>Main: constructs.Project
    Main->>J: build JSON tree
    J-->>Main: write to -o file (or stdout)
```

## javaAbstractor — Internal Flow

```mermaid
sequenceDiagram
    participant App as App.main
    participant Cfg as Config
    participant Ab as Abstractor
    participant Spoon as MavenLauncher
    participant Proj as Project
    participant Json as JsonHelper

    App->>Cfg: FromArgs(args, null)
    App->>Ab: new Abstractor(log, proj)
    App->>Ab: prepareMavenProject(cfg.input)
    Ab->>Spoon: MavenLauncher.buildModel()
    Spoon-->>Ab: CtModel
    App->>Ab: performAbstraction()
    Ab->>Ab: pendingPackages → processPackage
    Ab->>Ab: processPendingMetrics() (batch drain)
    Ab->>Ab: consolidateCons → crossConnectConstructs
    Ab->>Ab: Validator.validate
    App->>Proj: toJson(JsonHelper)
    Proj-->>App: JsonNode
    App->>Json: format(stream, node, "")
```

Notable details:

- `addTypeDesc` uses `getTypeDeclaration()`; anonymous/local types return `null` (notice); `<nulltype>` → `anyDesc`.
- Wildcards map to bounds or `anyDesc` when unbounded (`Object` bound treated as unbounded).
- Boxed / `String` → `Baker.basicForBoxedOrString`; **`tr.isShadow()`** → `addShadowTypeDesc` → **`anyDesc`** (stubs planned).
- `InterfaceDesc.inherits` filled from `getSuperInterfaces()` on classes and interfaces.
- `processPendingMetrics` batch-drains to avoid `ConcurrentModificationException`.
- Enum constants become package `Value`s; full enum modeling is plan Step 1.

## Developer Workflow (per `AGENTS.md`)

The researcher controls all changes; agents follow this strict iterative loop:

```mermaid
flowchart LR
    A[1. User: 'work on X'] --> B[2. Agent: produce plan, no code]
    B --> C[3. User: adjust plan]
    C --> D[4. User: 'write the changes']
    D --> E[5. Agent: write changes, stop]
    E --> F[User reviews & commits]
    F --> A
```

- Each step ≈ one feature or coherent set of changes.
- Tests accompany code: `testData/java/test*` with expected `abstraction.yaml`.
- The user may request an integration test first — even with no expected YAML — to see the "shape of constructs" before implementation lands.
- Agents may modify files in `.agents/`, `.cursor/`, and `AGENTS.md` without asking; for any other path, ask first.
- Agents must **never** run `git add`, `git commit`, `git push`, or create PRs/branches. Only `git status`, `git fetch`, and `git diff` are permitted.

## Running Tests

| Component | Command |
| --- | --- |
| All | `make test` |
| Go | `cd goAbstractor && go test -count=1 ./...` |
| Java | `cd javaAbstractor && mvn test` (filter: `-Dtest="abstractor.AppTests#test0001"`) |
| .NET | `cd techDebtMetrics && dotnet test` (filter: `--filter StubTest0007`) |

For the Java side, `mvn clean compile assembly:single` builds the runnable jar. The fixtures under `testData/java/test10NN` are single-file Tester fixtures; lower-numbered fixtures are full Maven projects exercised by `AppTests`.

## Adding a New Construct or Schema Change

1. Update `docs/genFeatureDef.md` to define/describe the construct.
2. Add the construct in `goAbstractor/internal/constructs/<kind>` and wire it through factories, JSON output, and the resolver as needed.
3. Add the mirror in `javaAbstractor/src/main/java/abstractor/core/constructs/<Kind>.java`.
4. Add the consumer in `techDebtMetrics/Constructs/<Kind>.cs`.
5. Add/extend a fixture in `testData/<lang>/test*` and update the `abstraction.yaml` golden.
6. Run `make test`.

## Java Abstractor Plan Iteration

Active plan: `.agents/planning/2026-05-01-java-abstractor-completion/implementation/plan.md` — **11 steps**: enum completion, package values, `nest`, anonymous/lambda metrics, JDK stubs, metrics writes/refs, generics, resolver, imports, cleanup, TDD script.
