# Spoon Notes

To abstract Java I am using Spoon.

- [Spoon on gforge](https://spoon.gforge.inria.fr/)
- [Spoon on Github](https://github.com/INRIA/spoon)

## About Generics

- Spoon calls instances of a generic type, references.
  - I assume they call them references since instances only exist when a
    generic type is referenced with a type argument.

## About Methods

- There are 3 kinds of functions all called executables:
  - a method
  - a constructor
  - an anonymous block
- Method is in a collection of the same method overridden another.
- Signatures in Java are different that in Go. See Comment about Signature in
  [CtExecutable](https://github.com/INRIA/spoon/blob/master/src/main/java/spoon/reflect/declaration/CtExecutable.java#L122-L141)
  - > Note also that the signature of a method reference is the same as
      the signature of the corresponding method if and only if the method
      parameters does not involve generics in their types.
      Otherwise, one has e.g. `m(String)` (reference) and `m(T)` (declaration)
  - > In the Java programming language, a method signature is the method name
      and the number and type of its parameters.
      Return types and thrown exceptions are not considered
      to be a part of the method signature.

      see [Stackoverflow](https://stackoverflow.com/questions/16149285/does-a-methods-signature-in-java-include-its-return-type)
      and [Wikipedia](https://en.wikipedia.org/wiki/Type_signature)
- [CtMethod Uses](https://spoon.gforge.inria.fr/mvnsites/spoon-core/apidocs/spoon/reflect/declaration/class-use/CtMethod.html)
- [CtMethodImpl](https://github.com/INRIA/spoon/blob/master/src/main/java/spoon/support/reflect/declaration/CtMethodImpl.java)

## About Shadows

What Spoon means by "shadow". A CtType/CtElement is shadow when Spoon built it
from bytecode/reflection instead of source. It's not a value judgment — it just
means "this type is visible in the classpath but wasn't compiled from the sources
we handed to the model". In practice for BCEL: everything under `java.*`, `javax.*`,
everything from Maven-transitive jars, and anything else outside the -i root
all comes back as shadow. Project types, even ones in different packages under
the same source tree, are non-shadow.
