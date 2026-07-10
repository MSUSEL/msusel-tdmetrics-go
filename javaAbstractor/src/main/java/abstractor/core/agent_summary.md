# Agent shadow work

[Agent change](https://github.com/MSUSEL/msusel-tdmetrics-go/commit/b9616549cbcf09b40a1926cfc637f8097cd302f7#diff-dc0383eda13155e9f87c81061dce8ef06ddaa9e62492d8cffd997905335c1bca)

## Summary

Skips added (declarations only stub what's referenced):

- addInterfaceDecl (192‑196) — finisher: sets id.isShadow and returns before nested‑types walk.
- addInterfaceInst (231) — finisher: early return for shadow before nested‑types walk.
- addInterfaceDesc (247, 270) — creator skips the getAllMethods → addAbstract loop;
  finisher skips the super‑interfaces walk.
- addStructDesc (701) — early returns an empty StructDesc for shadows
  (skips fields, $super, $nest). TODO left about factory dedup of empty stubs.
- addObjectInst (1024‑1028) — finisher: early return for shadow before ctors/methods/nested walk.
  (ObjectInst itself has no isShadow; the shadow flag lives on its generic ObjectDecl.)

Synthesis commented out (deferred to your post‑abstraction pass):

- addObjectDecl (962), addObjectInst (1019 + resInterface = null), addEnum (1239).
  The synthesizeObjectInterface method is left intact and just becomes
  unused for now; feel free to delete/replace when you wire up the post pass.

On‑demand attachment:

- addAbstract (573‑585) — when the declaring interface is a shadow, attaches the
  returned Abstract ref to that interface's InterfaceDesc.abstracts.
  TODO left about the dedup risk you flagged.
- addSelection (792‑805) + helper attachShadowField (811‑817) — when the origin
  resolves to a shadow ObjectDecl or a shadow generic's ObjectInst, calls
  addField(field) and appends the field to the origin's StructDesc.fields
  (dedup by presence). TODO left about the same dedup risk.
- Method attachment for shadow objects — no change needed: addMethodDecl already
  does recv.methodDecls.add(ref), so shadow object methods accumulate naturally.
  The synthesized inter.abstracts for shadow objects is not touched here since synthesis is deferred.

Not changed (per your ask to skip synthesis rather than rework it):

- No new pending‑synthesis pass in performAbstraction.
- StructDesc / InterfaceDesc unchanged — no per‑owner discriminator added yet;
  the TODOs mark where empty stubs may collide.
