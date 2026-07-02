# TODO

## Features and debugging

1. [ ] **Add more logs before and during consolidation.**
  The consolidation phase takes a very long time with the target projects.
  Print a short summary of the sizes before consolidation and some logs
  during so that we can see the progress.

2. [ ] **Improve logs during normal runs.**
   Add an config options to write the logs to a file. Allow simplified output
   to help us determine how much longer it will take to finish some work and
   limit the noise. This means we might have a mode where it shows percentage
   complete. Select message level to output and to push/pop at.
   Also add timing output that tells how long something took.

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

1. [x] **Constructor type parameters are ignored at the instance level** — resolved by the
  call-site `addMethodInstForCall(CtConstructorCall)` path in `Abstractor`. The object-walk
  `addMethodInstForObjectInst(CtConstructor, ...)` intentionally still leaves the ctor's own
  type params generic in the resolved signature (baseline coverage); call-site MethodInsts
  bind them with the actual arguments from `cc.getActualTypeArguments()`.

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

5. [x] **`addMethodDeclOrAbstract` / `addMethodDeclForConstructor` don't gate on shadow types** —
  both single-arg overloads now bail out with a `Logger.notice` when
  `decl.isShadow()` (methods) or `c.isShadow()` (constructors), so the
  `addDeclaration` path no longer materializes shadow JDK/library methods into
  `MethodDecl`s whose typeParams (stripped by Spoon at shadow load) mismatch
  their `MethodInst`'s instanceTypes. The two-arg overloads are unchanged
  because their callers (e.g. `addObjectInst`'s ctor loop) have already
  committed to the receiver. Was responsible for ~1046 of the 1047 validation
  errors on `commons-bcel`; longer-term the named-stub-`InterfaceDecl` plan
  (Step 5 in Behavior gaps) will replace `null` with real stubs.

6. [x] **`addSelection` creates a Selection with `origin = null` when the field's declaring
  type can't be resolved** — `addSelection` now resolves `origin` up front and
  returns `null` (with a `Logger.notice`) when the declaring type can't be
  tracked. `Analyzer.addRead` / `addWrite` are null-safe. Was responsible for
  1 of the 1047 validation errors on `commons-bcel`.
