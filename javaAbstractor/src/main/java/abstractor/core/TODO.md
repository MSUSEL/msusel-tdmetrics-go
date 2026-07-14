# TODO

## Spoon usage caveats / corrections

1. [ ] **`CtType.getReference()` strips formal type parameters.** Discovered while
  fixing test1006: the returned reference has empty `getActualTypeArguments()` and
  a `getDeclaringType()` that itself has no type args. Anywhere that needs the
  formal chain must use `SpoonUtils.parameterizedRef(type)`.
  Triage of the current `nt.getReference()` / `nest.getReference()` call sites:
    - Lines 199, 232, 898, 996, 1161 (the nested-type registration loops):
      benign. They pass the raw ref through `addTypeDesc` → `addObjectInst`,
      which hits the "no actual type args" branch and returns the generic
      decl. Storing the decl in `nestedTypes` is what we want here.
    - Line 265 (`addInterfaceDesc` pin when the parent is a `CtType<?>`) and
      line 701 (`addStructDesc`'s `$nest` field): these DO care about the
      formal chain. A nested type inside `Outer<T>` currently gets a raw
      `Outer` reference instead of `Outer<T>`. Switch these two to
      `SpoonUtils.parameterizedRef(nest)`.
    - `addStructDesc`'s `$super` (line 693) uses `c.getSuperclass()`, which
      is a real Spoon API returning the parameterized ref directly — fine.

2. [ ] **Synthetic references built by `parameterizedRef` have no AST parent.**
  That means `tr.hasParent(c.getParent())` (the `definedInNest` check in
  `addObjectInst`/`addInterfaceInst`) will return `false` for them.
  Verification: Spoon's `CtType.getReference()` returns a fresh reference
  via `getFactory().Type().createReference(this)` whose parent is
  UNINITIALIZED (reads as null). `parameterizedRef` then only wires up
  `setDeclaringType` (a synthetic ref chain), never the real AST parent —
  so the concern is real for any synthetic ref, not just after our
  adjustments.
  The only known caller that hands a synthetic ref to `inSameNested` today
  is `getReceiverForCall`'s fallback branch (`in.getTarget() == null`),
  which fires in static-like call sites where outer type-param bindings
  don't apply, so the practical impact is limited right now. Still,
  `pushFrame` vs `pushCleanFrame` will silently pick the wrong one for any
  future caller that hands a synthetic ref in from a generic nested
  context. Cheapest replacement is
  `t.getDeclaringType() != null && !t.isStatic()` (interfaces are
  implicitly static so they always pick `pushCleanFrame`, which is correct
  because a nested interface never inherits the outer type's parameters).

3. [ ] **`CtTypeParameter.getTypeErasure()` only returns the first bound** for
  multi-bounded type params (`T extends A & B` → just `A`). Already noted
  inline at `addTypeParam` line 837–838. Spoon does not expose multi-bound
  directly here; `tp.getSuperclass()` plus `tp.getSuperInterfaces()` (or
  walking `tp.getReference().getBoundingType()`) gives the full list.

## Behavior gaps

1. [ ] **`addStructDesc` doesn't include `$super` chain types.** Only `getSuperclass()`
  is added as a single `$super` field. A class with both `extends` and `implements`
  ignores the interface side here (those are handled via `synthesizeObjectInterface.inherits`),
  but the lack of any link in the data view means TD metrics computed from struct data
  alone will miss inherited fields. Possibly intentional; document if so.

## Additional cleanup candidates

1. [ ] **`addDeclaration` unchecked `(CtTypeReference<?>)ref` cast.** Lines
  122–123 (old line numbers) narrow with `elem instanceof CtReference ref` and then blind-cast
  to `CtTypeReference<?>` inside the CtClass / CtInterface branches. In
  practice only type references resolve to class/interface declarations, but
  a rogue `CtExecutableReference` or `CtFieldReference` reaching this point
  would throw `ClassCastException` instead of hitting the "unhandled decl"
  log. Narrow with `ref instanceof CtTypeReference<?> tr` for defensive
  parity with the rest of the file.
