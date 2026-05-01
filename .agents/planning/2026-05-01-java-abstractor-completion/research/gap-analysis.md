# Research: Gap Analysis

## What the Abstractor Needs to Produce

Per `genFeatureDef.md`, a valid output JSON must contain:

- `language`: "java" ✅
- `locs`: location map ✅
- `abstracts`: interface method signatures ✅
- `arguments`: method parameters/results ✅
- `basics`: primitive types ✅
- `fields`: struct/class fields ✅
- `interfaceDecls`: named interface declarations ⚠️ Partial
- `interfaceDescs`: interface type descriptions ⚠️ Partial
- `interfaceInsts`: generic interface instantiations ⚠️ Only arrays
- `methods`: method declarations ✅
- `methodInsts`: generic method instantiations ❌ Not populated
- `metrics`: method metrics ⚠️ Partial (usage tracking incomplete)
- `objects`: class/object declarations ✅
- `objectInsts`: generic object instantiations ❌ Not populated
- `packages`: package declarations ⚠️ Imports missing
- `selections`: field/method selections ✅
- `signatures`: method signatures ✅
- `structDescs`: structure descriptions ✅
- `typeParams`: type parameters ⚠️ Partial
- `values`: package-level variables ❌ Not populated

## Priority-Ordered Gap List

### P0: Must Fix for Basic Functionality

These will cause crashes or fundamentally incorrect output on real projects:

1. **Error handling / robustness**: The abstractor throws exceptions on many
   unhandled cases (unknown type descriptors, null declarations from shadow
   types, etc.). Real projects will have annotation types, anonymous classes,
   lambda expressions, and other constructs not currently handled in
   `addTypeDesc()`. The abstractor needs to gracefully handle or skip these.

2. **Package imports**: The `getImports()` method is a debug stub. Package
   import relationships are needed for the downstream metrics pipeline.

3. **Interface inheritance**: `getSuperInterfaces()` for both `CtClass` and
   `CtInterface` are not connected. This means the `inherits` field on
   `InterfaceDesc` is always empty.

4. **Class super-interfaces**: `addObjectDecl` doesn't connect
   `c.getSuperInterfaces()` to the object's synthesized interface description.

5. **Enum completion**: Enums are partially handled but enum constant values
   are not added to packages.

### P1: Needed for Correct Output

6. **Values**: Package-level static fields and constants need to be extracted.

7. **Interface pinning**: `InterfaceDesc.pin` is not set for declared interfaces.

8. **Nested type handling**: Both nested classes and nested interfaces need
   proper scoping (differentiating `Outer.Inner` from `Other.Inner`).

9. **Generic instantiation tracking**: `ObjectInst`, `MethodInst`, and real
   `InterfaceInst` (beyond arrays) are not populated.

### P2: Needed for Accurate Metrics

10. **Metrics: assignment usage** (`addAssignmentUsage`): Writes are not tracked.

11. **Metrics: executable reference usage** (`addExecutableReferenceUsage`):
    Invocations through references not tracked.

12. **Constructor flag**: MethodDecl doesn't distinguish constructors from methods.

### P3: Cleanup

13. Remove debug `println` statements from `getImports()`.
14. Set `logElementTree` and `logUsage` to `false` or make configurable.
15. Fix `test0001/abstraction.yaml` to remove `youShallNotPass`.
16. Cross-connection: Add interface declarations and values to packages.

## Research Context from Dissertation

The dissertation focuses on **participation** as a fuzzy estimate of membership
for TD analysis. The key metrics needed from the abstractor are:

- **Membership**: Which methods belong to which classes (objects). ✅ Available
  via `ObjectDecl.methods`.
- **Metrics per method**: Complexity, line count, code count, indents. ⚠️ Partial.
- **Invocations**: What methods call what other methods. ⚠️ Partial.
- **Reads/Writes**: What types (fields) are accessed by each method. ⚠️ Partial.
- **Package structure**: Import dependencies between packages. ❌ Missing.

The participation matrix requires knowing which objects are **accessed** (read/written)
by each method. This means the metrics reads/writes tracking is especially important
for the research goals.

## Suggested Implementation Order

Given the research goals and the need to work on real TDD projects:

1. **Robustness first**: Handle all type descriptor cases gracefully so the
   abstractor doesn't crash on real projects (annotations, anonymous classes,
   lambdas, shadow types, wildcards, etc.).
2. **Package imports**: Complete import tracking.
3. **Inheritance**: Connect interface/class inheritance chains.
4. **Metrics completeness**: Finish reads/writes/invocations tracking.
5. **Values and enums**: Extract package-level constants and enum values.
6. **Generic instances**: Track real generic instantiations.
7. **Cleanup**: Remove debug code, fix tests.
