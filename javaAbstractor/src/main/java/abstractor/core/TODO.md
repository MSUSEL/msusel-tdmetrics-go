# TODO

The following is a list of things the Agent thinks is still outstanding issues:

## Spoon usage caveats / corrections

- **`tr.getTypeDeclaration()` is not wrapped in try/catch in `Abstractor.addTypeDesc`** despite the AGENTS.md convention. For unresolved/shadow references this can throw (e.g. `SpoonClassNotFoundException`, `IllegalArgumentException` from internal class loading). Wrap it and fall back to `Baker.anyDesc()` (or `null` in declaration paths) with a logged warning.

- **`CtType.getReference()` strips formal type parameters.** Discovered while fixing test1006: the returned reference has empty `getActualTypeArguments()` and a `getDeclaringType()` that itself has no type args. Anywhere that needs the formal chain must use `SpoonUtils.parameterizedRef(type)`. Candidates worth auditing for the same trap: `addStructDesc` (when constructing references for `$super`/`$nest`), and any future code that walks `getNestedTypes()` and calls `nt.getReference()`.

- **`CtConstructor.getType()` returns a raw reference to the declaring type.** Same lesson as above and the bug we just fixed — for the constructor result type use `parameterizedRef(m.getDeclaringType())`, not `m.getType()`.

- **Synthetic references built by `parameterizedRef` have no AST parent.** That means `tr.hasParent(c.getParent())` (the `definedInNest` check in `addObjectInst`/`addInterfaceInst`) will return `false` for them. Today this is benign because the frame `nestCount` ends up at 0 in the cases we hit, but a future caller that depends on `definedInNest` being correct for a synthetic ref will be surprised. Consider deriving `definedInNest` from `c.getDeclaringType() != null && c.isStatic() == false` instead of `tr.hasParent(...)`.

- **`CtTypeParameter.getTypeErasure()` only returns the first bound** for multi-bounded type params (`T extends A & B` → just `A`). Already noted in code; agree it's a real correctness gap and matters for test1010. Spoon does not expose multi-bound directly here; `tp.getSuperclass()` plus `tp.getSuperInterfaces()` (or `tp.getReference().getBoundingType()` after walking) gives the full list.

- **`CtWildcardReference.getBoundingType()` has the same single-bound limitation** as `getTypeErasure`. `addWildcard` currently only handles `? extends Foo` / `? super Bar` with a single bound. `? extends A & B` (rare) will silently drop `B`.

- **`tpr.getDeclaration()` does not lazy-load shadow types** (we confirmed). But the symmetrically-named `tr.getTypeDeclaration()` does. The naming similarity is a footgun; worth a one-line code comment on each call site so future maintainers don't conflate them.

- **`getAllMethods()` pulls inherited methods from JDK shadow super-interfaces.** `addInterfaceDesc` iterates `i.getAllMethods()` and then filters via `isObjectMethod`; for any interface that extends `java.util.Map` (or similar) this pulls dozens of abstracts and signatures into the project, as seen in test1005. The cheapest mitigation is `i.getMethods()` for the declared set plus an explicit walk of declared `getSuperInterfaces()`. The current approach is workable but produces noisy output and slow tests.

## Factory / Instantiator gotchas

- **`Factory.elemInProg` is never removed from.** Despite its name, "in progress" effectively becomes "ever created". This currently only matters as a recursion guard, but the semantics are misleading. Either remove the entry in a `finally` block at the end of `create`, or rename it to `elemCreated` and document the intent.

- **`Instantiator.add(param, arg)` overwrites silently** when the same param is added twice (re-keyed in `paramOrder`, prior value replaced). For nested-class instantiations whose outer params shadow inner ones, this is the right behavior. For accidental double-adds it hides bugs. Consider asserting that the second `add` for the same param either has the same arg or is in a different nest level.

- **`setToCompareResolved` clears `nonElemRef`.** Refs added through `addOrGetRef` (synthesized constructs, not tied to a CtElement) lose their entry. If anything still holds those refs outside the factory, they become orphaned. Worth a code comment explaining why this is safe (because consolidation will fix them).

## Behavior gaps

- **Anonymous classes and lambdas are *skipped*, not folded.** AGENTS.md says they should fold into the enclosing method's metrics, but `addObjectDecl` / `addInterfaceDecl` short-circuit on `isAnonymous()` / `isLocalType()` with `return null`. The contained code's reads/writes/invokes are never analyzed. Likely fix: visit the body via `Analyzer` even when the type itself is skipped.

- **Constructor type parameters are ignored at the instance level.** `addMethodDeclForConstructor` populates `typeParams` from `addTypeParams(ctor)`, but `addMethodInstForObjectInst(CtConstructor, ...)` only varies on the receiver's `instanceTypes`. A ctor with its own type params (`public <U> Bar(U u)`) cannot currently produce per-call ctor instantiations.

- **`addObjectInst` calls `synthesizeObjectInterface(c, null)` in its supplier.** That passes `null` for the pin, so the InterfaceDesc has no `pin`, while the ObjectDecl path passes the decl ref. After consolidation these will collapse to the same InterfaceDesc only if the abstracts match exactly. Inconsistent — either both should pin, or the comment in `synthesizeObjectInterface` should explain why ObjectInsts intentionally lose the pin.

- **`addStructDesc` doesn't include `$super` chain types.** Only `getSuperclass()` is added as a single `$super` field. A class with both `extends` and `implements` ignores the interface side here (those are handled via `synthesizeObjectInterface.inherits`), but the lack of any link in the data view means TD metrics computed from struct data alone will miss inherited fields. Possibly intentional; document if so.

- **Package imports are not derived yet.** Already noted in `Abstractor.performAbstraction`'s comment; agree it's still pending and ties into the planned Resolver pipeline.

- **`Validator` runs after `crossConnectConstructs`, but `crossConnectConstructs` itself does no nesting checks.** If `obj.setNest(...)` was missed for a nested type, the validator may not catch it. Worth adding a validator pass that confirms every `obj.nest != null` iff `obj.getDeclaringType() instanceof CtType`.

## Minor cleanups

- **`isObjectMethod` does an O(n) signature-string compare against every `java.lang.Object` method on every call.** Cache the set of object signatures once in `SpoonUtils` (or pre-compute on `Baker`).

- **`SpoonUtils.describeElem` falls through three try/catches** to handle different element kinds. Fine, but the fall-through order (`CtNamedElement` → `CtTypeInformation` → `CtExecutable` → `CtReference`) silently swallows real errors via `catch (Exception ignore)`. Consider logging at debug level so unexpected failures aren't invisible.

- **`addObjectInst`/`addInterfaceInst` push a frame even when `typeArgs == null`.** Actually they early-return before pushing — but the early-return path through `decl` skips the finisher loops in `addObjectInst`. If a future change makes the early return depend on partial state, it could leak. Worth either commenting the invariant or restructuring so push/pop bracket the entire body uniformly.
