# Rough Idea: Finish the Java Abstractor

## Goal

Complete the Java Abstractor so it can successfully process all Java projects
listed in the Technical Debt Dataset (https://github.com/clowee/The-Technical-Debt-Dataset).

A local copy of the dataset's SQLite database is at `javaAbstractor/tdd/td_V2.db`.

## Definition of "Finished"

The abstractor must be able to run on each Java project in the dataset and produce
JSON (or YAML) output that conforms to the Generalized Feature Definition
specified in `docs/genFeatureDef.md`.

## Technical Approach

- Uses Spoon to parse and analyze Java source code.
- The existing codebase already has significant infrastructure: constructs,
  JSON serialization, diff tooling, metrics, comparison, validation, etc.
- The output must match the schema: Project, Abstract, Argument, Basic, Field,
  InterfaceDecl, InterfaceDesc, InterfaceInst, Locations, Method, MethodInst,
  Metrics, Object, ObjectInst, Package, Selection, Signature, StructDesc,
  TypeParam, Value.

## Constraints

- This is a college research project (PhD research). The author MUST remain in
  full control of all code changes.
- The agent may NOT run `git add`, `git commit`, `git push`, or any destructive
  git commands. Only `git list`, `git fetch`, and `git diff` are permitted.
- Every step must be iterative and interactive:
  1. User tells the agent what to work on.
  2. Agent returns a plan for that step WITHOUT changing code.
  3. User adjusts the plan as needed.
  4. User asks the agent to write the changes.
  5. Agent stops when changes are done so user can review.
- No help with PRs is needed.
