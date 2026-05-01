# Idea Honing: Java Abstractor Completion

Requirements clarification questions and answers will be recorded below.

---

## Q1: Definition of "works" for each TD Dataset project

When you say the abstractor must "work" for each of the 31 projects in the TD Dataset, what does success look like concretely? For example:

- **(a)** The abstractor runs without crashing and produces valid JSON/YAML (even if some constructs are incomplete or skipped).
- **(b)** The output must be complete — every class, interface, method, field, etc. in the project is represented in the output.
- **(c)** The output must be accurate enough to feed into the `techDebtMetrics` pipeline and produce meaningful participation/membership matrices.
- **(d)** Some other definition?

**A1:** The abstractor must run without crashing and produce valid JSON/YAML.
The output must be accurate enough to feed into the `techDebtMetrics` pipeline
and produce meaningful participation/membership matrices (option c). This means
the key constructs for membership (objects, methods, method-to-object receivers)
and participation (metrics: reads, writes, invocations) must be correct.

---

## Q2: Handling of external/library types

Real-world Java projects reference many types from the JDK standard library
(e.g., `java.util.List`, `java.lang.String`) and third-party dependencies.
Currently the abstractor maps unresolvable types to an empty interface
(`baker.objectDesc()`), effectively treating them as `Object`.

For your research, how should external types be handled?

- **(a)** Treat all external types as opaque/`Object` (current behavior) — the
  abstractor only models types defined in the project's own source code.
- **(b)** Model external types as named basics or named interface stubs (e.g.,
  create a basic `String`, an interface `List` with its abstracts) so they appear
  as distinct types in the output, but don't deeply analyze them.
- **(c)** Something else?

**A2:** Option (b) — model external types as named stubs. External types should
appear as distinct named types (not collapsed to `Object`) so that reads/writes/
invocations referencing them are distinguishable. However, the initial
implementation should NOT deeply analyze JDK or third-party code — just create
stub interface declarations or basics with their names.

The design should be flexible enough that if later research shows we need to
trace data/call paths *through* the JDK (e.g., value A wrapped in a JDK
container, then accessed by B), the code can be extended to include full JDK
models without a major rewrite. Keep the boundary between "project source" and
"external" clean and configurable.

---

## Q3: Which specific commit/version of each TD Dataset project?

The TD Dataset database has commit history in `GIT_COMMITS`. When running the
abstractor against each of the 31 projects, do you need to run it against:

- **(a)** Just the latest/HEAD version of each project (simplest to validate).
- **(b)** Specific commits referenced in the dataset (e.g., commits that have
  corresponding SonarQube analysis snapshots in `SONAR_ANALYSIS`).
- **(c)** Multiple commits per project to analyze TD evolution over time.
- **(d)** You haven't decided yet — for now just get it working on one version
  per project and figure out commit selection later.

**A3:** Option (b) — run against the latest commit per project that has the most
complete technical debt analysis in the TDD. This means we'll need to query the
TDD database to find commits that have corresponding SonarQube analysis data
(likely joining `GIT_COMMITS` with `SONAR_ANALYSIS`) and pick the best one per
project. The abbreviation "TDD" will be used for "Technical Debt Dataset" going
forward.

---

## Q4: Error tolerance vs. strictness

When the abstractor encounters a Java construct it can't fully model (e.g., a
complex annotation processor, a deeply nested anonymous class inside a lambda,
or an edge case in generics), what should happen?

- **(a)** Log a warning and skip/stub the construct — the abstractor should be
  resilient and produce output for everything it *can* handle, even if some
  constructs are simplified or omitted.
- **(b)** Fail with a clear error — correctness is more important than
  completeness, and you'd rather know about gaps than have silent data loss.
- **(c)** Configurable — a flag to choose between strict and lenient modes.

**A4:** Option (a) — log warnings and continue. The abstractor should try to be
as strict/correct as possible, but since TD analysis is already estimation
interpreted by humans, the results can tolerate some imprecision. When an
unhandled construct is encountered, log a warning and skip/stub it rather than
crashing. This ensures the abstractor produces output for all 31 TDD projects
even if edge cases exist.

---

## Q5: Scope of metrics needed for participation

Your dissertation defines participation based on how objects participate in
functions — via reads, writes, and invocations. Looking at the `Metrics`
construct in `genFeatureDef.md`, it tracks:

- `codeCount`, `complexity`, `indents`, `lineCount` (line-level metrics)
- `getter`, `setter` (pattern detection)
- `invokes` (methods called)
- `reads` (types read from)
- `writes` (types written to)

For producing meaningful participation matrices, which of these are essential
vs. nice-to-have?

- **(a)** All of them are needed — the full metrics set drives the analysis.
- **(b)** `invokes`, `reads`, and `writes` are the critical ones for
  participation; the line-level metrics (`codeCount`, `complexity`, etc.) are
  secondary but still useful for other TD metrics like God Class detection.
- **(c)** Something else?

**A5:** Option (a) — all metrics are needed. The full set (`codeCount`,
`complexity`, `indents`, `lineCount`, `getter`, `setter`, `invokes`, `reads`,
`writes`) should be implemented for the best possible analysis. The line-level
metrics feed into TD metrics like God Class detection (WMC), while
invokes/reads/writes drive participation matrices.

---

## Q6: Handling of annotations

Java annotations (`@Override`, `@Deprecated`, custom annotations like
`@Autowired`, `@Entity`, etc.) are pervasive in real projects. The
`genFeatureDef.md` schema doesn't have a dedicated construct for annotations.
Spoon represents annotation types as `CtAnnotationType` (a subtype of `CtType`
but not `CtClass` or `CtInterface`), which currently causes the abstractor to
crash.

How should annotations be handled?

- **(a)** Ignore annotation types entirely — skip `CtAnnotationType`
  declarations and annotation usages. They don't define methods with bodies
  or fields that participate in the program's data flow.
- **(b)** Model annotation types as interface declarations (since Java
  annotations are conceptually interfaces), even though their "methods" are
  really attribute declarations.
- **(c)** Something else?

**A6:** Option (c) — use annotations to inform the model, but don't output them
as standalone constructs. Annotations should help determine how a class is "put
together" (e.g., `@Override` confirms method inheritance, `@FunctionalInterface`
marks single-abstract-method interfaces). However, since the feature definition
must be language-agnostic (Go doesn't have the same annotation system),
annotation types themselves should NOT be modeled as declarations in the output.

When the abstractor encounters a `CtAnnotationType`, it should skip it as a type
declaration. When encountering annotation *usages* on classes/methods/fields,
it may use them to inform analysis but should not include annotation-specific
output.

This may be revisited if the TDD projects use many custom annotations that
significantly alter class behavior, but that's unlikely for these Apache
projects.

---

## Q7: Anonymous classes and lambdas

Real Java projects heavily use anonymous classes and lambdas, especially in
the Apache commons libraries (callbacks, comparators, event handlers, etc.).
These are code blocks with method bodies that contain reads, writes, and
invocations — potentially important for metrics.

How should they be handled?

- **(a)** Attribute the anonymous class/lambda body metrics to the enclosing
  method — treat them as part of the method they appear in, not as separate
  entities. This is simpler and keeps the method-to-class membership clean.
- **(b)** Model anonymous classes as nested objects and lambdas as nested
  methods with their own metrics — more accurate but more complex, and
  creates "objects" that have no meaningful name or reuse.
- **(c)** Something else?

**A7:** Between (a) and (b). The rule is:

- **Anonymous classes and lambdas**: Treat as part of the containing object/method.
  Their metrics (reads, writes, invocations, complexity) are attributed to the
  enclosing method. They are NOT modeled as separate objects.
- **Named nested classes**: These ARE modeled as separate nested objects (the
  `nest` field on `Object` in `genFeatureDef.md` supports this).

This should be consistent with whatever the Go abstractor does for nested types
and functions-within-functions. The Go abstractor should be consulted later to
verify alignment, but for now this is the working rule.

---

## Q8: Package imports — what does "imports" mean for your research?

The `Package` construct in `genFeatureDef.md` has an `imports` field listing
other packages that a package depends on. The current implementation is a stub.

For your research, is the package import list:

- **(a)** Critical — the participation/membership analysis needs to know
  inter-package dependencies to understand the project's structure.
- **(b)** Nice-to-have — the main analysis works at the method/class level and
  package imports are secondary information.
- **(c)** Not needed at all for the initial research — can be added later.

**A8:** Option (b) — nice-to-have. The main analysis works at the method/class
level. Package imports are secondary and can be implemented after the core
constructs are working correctly. Should still be done but is lower priority
than getting objects, methods, metrics, and inheritance right.

---

## Q9: Generics and instantiations

The `genFeatureDef.md` schema has `ObjectInst`, `MethodInst`, and
`InterfaceInst` for tracking generic instantiations (e.g., `List<String>` is
an instance of the generic `List<T>`). Currently only array types use
`InterfaceInst` via the Baker.

For your research, how important is tracking generic instantiations?

- **(a)** Important — different instantiations of the same generic type should
  be tracked as distinct types because they affect participation (e.g., a
  method using `Map<String, User>` participates with `User` differently than
  one using `Map<String, Order>`).
- **(b)** Low priority — for now, generic types can be modeled as their raw/
  erased form (e.g., just `Map` without type arguments). The type parameters
  on declarations should still be captured, but we don't need to track every
  concrete instantiation.
- **(c)** Something else?

**A9:** Option (a) — important. Different instantiations of the same generic
type should be tracked as distinct types. For example, a method using
`Map<String, User>` participates with `User` differently than one using
`Map<String, Order>`. This means `ObjectInst`, `MethodInst`, and
`InterfaceInst` all need to be properly populated by the abstractor when
concrete type arguments are encountered in the source code.

---

## Q10: Testing strategy for the 31 TDD projects

Running the abstractor against all 31 Apache projects requires cloning them
(potentially large repos). For validating correctness during development:

- **(a)** Start with a few small TDD projects (e.g., `commons-cli`,
  `commons-exec`, `commons-dbutils`) as integration tests, expand to all 31
  once core functionality is solid.
- **(b)** Continue building up the existing unit test suite (`testData/java/`)
  with targeted small test cases for each feature (enums, generics, inheritance,
  etc.), and only test against TDD projects manually at milestones.
- **(c)** Both — unit tests for individual features AND a small set of TDD
  projects for integration testing.
- **(d)** Something else?

**A10:** Option (b) with manual validation. Build up the unit test suite in
`testData/java/` with targeted small test cases for each feature (enums,
generics, inheritance, nested classes, etc.) with expected output YAML files.

For TDD projects, use a script (not integration tests) to run the abstractor
against a few small projects (e.g., `commons-cli`, `commons-exec`) to verify
it doesn't crash and the output looks reasonable. Don't write automated
assertions against the TDD project output since there's no way to validate
the JSON correctness without the `techDebtMetrics` pipeline being further
along. The script approach also allows adding more pipeline stages later.

---

## Q11: Interaction model for implementation

You described a strict iterative workflow: you tell the agent what to work on,
it returns a plan, you adjust, then you ask it to write code, and it stops for
review. A couple of clarifying questions about this:

- Should each "step" be roughly one feature/fix (e.g., "handle annotation
  types", "implement interface inheritance"), or do you prefer smaller or
  larger increments?
- When writing code changes, do you want the agent to also update/create the
  corresponding unit test and expected YAML file in the same step, or should
  tests be a separate step you request?

**A11:** Each step should be roughly one feature/fix, or a related set of
changes when features interact and must change together. The agent should
typically create unit test(s) alongside the code change in the same step.

However, if a complication arises, the user may specifically request an
integration test first (with or without expected YAML) to understand the
"shape of the constructs" before committing to the implementation approach.
The agent should be flexible and follow the user's lead on test strategy
per-step.

---

## Q12: Existing test failures

The current test `test0001/abstraction.yaml` has an intentional
`youShallNotPass: true` field that prevents it from passing. The existing
tests may also fail due to debug `println` statements or incomplete features.

Should fixing the existing tests be one of the first steps, or should we
focus on new functionality first and circle back to fix tests later?

**A12:** Focus on new functionality first. The `youShallNotPass` flag exists
because the test passes with current output but the output is incomplete —
it's a marker that more features are needed.

The agent may:
- Remove `println` debug statements (artifacts of exploration).
- Remove or fill out incomplete feature stubs.
- Add TODO comments and warning logs for features known to be needed but not
  yet implemented, so that "unknown feature" messages are replaced with
  clear indicators of planned work.
- Clean up code as it goes, as long as it doesn't break existing functionality.

Existing tests will naturally start passing as features are completed and
the expected YAML files are updated to match the now-complete output.

---

## Q13: Do you have any constraints on code style or patterns?

Looking at the existing codebase, I see some established patterns:
- `Factory<T>` + `Ref<T>` pattern for construct management.
- `Cmp` / `CmpOptions` for comparison-based deduplication.
- `Jsonable` interface with `toJson(JsonHelper)` for serialization.
- `Logger` with push/pop indentation for structured logging.

Should the agent strictly follow these existing patterns, or are you open to
introducing new patterns where they'd help (e.g., a visitor pattern for
Spoon element dispatch)?

**A13:** Open to new patterns where helpful. Context on existing patterns:

- **Jsonable**: Intentionally avoids YAML library dependency. The custom JSON
  parser supports a relaxed format (unquoted strings, comments, trailing
  commas) which makes test YAML files easier to write and read.
- **Logger**: Designed for human-readable output with indentation, NOT for
  structured logging tools like Logrus/Datadog/Splunk.
- **Cmp/CmpOptions**: Solves comparison of constructs that contain `Ref<T>`
  values — allows consistent ordering in Sets/Maps when refs aren't yet
  resolved, then switches to comparing resolved values. Needed for the
  consolidation phase.
- **Patterns loosely mirror Go abstractor** for maintainability across both
  codebases, but this is a preference not a hard requirement.

New patterns (e.g., visitor pattern for Spoon dispatch) are welcome if they
make the code cleaner or more maintainable.

---

---

## Q14: Code quality expectations

**A14:** (Provided directly by user.) The code must be debuggable and follow
good practices. However, the Java abstractor is NOT a long-lived project with
many developers, multiple versions, or years of maintenance ahead. It has a
finite lifetime — it's a research tool, not a product.

The code needs to be:
- **Good enough** to reach the success criteria (run on all 31 TDD projects,
  produce accurate output for TD analysis).
- **Readable enough** that other researchers can understand it.
- **Debuggable** — clear logging, sensible structure.

It does NOT need to be:
- Perfectly polished or over-engineered.
- Designed for extensibility beyond the research goals.
- Hardened for arbitrary Java projects outside the TDD scope.

---

## Requirements Clarification Complete

All key requirements have been discussed and recorded.
