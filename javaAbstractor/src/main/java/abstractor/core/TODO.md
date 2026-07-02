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

5. [ ] **`addMethodDeclOrAbstract` / `addMethodDeclForConstructor` don't gate on shadow types.**
  When analyzer sees a call to a JDK / third-party method (e.g. `x.getClass().getName()`,
  `throw new RuntimeException(...)`), `addDeclaration` routes through
  `addMethodDeclOrAbstract` which unconditionally does `addObjectDecl((CtClass<?>)decl)` for
  the shadow declaring class, then builds a `MethodDecl` for the shadow method. Shadow decls
  loaded from `.class` reflectively have `getFormalCtTypeParameters()` return an empty list
  even when the source declares type params (e.g. `Class<T>`), so the resulting `MethodDecl`
  has `typeParams: []`. When any generic instantiation of that class walks its methods, the
  produced `MethodInst` carries the class-level `instanceTypes`, and the validator rejects it
  (`[0600]`/`[0610]`/`[0620]`). Simplest fix: bail out to `anyDesc()` (or `null`) in both
  methods when `decl.isShadow()`; longer-term this ties into Step 5's named-stub-`InterfaceDecl`s
  target. Was responsible for ~1046 of the 1047 validation errors on `commons-bcel`.
  Note: three attempts to reproduce with a small single-file test (`this.getClass().getName()`,
  `throw new RuntimeException` inside `<T>` method, generic wrapper class using `ArrayList<T>` +
  `IndexOutOfBoundsException`) all produced clean runs after the call-site MethodInst work,
  so the specific bcel repro pattern isn't obvious. Keeping the guard planned since the code
  path that made the bad `MethodDecl`s is still reachable and bcel exercises it.

6. [ ] **`addSelection` creates a Selection with `origin = null` when the field's declaring
  type can't be resolved.** `addDeclaration(field.getDeclaringType())` legitimately returns
  `null` for cross-artifact refs. The current code stores that null and the validator flags it
  (`[0800]`). Simplest fix: have `addSelection` return `null` when `origin` is null; callers
  already handle null refs via `addRead`/`addWrite`. Was responsible for 1 of the 1047
  validation errors on `commons-bcel`.
