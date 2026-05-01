# Agent Guidelines

This is a PhD research project for technical debt analysis of procedural and
object-oriented languages. The researcher must remain in full control of all
code changes.

## Git Restrictions

- **NEVER** run `git add`, `git commit`, `git push`, or any destructive git command.
- Only `git list`, `git fetch`, and `git diff` are permitted.
- Do not create PRs or branches.

## Interaction Model

Follow this strict iterative workflow for every code change:

1. **User tells the agent what to work on.**
2. **Agent returns a plan** for that step WITHOUT changing any code.
3. **User adjusts the plan** as needed.
4. **User asks the agent to write the changes.**
5. **Agent writes the changes and stops** so the user can review.

Never proceed to the next step without the user's explicit direction.

## Code Changes

- Each step should be roughly one feature/fix or a related set of changes.
- Include unit tests alongside code changes (test files in `testData/java/`
  with expected YAML output).
- The user may request an integration test first (with or without expected
  YAML) to understand the "shape of constructs" before implementation.
- Clean up debug artifacts (`println` statements, hardcoded log flags) as
  you encounter them.
- Add TODO comments for features known to be needed but not yet implemented.

## Error Handling

- Log warnings and continue; never crash on unhandled constructs.
- Since technical debt analysis is estimation interpreted by humans, the
  results can tolerate some imprecision.

## Code Quality

- Follow good practices but do not over-engineer. This is a research tool
  with a finite lifetime, not a long-lived product.
- Code must be debuggable and readable by other researchers.
- Follow existing patterns: `Factory<T>` + `Ref<T>`, `Cmp`/`CmpOptions`,
  `Jsonable` with `toJson(JsonHelper)`, `Logger` with push/pop indentation.
- New patterns are welcome where they make the code cleaner.
- Patterns loosely mirror the Go abstractor (`goAbstractor/`) for
  maintainability across both codebases.

## Project Context

- **Goal**: Complete the Java Abstractor (`javaAbstractor/`) so it can
  process all 31 Apache Java projects in the Technical Debt Dataset (TDD).
- **Output format**: JSON/YAML conforming to `docs/genFeatureDef.md`.
- **Parser**: Spoon (v11.2.0) via `MavenLauncher` for Maven projects.
- **Target**: Java 17, Maven build.
- **TDD database**: `javaAbstractor/tdd/td_V2.db` (SQLite).
- **Design docs**: `.agents/planning/2026-05-01-java-abstractor-completion/`.

## Key Design Decisions

- External (JDK/library) types: named stubs, not collapsed to `Object`.
- Annotations: use to inform analysis, do not output as constructs.
- Anonymous classes and lambdas: fold into enclosing method metrics.
- Named nested classes: separate objects with `nest` field.
- Generic instantiations: track as distinct types (`ObjectInst`,
  `MethodInst`, `InterfaceInst`).
- Package imports: derive from actual type usage, not `import` statements.
