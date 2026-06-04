# Implementation Plan: Java Abstractor Completion

Remaining work only. Earlier milestones (robust type dispatch, boxing, interface
`inherits`, enum constants as `Value`, constructors as `MethodDecl`, batched
`processPendingMetrics`, `crossConnectConstructs` for objects/interfaces/methods/values)
are already in `javaAbstractor/` and are not listed as steps below.

## Checklist

- [ ] Step 1: Enum completion (methods, super-interfaces, fixture)
- [ ] Step 2: Package-level values (static fields and constants)
- [ ] Step 3: Named nested classes (`nest` on `ObjectDecl`)
- [ ] Step 4: Anonymous class and lambda metrics folding
- [ ] Step 5: External JDK/library type stubs (shadow types)
- [ ] Step 6: Metrics — writes, method references, constructor calls
- [ ] Step 7: Generic instantiation (`ObjectInst`, then `MethodInst`)
- [ ] Step 8: Resolver pipeline extraction
- [ ] Step 9: Package imports from type usage
- [ ] Step 10: Analyzer/debug cleanup and golden refresh
- [ ] Step 11: TDD project validation script

---

## Step 1: Enum completion

**Objective:** Enums match the schema and Go abstractor: full object shape, package
constants, user methods only.

**Already in code:** `addEnum` emits `ObjectDecl` with struct field `$value`;
enum constants become `Value` entries (`const = true`, type ref to enum object).
`addDeclaration` checks `CtEnum` before `CtClass`.

**Still needed:**
- Add user-defined enum methods via `addMethod` (skip compiler-generated
  `values()` / `valueOf()` and similar).
- Connect `e.getSuperInterfaces()` on the enum’s synthesized interface (mirror
  `addObjectDecl` finisher).
- Set `Value.type` to the enum `ObjectDecl` (currently `null` in the finisher).
- Replace or add fixture: `testData/java/test1006` should exercise enums (today
  it holds nested generics — move that scenario to a new `test1007+` if needed).

**Tests:** `AppTests#test1006` + `abstraction.yaml` golden for a small enum
(constants, field, constructor, one user method).

---

## Step 2: Package-level values (static fields and constants)

**Objective:** `public` / package-visible `static` fields become `Value` constructs
in the enclosing package (not only enum constants).

**Guidance:**
- In `addObjectDecl` finisher (or a dedicated helper), scan static fields on
  the top-level type; emit `Value` with `const` from `field.isFinal()`, type from
  `addTypeDesc`, optional initializer metrics (match Go abstractor).
- Respect visibility; skip private statics unless the schema requires them.

**Tests:** New fixture (e.g. `test1007`) with `public static final int MAX` and
`public static String name`.

---

## Step 3: Named nested class support

**Objective:** Named inner/nested classes are separate `ObjectDecl`s with schema
`nest` pointing at the enclosing declaration.

**Already in code:** `$nest` field on `StructDesc` for nested types; nested types
are discovered via `addTypeDesc` on `getNestedTypes()`.

**Still needed:**
- `ObjectDecl.nest` field + JSON/Cmp output.
- Set `nest` when `CtRole.NESTED_TYPE` and not anonymous/local.
- Same for nested `InterfaceDecl` if required by schema.

**Tests:** Fixture with outer + named inner class; assert distinct objects and
`nest` link.

---

## Step 4: Anonymous class and lambda folding

**Objective:** No separate objects for anonymous/local types; their bodies
contribute to the enclosing method’s `Metrics`.

**Already in code:** `addTypeDesc` returns `null` for anonymous/local types
(logged at notice level).

**Still needed:**
- Ensure `Analyzer` walks `CtNewClass` bodies and `CtLambda` bodies under the
  enclosing executable.
- Confirm complexity / invokes / reads include that code.
- Verify named nested classes from Step 3 are not folded.

**Tests:** Method with anonymous `Runnable` and a lambda; no extra `objects`,
metrics on enclosing method reflect both.

---

## Step 5: External JDK/library type stubs (shadow types)

**Objective:** Shadow/external references use bounded named stubs, not only
`Baker.anyDesc()`.

**Already in code:** `Baker.basicForBoxedOrString` maps boxed primitives and
`String` to shared `Basic`s. `addShadowTypeDesc` currently returns `anyDesc()`.

**Still needed:**
- Restore or implement stub `InterfaceDecl` cache (erasure-qualified name);
  parameterized shadow refs → `InterfaceInst` with resolved desc.
- Route `tr.isShadow()` and unresolved declaration paths through stubs (keep
  output bounded — log-and-continue on failure).

**Tests:** Fixture using `List<String>`, `Map<…>`, etc.; goldens show stub
declarations / insts, not a single generic `any`.

---

## Step 6: Metrics — writes, references, constructor calls

**Objective:** Participation-oriented metrics: field writes, method-reference
invokes, constructor calls.

**Already in code:** `CtInvocation` → `invokes`; `CtFieldRead` / `CtFieldReference`
→ `reads` (with null-safe field resolution).

**Still needed:**
- Implement `addAssignmentUsage` (LHS `CtFieldWrite` → `writes`).
- Implement `addExecutableReferenceUsage` and constructor-call paths.
- Tie `logElementTree` / `logUsage` to `Config.verbose` (currently hardcoded
  `true` with TODOs).

**Tests:** Extend `MetricsTests` or add `test100N` fixture for reads/writes/invokes
lists in YAML.

---

## Step 7: Generic instantiation (`ObjectInst`, then `MethodInst`)

**Objective:** Concrete generic uses become distinct instance constructs.

**Already in code:** `addObjectInstances` scans the model for `CtTypeReference` /
`CtConstructorCall` matching a generic class; `addObjectInst` creates skeleton
`ObjectInst` (struct partially resolved; interface desc TODO; finisher commented).

**Still needed:**
- Finish `ObjectInst` resolved struct/interface and register on `ObjectDecl.instances`.
- Deduplicate by (generic + type args).
- `MethodInst` for generic method calls and per-`ObjectInst` method signatures.

**Tests:** `Box<T>` with `Box<String>` and `Box<Integer>`; later generic method
calls (`test1014`-style scenario).

---

## Step 8: Resolver pipeline extraction

**Objective:** Post-walk phases live in `abstractor.core.Resolver`, not only
inline in `Abstractor.performAbstraction`.

**Guidance:**
- Move `consolidateCons`, `crossConnectConstructs`, `validate`, and future
  `expandGenericInstances` / `generateInterfaces` / import derivation into ordered
  resolver steps.
- `performAbstraction()` ends with `new Resolver(log, proj).resolve()`.

**Tests:** Full `mvn test` after refactor — no golden changes expected if behavior
is unchanged.

---

## Step 9: Package imports from type usage

**Objective:** `PackageCon.imports` lists packages actually referenced by
constructs in that package (not Java `import` statements).

**Guidance:**
- Run after types and cross-connect are stable (Resolver step).
- Walk signatures, fields, metrics type refs; add foreign package refs.
- Remove any leftover import-stub / debug code.

**Tests:** Update `test0002` golden or add two-package fixture with cross-package
type use.

---

## Step 10: Analyzer/debug cleanup and golden refresh

**Objective:** Tests and output are stable for CI and researchers.

**Guidance:**
- Set Analyzer debug flags from verbose config.
- Refresh `AppTests` Maven goldens (`test0001`, `test0002`) and single-file
  `test1001`–`test1006` as needed.
- Align misplaced fixtures (e.g. enum vs nested-generics for `test1006`).

**Tests:** `mvn test` clean run.

---

## Step 11: TDD project validation script

**Objective:** Script to run the abstractor on cloned TDD projects without crashing.

**Guidance:**
- `javaAbstractor/scripts/run-tdd.sh` (or `.ps1` on Windows): build jar, run on
  project path, sanity-check JSON size / object+method counts.
- Document starter projects (`commons-cli`, `commons-exec`, `commons-dbutils`).
- Query `td_V2.db` for recommended commits; record remaining issues as TODOs.

**Tests:** Manual on 2–3 small Apache projects.
