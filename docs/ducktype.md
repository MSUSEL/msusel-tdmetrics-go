# Duck type hunting

Go uses [duck typing](https://en.wikipedia.org/wiki/Duck_typing).
To be able to do design recovery we need to know how types interact.
We need to be able to determine what are the sub-types and super-types.

Go's type system will determine if
[a type implements another](https://pkg.go.dev/go/types#Implements).
Therefore, part of Go abstraction will require to determine
type relationships for all types to remove duck-typing.

This can either be $\frac{n(n-1)}{2}$ comparisons where $n$
is the number of types, or a tree of inheritance can be used
to perform less comparisons.

Assuming a type can be represented by a set of numbers,
the type $S$ is a subtype of $T$ ( $S <: T$ ) iff the set of
numbers in each $S \supset T$. In OO terms $S$ inherits $T$,
as in $S$ may be used in anywhere that $T$ maybe used.

<div style="text-align: center;">

```Mermaid
flowchart LR
    S --> T
```

</div>

If $S \supset T \supset U$, then $U$ inherits $T$ and $T$ inherits $S$.
$U$ inherits $S$ too, in OO this could be defined by the developer,
however with duck typing we assume it is only a transient inheritance
through $T$.

<div style="text-align: center;">

```Mermaid
flowchart LR
    S --> T --> U
```

</div>

If $S \supset T$, $S \supset V$, and $T \nsupseteq V$,
then $S$ inherits both $T$ and $V$ directly.

<div style="text-align: center;">

```Mermaid
flowchart LR
    S --> T
    S --> U
```

</div>

Since $T$ and $U$ may overlap, $V$ might be a supertype of both iff
$(T \cap U) \supset V$.

<div style="text-align: center;">

```Mermaid
flowchart LR
    S --> T --> V
    S --> U --> V
```

</div>
