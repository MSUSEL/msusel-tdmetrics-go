# TODO

## Spoon usage caveats / corrections

1. [ ] **`CtType.getReference()` strips formal type parameters.** Discovered while
  fixing test1006: the returned reference has empty `getActualTypeArguments()` and
  a `getDeclaringType()` that itself has no type args. Anywhere that needs the
  formal chain must use `SpoonUtils.parameterizedRef(type)`. Candidates worth auditing
  for the same trap: `addStructDesc` (when constructing references for `$super`/`$nest`),
  and any future code that walks `getNestedTypes()` and calls `nt.getReference()`.
  Update: nt.getReference() still appears at lines 199, 232, 265, 701, 898, 996, 1161.
  See "Problematic" below for which of those actually matter.

2. [ ] **Synthetic references built by `parameterizedRef` have no AST parent.**
  That means `tr.hasParent(c.getParent())` (the `definedInNest` check in
  `addObjectInst`/`addInterfaceInst`) will return `false` for them. Today this is
  benign because the frame `nestCount` ends up at 0 in the cases we hit, but a
  future caller that depends on `definedInNest` being correct for a synthetic ref
  will be surprised. Consider deriving `definedInNest` from
  `c.getDeclaringType() != null && c.isStatic() == false` instead of `tr.hasParent(...)`.
  Before doing this change, ensure that `getDeclaringType` doesn't cause any lazy loading.
  Also, determine if this is actually a problem since `parameterizedRef` may carry some
  parent information (I don't know how the reference is created and what it carries
  prior to me adjusting it).

3. [ ] **`CtTypeParameter.getTypeErasure()` only returns the first bound** for multi-bounded
  type params (`T extends A & B` → just `A`). Already noted in code; agree it's a real
  correctness gap. Spoon does not expose multi-bound directly
  here; `tp.getSuperclass()` plus `tp.getSuperInterfaces()` (or
  `tp.getReference().getBoundingType()` after walking) gives the full list.
  Update:  addTypeParam still uses tp.getTypeErasure() at
  line 838 (the inline TODO on line 837 is still there).

4. [ ] **`CtWildcardReference.getBoundingType()` has the same single-bound limitation**
  as `getTypeErasure`. `addWildcard` seems to currently only handle `? extends Foo` / `? super Bar`
  with a single bound. `? extends A & B` will silently drop `B`. Re-evaluate and
  determine how to use spoon to fix the issue if it still exists.

5. [ ] **`getAllMethods()` pulls inherited methods from JDK shadow super-interfaces.**
  `addInterfaceDesc` iterates `i.getAllMethods()` and then filters via `isObjectMethod`;
  for any interface that extends `java.util.Map` (or similar) this pulls dozens of abstracts
  and signatures into the project, as seen in test1005. The cheapest mitigation is
  `i.getMethods()` for the declared set plus an explicit walk of declared `getSuperInterfaces()`.
  The current approach is workable but produces noisy output and slow tests.

## Behavior gaps

1. [ ] **`addObjectInst` calls `synthesizeObjectInterface(c, null)` in its supplier.**
  That passes `null` for the pin, so the InterfaceDesc has no `pin`, while the
  ObjectDecl path passes the decl ref. After consolidation these will collapse
  to the same InterfaceDesc only if the abstracts match exactly.
  It is not an issue to be inconsistent (both do not need to pin), should add a
  comment in `synthesizeObjectInterface` should explain why ObjectInsts intentionally
  lose the pin. Pins should only be added when there are unexported (private) methods
  in the interface meaning it may only be used within the current package, otherwise
  the pin can be ignored, however, we need to check that inheritance is being used
  to differentiate so that they remain different if constructed differently since
  Java doesn't have duck-typing like Go does.

2. [ ] **`addStructDesc` doesn't include `$super` chain types.** Only `getSuperclass()`
  is added as a single `$super` field. A class with both `extends` and `implements`
  ignores the interface side here (those are handled via `synthesizeObjectInterface.inherits`),
  but the lack of any link in the data view means TD metrics computed from struct data
  alone will miss inherited fields. Possibly intentional; document if so.
