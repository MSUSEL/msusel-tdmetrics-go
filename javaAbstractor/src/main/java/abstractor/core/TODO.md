# TODO

## Features and debugging

1. [ ] **Add more logs before and during consolidation.**
  The consolidation phase takes a very long time with the target projects.
  Print a short summary of the sizes before consolidation and some logs
  during so that we can see the progress.

2. [ ] **Try to do post write validation and/or validation of yaml.**
   The validation is failing on the target projects but it would be better
   to be able to review the yaml when trying to handle the validation problems.

3. [ ] **Improve logs during normal runs.**
   Add an config options to write the logs to a file. Allow simplified output
   to help us determine how much longer it will take to finish some work and
   limit the noise. This means we might have a mode where it shows percentage
   complete. Select message level to output and to push/pop at.
   Also add timing output that tells how long something took.

4. [ ] **Add better file/line outputs for debugging.** It appears some of the
   output isn't showing the file path but only the line numbers.

5. [ ] **Improve validation errors output.** Add an optional limit to the number
   of errors. Add a unique number for the different types of checks so that they
   can be laster grouped by how many of one type is failing.

## Spoon usage caveats / corrections

1. [ ] **`CtType.getReference()` strips formal type parameters.** Discovered while
  fixing test1006: the returned reference has empty `getActualTypeArguments()` and
  a `getDeclaringType()` that itself has no type args. Anywhere that needs the
  formal chain must use `SpoonUtils.parameterizedRef(type)`. Candidates worth auditing
  for the same trap: `addStructDesc` (when constructing references for `$super`/`$nest`),
  and any future code that walks `getNestedTypes()` and calls `nt.getReference()`.

2. [ ] **`CtConstructor.getType()` returns a raw reference to the declaring type.**
  Same lesson as above and the bug we just fixed — for the constructor result type
  use `parameterizedRef(m.getDeclaringType())`, not `m.getType()`.

3. [ ] **Synthetic references built by `parameterizedRef` have no AST parent.**
  That means `tr.hasParent(c.getParent())` (the `definedInNest` check in
  `addObjectInst`/`addInterfaceInst`) will return `false` for them. Today this is
  benign because the frame `nestCount` ends up at 0 in the cases we hit, but a
  future caller that depends on `definedInNest` being correct for a synthetic ref
  will be surprised. Consider deriving `definedInNest` from
  `c.getDeclaringType() != null && c.isStatic() == false` instead of `tr.hasParent(...)`.

4. [ ] **`CtTypeParameter.getTypeErasure()` only returns the first bound** for multi-bounded
  type params (`T extends A & B` → just `A`). Already noted in code; agree it's a real
  correctness gap. Spoon does not expose multi-bound directly
  here; `tp.getSuperclass()` plus `tp.getSuperInterfaces()` (or
  `tp.getReference().getBoundingType()` after walking) gives the full list.

5. [ ] **`CtWildcardReference.getBoundingType()` has the same single-bound limitation**
  as `getTypeErasure`. `addWildcard` currently only handles `? extends Foo` / `? super Bar`
  with a single bound. `? extends A & B` will silently drop `B`.

6. [ ] **`getAllMethods()` pulls inherited methods from JDK shadow super-interfaces.**
  `addInterfaceDesc` iterates `i.getAllMethods()` and then filters via `isObjectMethod`;
  for any interface that extends `java.util.Map` (or similar) this pulls dozens of abstracts
  and signatures into the project, as seen in test1005. The cheapest mitigation is
  `i.getMethods()` for the declared set plus an explicit walk of declared `getSuperInterfaces()`.
  The current approach is workable but produces noisy output and slow tests.

## Behavior gaps

1. [ ] **Constructor type parameters are ignored at the instance level.** `addMethodDeclForConstructor`
  populates `typeParams` from `addTypeParams(ctor)`, but `addMethodInstForObjectInst(CtConstructor, ...)`
  only varies on the receiver's `instanceTypes`. A ctor with its own type params
  (`public <U> Bar(U u)`) cannot currently produce per-call ctor instantiations.

2. [ ] **`addObjectInst` calls `synthesizeObjectInterface(c, null)` in its supplier.**
  That passes `null` for the pin, so the InterfaceDesc has no `pin`, while the
  ObjectDecl path passes the decl ref. After consolidation these will collapse
  to the same InterfaceDesc only if the abstracts match exactly.
  Inconsistent — either both should pin, or the comment in `synthesizeObjectInterface`
  should explain why ObjectInsts intentionally lose the pin.

3. [ ] **`addStructDesc` doesn't include `$super` chain types.** Only `getSuperclass()`
  is added as a single `$super` field. A class with both `extends` and `implements`
  ignores the interface side here (those are handled via `synthesizeObjectInterface.inherits`),
  but the lack of any link in the data view means TD metrics computed from struct data
  alone will miss inherited fields. Possibly intentional; document if so.

4. [ ] **Package imports are not derived yet.** Already noted in `Abstractor.performAbstraction`'s
  comment; agree it's still pending and ties into the planned Resolver pipeline.

5. [ ] **`Validator` runs after `crossConnectConstructs`, but `crossConnectConstructs`**
  **itself does no nesting checks.** If `obj.setNest(...)` was missed for a nested type,
  the validator may not catch it. Worth adding a validator pass that confirms
  every `obj.nest != null` iff `obj.getDeclaringType() instanceof CtType`.
