# Implementation Plan: Java Abstractor Completion

## Checklist

- [ ] Step 1: Robustness — Type dispatch hardening
- [ ] Step 2: External type stubs and boxing
- [ ] Step 3: Enum completion
- [ ] Step 4: Values (package-level variables and constants)
- [ ] Step 5: Named nested class support
- [ ] Step 6: Anonymous class and lambda folding
- [ ] Step 7: Interface inheritance and super-class connections
- [ ] Step 8: Interface pinning and generation cleanup
- [ ] Step 9: Metrics — complete reads/writes/invokes tracking
- [ ] Step 10: Generic instantiation tracking (ObjectInst, InterfaceInst)
- [ ] Step 11: Generic instantiation tracking (MethodInst)
- [ ] Step 12: Resolver pipeline extraction
- [ ] Step 13: Package imports from type usage
- [ ] Step 14: Cross-connection and cleanup
- [ ] Step 15: TDD project validation script

---

## Step 1: Robustness — Type dispatch hardening

**Objective:** Make `addTypeDesc` and `addDeclaration` handle all Spoon type
cases without crashing, so the abstractor can run on real-world projects.

**Guidance:**
- Change `addTypeDesc` to use `tr.getTypeDeclaration()` instead of
  `tr.getDeclaration()` so shadow types are available rather than null.
- Add handling for `CtAnnotationType` → return `baker.objectDesc()` with a
  log notice (not warning, since this is expected behavior).
- Add handling for `CtWildcardReference` → map to the bounding type if
  upper-bounded, or `baker.objectDesc()` if unbounded.
- Add handling for anonymous types (`tr.isAnonymous()`) → return the
  superclass or implemented interface type.
- Add handling for local types (`tr.isLocalType()`) → same as anonymous.
- Wrap the body of `addTypeDesc` in a try/catch that logs a warning and
  returns `baker.objectDesc()` on any unexpected exception.
- Similarly wrap `addDeclaration` to log and return null on failure.
- Remove or replace the `unknownTypeDesc` method that throws an exception
  with one that logs a warning and returns `baker.objectDesc()`.

**Tests:**
- test1004: A class that uses annotation types, wildcards (`List<?>`),
  and references to JDK types. Verify it produces output without crashing
  and unhandled types map to `objectDesc` or appropriate stubs.

**Integration with previous work:** This is the foundation — all subsequent
steps depend on the abstractor not crashing.

**Demo:** Run `mvn test` and the new test passes. Run the abstractor on a
simple Maven project that uses annotations and generics without crashing.

---

## Step 2: External type stubs and boxing

**Objective:** External (JDK/library) types appear as named stubs in the
output instead of being collapsed to `Object`.

**Guidance:**
- Add a boxing map to `Baker`: `java.lang.Integer` → `int`,
  `java.lang.String` → basic `string`, `java.lang.Boolean` → `boolean`,
  `java.lang.Long` → `long`, `java.lang.Double` → `double`,
  `java.lang.Float` → `float`, `java.lang.Character` → `char`,
  `java.lang.Byte` → `byte`, `java.lang.Short` → `short`.
- Add `addExternalStub(CtTypeReference)` method to `Abstractor`:
  - If the type is a known boxed type, return the corresponding basic.
  - Otherwise, create a stub `InterfaceDecl` with the type's simple name
    in a package matching the external type's package path.
  - Cache stubs by qualified name to avoid duplicates.
- If the external type reference has actual type arguments (parameterized),
  create an `InterfaceInst` pointing to the stub declaration with the
    resolved type arguments.
- Route `addTypeDesc` to `addExternalStub` when `ty.isShadow()` or
  `ty == null` (after `getTypeDeclaration()`).

**Tests:**
- test1005: A class with fields of type `String`, `Integer`, `List<String>`,
  and a method returning `Map<String, Integer>`. Verify `String` becomes
  basic `string`, `Integer` becomes basic `int`, `List` and `Map` appear
  as named interface declaration stubs.

**Integration with previous work:** Builds on Step 1's dispatch changes.

**Demo:** Run `mvn test`, inspect output JSON showing named external types
instead of empty interfaces.

---

## Step 3: Enum completion

**Objective:** Enums are fully modeled as objects with their constants
extracted as values.

**Guidance:**
- Update `addEnum` to treat enums as `ObjectDecl` (similar to classes).
  Use the same `addObjectDecl` path since `CtEnum extends CtClass`, but
  handle the struct description differently: include the enum's fields
  plus a synthetic `$value` field of the enum's superclass type.
- Add enum constants as `Value` constructs in the enclosing package.
  Each enum constant is `const = true`, type = the enum's `ObjectDecl`.
- Add enum methods (user-defined methods, not compiler-generated ones
  like `values()` and `valueOf()`) via the normal `addMethod` path.
- Connect enum super-interfaces via `e.getSuperInterfaces()`.
- Handle `addDeclaration` for `CtEnum` elements (check before `CtClass`).

**Tests:**
- test1006: An enum with constants, fields, a constructor, and a method.
  Verify the enum appears as an object with struct, methods, and its
  constants appear as values.

**Integration with previous work:** Uses Step 1's robust dispatch and
Step 2's external stubs for any external types referenced by the enum.

**Demo:** Run `mvn test`, inspect output showing enum as object with
constant values.

---

## Step 4: Values (package-level variables and constants)

**Objective:** Package-level static fields and constants are extracted as
`Value` constructs.

**Guidance:**
- In `addObjectDecl`, after processing fields and methods, iterate over
  static fields that are `public` or package-visible.
- For each qualifying static field, create a `Value` construct:
  - `name` = field's simple name
  - `type` = `addTypeDesc(field.getType())`
  - `const` = `field.isFinal()`
  - `package` = the enclosing package
  - `loc` = field's source position
- Add the value to the package's `values` set.
- If the static field has an initializer expression, create `Metrics`
  for it (matching the Go abstractor's behavior for value initializers).

**Tests:**
- test1007: A class with `public static final int MAX = 100;` and
  `public static String name = "test";`. Verify both appear as values
  with correct types and const flags.

**Integration with previous work:** Uses Step 2's type resolution for
field types.

**Demo:** Run `mvn test`, inspect output showing values with correct
types, const flags, and packages.

---

## Step 5: Named nested class support

**Objective:** Named inner/nested classes are modeled as separate
`ObjectDecl` with the `nest` field set.

**Guidance:**
- In `addObjectDecl`, when `c.getRoleInParent() == CtRole.NESTED_TYPE`
  and `!c.isAnonymous()` and `!c.isLocalType()`:
  - Determine the enclosing context: if the parent is a `CtClass`, the
    nest target is the enclosing object's method or the object itself.
  - Set `obj.nest` to the enclosing construct's ref.
- Add `nest` field to `ObjectDecl` if not already present.
- Update `ObjectDecl.toJson()` to output `nest` when set.
- Update `ObjectDecl.getCmp()` to include `nest` in comparison so that
  `Outer.Inner` is distinct from `Other.Inner`.
- Named nested interfaces: similar treatment for `addInterfaceDecl` when
  the interface has `NESTED_TYPE` role.

**Tests:**
- test1008: Outer class with a named inner class. Inner class has fields
  and methods. Verify the inner class appears as a separate object with
  `nest` set, and its methods are in the output.

**Integration with previous work:** Uses Step 1's dispatch. The struct
for the nested class already adds a `$nest` field for the parent reference.

**Demo:** Run `mvn test`, inspect output showing nested object with
`nest` field pointing to enclosing method or object.

---

## Step 6: Anonymous class and lambda folding

**Objective:** Anonymous classes and lambdas are NOT separate objects;
their code is attributed to the enclosing method's metrics.

**Guidance:**
- In `addTypeDesc`, when `tr.isAnonymous()`:
  - Return the type of the superclass or implemented interface
    (the anonymous class's "shape"), not a new object declaration.
- In `addDeclaration`, skip anonymous and local classes.
- In `Analyzer.addElement`, when encountering `CtNewClass` (anonymous
  class instantiation):
  - Continue walking into its children (the anonymous class body).
  - This means the anonymous class's method bodies contribute to the
    enclosing method's complexity, reads, writes, and invocations.
- In `Analyzer.addElement`, when encountering `CtLambda`:
  - Walk the lambda body (block or expression) as children.
  - Track reads/writes/invocations from the lambda.
- Verify that `addComplexity` correctly counts control flow in anonymous
  class methods and lambda bodies.

**Tests:**
- test1009: A method containing an anonymous `Runnable` implementation
  and a lambda. Verify no separate objects are created for them, and the
  enclosing method's metrics include their complexity and invocations.

**Integration with previous work:** Depends on Step 1 (anonymous type
handling) and Step 5 (to verify named nested classes are NOT folded).

**Demo:** Run `mvn test`, verify metrics on enclosing method include
anonymous/lambda code, and no spurious objects appear.

---

## Step 7: Interface inheritance and super-class connections

**Objective:** `InterfaceDesc.inherits` is populated for both declared
interfaces and object-synthesized interfaces.

**Guidance:**
- In `addInterfaceDecl` finisher: iterate `i.getSuperInterfaces()`,
  call `addInterfaceDecl` or `addTypeDesc` for each, and add the
  resulting interface description to `id.inherits`.
- In `addObjectDecl` finisher: iterate `c.getSuperInterfaces()`,
  resolve each to an interface description, and add to the object's
  synthesized `InterfaceDesc.inherits`.
- For superclass (`c.getSuperclass()`): if the superclass is a project
  class (not `Object`), add the superclass's synthesized interface
  to `inherits` as well. This models Java's single inheritance.
- Handle diamond inheritance gracefully (interface implemented by both
  a class and its superclass).

**Tests:**
- test1010: An interface `A`, an interface `B extends A`, a class `C`
  implementing `B`, and a class `D extends C` implementing `A`.
  Verify `inherits` chains are correct and no duplicates.

**Integration with previous work:** Uses Step 2 for external interface
stubs when a class implements a JDK interface (e.g., `Serializable`).

**Demo:** Run `mvn test`, inspect output showing correct `inherits`
lists on interface descriptions.

---

## Step 8: Interface pinning and generation cleanup

**Objective:** `InterfaceDesc.pin` is set correctly, and interface
generation for objects is robust.

**Guidance:**
- In `addInterfaceDecl`: set `inter.pin` to the interface declaration
  ref (the `InterfaceDecl` itself). Already partially supported in
  `InterfaceDesc` constructor but not used for declared interfaces.
- In `addObjectDecl`: the synthesized `InterfaceDesc` for the class
  should have `pin` set to the `ObjectDecl` ref. This is already done
  in the existing code (`new InterfaceDesc(abstracts, ref)`).
- Move interface generation for objects to a later phase (in the
  Resolver) to ensure all methods have been added before synthesizing
  the interface. Currently it happens inline in `addObjectDecl`'s
  finisher which may miss methods added later.
- Ensure that constructor methods are NOT included in the synthesized
  interface (constructors are not abstracts).
- Ensure that static methods are NOT included in the synthesized
  interface.

**Tests:**
- test1011: A class with static methods, constructors, and instance
  methods. Verify the synthesized interface only contains instance
  method abstracts, and `pin` is set to the object.

**Integration with previous work:** Works with Step 7's inheritance.

**Demo:** Run `mvn test`, verify `pin` fields and correct abstract
lists in synthesized interfaces.

---

## Step 9: Metrics — complete reads/writes/invokes tracking

**Objective:** The `Analyzer` fully tracks invocations, field reads,
and field writes for accurate participation metrics.

**Guidance:**
- Complete `addAssignmentUsage`:
  - Get the LHS of the assignment.
  - If LHS is a `CtFieldWrite`, resolve the field declaration and
    create a `Selection`, add to `this.writes`.
  - If LHS is a `CtVariableWrite` to a local variable, skip (local
    variables don't affect participation).
- Complete `addExecutableReferenceUsage`:
  - Get the declaration from the executable reference.
  - If it's a `CtMethod`, add to `this.invokes` via `addDeclaration`.
  - If it's a `CtConstructor`, add to `this.invokes`.
- Handle `CtConstructorCall` in `addUsage`: add to invokes.
- Handle `CtExecutableReferenceExpression` (method references like
  `Foo::bar`): resolve and add to invokes.
- Handle `CtFieldWrite` in `addUsage`: resolve field, create Selection,
  add to writes.
- Remove `logElementTree` and `logUsage` hardcoded `true` values; tie
  them to `Config.verbose` or a constructor parameter.

**Tests:**
- test1012: A class with methods that read fields, write fields, call
  other methods, and use method references. Verify the metrics on each
  method show correct invokes, reads, and writes lists.

**Integration with previous work:** Uses Steps 1-2 for type resolution
of field types and method references.

**Demo:** Run `mvn test`, inspect metrics in output JSON showing
correct invokes/reads/writes with selection references.

---

## Step 10: Generic instantiation tracking (ObjectInst, InterfaceInst)

**Objective:** When a parameterized type is used (e.g., `HashMap<String, User>`),
create `ObjectInst` or `InterfaceInst` with resolved type arguments.

**Guidance:**
- In `addTypeDesc`, when processing a type reference that has actual type
  arguments (`tr.getActualTypeArguments()` is non-empty):
  - First resolve the base declaration (ObjectDecl or InterfaceDecl).
  - Collect the type argument refs via `addTypeDesc` on each argument.
  - If base is `ObjectDecl`: create `ObjectInst` with the generic, instance
    types, resolved struct desc (with concrete field types), and resolved
    interface desc.
  - If base is `InterfaceDecl`: create `InterfaceInst` with the generic,
    instance types, and resolved interface desc.
  - Register the instance on the generic declaration's `instances` list.
  - Return the instance ref as the type description.
- Handle wildcard type arguments: `? extends Foo` → use the bound type.
  `?` → use `baker.objectDesc()`.
- Cache instances by (generic + instance types) to avoid duplicates.

**Tests:**
- test1013: A generic class `Box<T>` with a field `T value` and method
  `T get()`. Usage: `Box<String>` and `Box<Integer>`. Verify two
  `ObjectInst` entries with correct resolved struct and interface.

**Integration with previous work:** Uses Steps 1-2 for type dispatch and
external stubs (type arguments may be external types).

**Demo:** Run `mvn test`, inspect output showing `objectInsts` and
`interfaceInsts` with correct generic back-references and instance types.

---

## Step 11: Generic instantiation tracking (MethodInst)

**Objective:** When a generic method is called with concrete type arguments,
create `MethodInst` with resolved signature.

**Guidance:**
- In `Analyzer.addUsage` for `CtInvocation`: check if the invoked method
  is generic and the call provides concrete type arguments.
  - `inv.getExecutable().getActualTypeArguments()` gives the type args.
  - If type args are present and not type parameters, create `MethodInst`:
    - `generic` = the generic `MethodDecl`
    - `instanceTypes` = resolved type argument refs
    - `resolved` = signature with concrete types substituted
- Also handle method instances from receiver instantiation: when an
  `ObjectInst` exists (e.g., `Box<String>`), each method on `Box<T>` needs
  a `MethodInst` for `Box<String>` with `T→String` resolved in the
  signature.
- Register on the generic method's `instances` list.

**Tests:**
- test1014: A generic method `<T> T identity(T x)` called as
  `identity("hello")` and `identity(42)`. A generic class `Pair<A,B>`
  with method `A getFirst()`, used as `Pair<String, Integer>`. Verify
  method instances are created with correct resolved signatures.

**Integration with previous work:** Depends on Step 10 (ObjectInst
triggers MethodInst creation for receiver methods).

**Demo:** Run `mvn test`, inspect `methodInsts` in output with correct
generic back-references and resolved signatures.

---

## Step 12: Resolver pipeline extraction

**Objective:** Extract post-processing logic from `Abstractor.finish()`
into a structured `Resolver` class.

**Guidance:**
- Create `abstractor/core/Resolver.java` with the pipeline:
  1. `expandGenericInstances()` — fixed-point loop (from Steps 10-11)
  2. `generateInterfaces()` — move from `addObjectDecl` finisher
  3. `computeInheritance()` — from Step 7
  4. `processMetrics()` — move from `processPendingMetrics()`
  5. `crossConnect()` — move from `crossConnectConstructs()`
  6. `consolidate()` — move from `consolidateCons()`
  7. `validate()` — move from `validate()`
- Update `Abstractor.finish()` to call `new Resolver(log, proj).resolve()`.
- The Resolver has access to the `Project` and `Logger`.
- Each step is a separate method for clarity and debuggability.
- The fixed-point loop runs `expandGenericInstances` and
  `resolveReferences` until no new constructs are created.

**Tests:**
- No new test cases — existing tests should continue to pass after
  the refactoring. Run full test suite.

**Integration with previous work:** This reorganizes Steps 7-11's work
into the correct pipeline order.

**Demo:** Run `mvn test`, all existing tests pass. The `finish()` method
delegates to the Resolver.

---

## Step 13: Package imports from type usage

**Objective:** Populate `Package.imports` by analyzing which packages
are actually referenced by each package's constructs.

**Guidance:**
- In the Resolver, after all constructs are finalized:
  - For each package P, iterate all its objects, interfaces, methods,
    and values.
  - For each type reference in their signatures, fields, and metrics:
    - Determine the type's package.
    - If that package is different from P, add it to P's imports.
  - Deduplicate and sort the imports list.
- Remove the old `getImports()` stub method and its debug println
  statements from `Abstractor.java`.
- This approach is more reliable than parsing `import` statements
  because it captures actual usage and handles wildcard imports.

**Tests:**
- test0002 (existing Maven test): Update the expected YAML to include
  import lists if the test uses types from multiple packages.
- Or add a test with two packages where one references the other.

**Integration with previous work:** Depends on Step 12 (runs in the
Resolver pipeline after all types are resolved).

**Demo:** Run `mvn test`, inspect package output showing correct
import lists.

---

## Step 14: Cross-connection and cleanup

**Objective:** All constructs are correctly wired to their packages,
debug artifacts are removed, and tests are updated.

**Guidance:**
- Ensure `crossConnect` in the Resolver adds to packages:
  - Methods (already partially done)
  - Interface declarations
  - Object declarations
  - Values
- Remove all debug `System.out.println` statements.
- Remove the `youShallNotPass` field from `test0001/abstraction.yaml`
  and update the expected output to match the now-complete abstraction.
- Update `test0002/abstraction.yaml` similarly.
- Update test1001-1003 expected YAMLs if the output has changed due
  to new features (external stubs, interface pinning, etc.).
- Set `logElementTree` and `logUsage` to be controlled by verbose flag
  rather than hardcoded.

**Tests:**
- All existing tests updated and passing.
- Run full `mvn test` suite.

**Integration with previous work:** Final polish step after all features.

**Demo:** `mvn test` passes with zero failures. Output JSON from any
test shows complete, correct data with no debug artifacts.

---

## Step 15: TDD project validation script

**Objective:** A script to run the abstractor against small TDD projects
and verify it doesn't crash.

**Guidance:**
- Create `javaAbstractor/scripts/run-tdd.sh`:
  - Accept a path to a cloned TDD project directory.
  - Build the abstractor jar if not already built.
  - Run the abstractor on the project.
  - Report success/failure and output file size.
  - Optionally run a basic sanity check (e.g., output is valid JSON,
    has non-zero objects and methods counts).
- Query `td_V2.db` to find the best commit per project (latest commit
  with SonarQube analysis data).
- Document which TDD projects to clone first for validation
  (recommend: `commons-cli`, `commons-exec`, `commons-dbutils` as they
  are small).
- Run the script on 2-3 small projects and document any remaining
  issues as TODOs for follow-up.

**Tests:**
- No automated tests. Manual validation.

**Integration with previous work:** This is the final validation that
all previous steps work together on real projects.

**Demo:** Run the script on `commons-cli`. The abstractor produces
a JSON file with objects, methods, metrics, and no crashes.
