# From Agent

Big win — from 147 down to 40 GT-missing. The remaining 40 split into two clean buckets:

## Bucket A

Java interfaces (≈34 of 40). These live in proj.InterfaceDecls, not proj.ObjectDecls. Everything in the list ending in Observer, Instruction, Consumer, Producer,
Thrower, Comparator, Reader, Visitor, Node, LoadClass, NamedAndTyped, Repository, Subroutine, InstructionContext, VerifierFactoryObserver, plus the constants-only interfaces
Constants, ExceptionConstants, InstructionConstants, ConstantObject — all interface in BCEL source.
Fix: union the two sets.

```Csharp
IEnumerable<IDeclaration> projDecls =
    proj.ObjectDecls.Cast<IDeclaration>().Concat(proj.InterfaceDecls);
```

## Bucket B

Nested-type name mismatches (6 of 40)

- classfile.ConstantUtf8$CACHE_HOLDER — nested class in a class; check whether Java abstractor emits it at all (it's a private static holder class for the string cache).

- generic.ConstantPoolGen$Index — same pattern.

- util.ClassPath$ClassFile, util.ClassPath$PathEntry, util.ClassPath$Zip — ClassPath in current BCEL is abstract class with a nested interface PathEntry and nested classes
  ClassFile/Zip. PathEntry will need step-3 (InterfaceDecl.FullName nest fix) once bucket A is unioned. ClassFile/Zip should already be in ObjectDecls with $ naming — if
  still missing, they're an actual emit gap.

- util.InstructionFinder$CodeConstraint — CodeConstraint is a nested interface. Blocked by the InterfaceDecl.FullName TODO — currently it would appear as
  util.CodeConstraint, not util.InstructionFinder$CodeConstraint.

So expect bucket A to shrink the miss set to ~6, and step 3 (nested-name in InterfaceDecl.FullName) plus one Java-side check on the private CACHE_HOLDER/Index holders to
eliminate the rest.

## The 23 proj-missing entries are on the GT side, not yours. Two sub-groups

- Java 9+ features (all 11 Module* + NestHost/NestMembers + Record/RecordComponentInfo + ModularRuntimeImage +
  ClassPath$Jar/JrtModule/JrtModules/Module/AbstractPathEntry/AbstractZip + LruCacheClassPathRepository + AbstractClassPathRepository + Args) — PMD/CK almost certainly
  skipped these; a common cause is a compile-target/JDK mismatch in the metrics pipeline, or the tool's parser bailing on module-info.java-adjacent code.

- Odd-ones-out (InvalidMethodSignatureException, ConstantUtf8$Cache) — worth a spot-check. Note the naming: your tool emits ConstantUtf8$Cache while GT has
  ConstantUtf8$CACHE_HOLDER. If both are supposed to be on commit 6ed18c5b…, one side is wrong — probably worth opening the file to see whether it's Cache or CACHE_HOLDER
  at that SHA and adjusting whichever tool has drifted.
