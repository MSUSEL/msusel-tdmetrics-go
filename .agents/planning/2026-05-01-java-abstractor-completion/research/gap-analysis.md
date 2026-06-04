# Research: Gap Analysis

Snapshot aligned with `javaAbstractor/` source and the **remaining** implementation plan.

## Schema coverage (`genFeatureDef.md`)

| Construct / field | Status |
| --- | --- |
| Core decls (objects, interfaces, methods, signatures, metrics skeleton) | ✅ Largely present |
| `InterfaceDesc.inherits` | ✅ Wired for classes and interfaces |
| `InterfaceDesc.pin` | ⚠️ Objects and nested interfaces; declared top-level interfaces partial |
| `Value` (package constants) | ⚠️ Enum constants only; static fields missing |
| `ObjectDecl.nest` | ❌ Not on declaration (only `$nest` struct field) |
| `Package.imports` | ❌ TODO in `Abstractor` |
| `objectInsts` / `methodInsts` | ⚠️ / ❌ Skeleton / empty |
| Shadow JDK types as named stubs | ❌ Currently `anyDesc` via `addShadowTypeDesc` |

## Priority gaps (for TDD-scale runs)

### P0 — Correctness on real projects

1. **Enum completion** — methods, `Value.type`, super-interfaces; fix `test1006`.
2. **Metrics writes & method references** — participation depends on `writes` and full `invokes`.
3. **External type stubs** — bounded named stubs instead of collapsing shadow types to `anyDesc` only.

### P1 — Schema completeness

4. **Package-level static values** — `public static final` / static fields as `Value`.
5. **`ObjectDecl.nest`** — named nested classes distinct in Cmp/JSON.
6. **Generic instances** — finish `ObjectInst`, then `MethodInst`.
7. **Anonymous / lambda metrics** — ensure body walk under enclosing method.

### P2 — Pipeline & polish

8. **Resolver extraction** — ordered post-walk phases (imports, interface generation timing).
9. **Package imports** — derive from type usage in Resolver step.
10. **Analyzer debug flags** — verbose-driven, not hardcoded `true`.
11. **Goldens** — `test0001`/`test0002` and single-file fixtures as features land.

### P3 — Validation

12. **TDD script** — run abstractor on cloned Apache projects from `td_V2.db`.

## Research-critical metrics

| Need | Status |
| --- | --- |
| Method ↔ class membership | ✅ `ObjectDecl.methodDecls` |
| Per-method complexity / size | ✅ Mostly |
| Invocations | ⚠️ `CtInvocation` yes; method refs / constructors incomplete |
| Reads / writes for participation | ⚠️ Reads yes; writes TODO |
| Package dependency graph | ❌ `imports` empty |

## Suggested order

Matches `implementation/plan.md` Steps 1–11: enums → values → nest → anonymous/lambda metrics → JDK stubs → metrics completion → generics → resolver → imports → cleanup → TDD script.
