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

## Tree of Inheritance

Assuming a type can be represented by a set of numbers,
the type $S$ is a subtype of $T$ ( $S <: T$ ) iff the set of
numbers in each $S \supset T$. In OO terms $S$ inherits $T$,
as in $S$ may be used in anywhere that $T$ maybe used.

```Mermaid
flowchart LR
    S --> T
```

If $S \supset T \supset U$, then $U$ inherits $T$ and $T$ inherits $S$.
$U$ inherits $S$ too, in OO this could be defined by the developer,
however with duck typing we assume it is only a transient inheritance
through $T$.

```Mermaid
flowchart LR
    S --> T --> U
```

If $S \supset T$, $S \supset V$, and $T \nsupseteq V$,
then $S$ inherits both $T$ and $V$ directly.

```Mermaid
flowchart LR
    S --> T
    S --> V
```

Since $T$ and $V$ may overlap, $U$ might be a supertype of both iff
$(T \cap V) \supseteq U$.

```Mermaid
flowchart LR
    S --> T --> U
    S --> V --> U
```

## Insert to Tree

1. $N$ is a node, i.e. type, in the tree, $\mathbb{T}$:
    1. $N$ has zero or more parent nodes in the parent set $Np$.
    2. The number of parents is denoted as $|Np|$.
    3. The $i^{th}$ parent is denoted as $N_i$
       where $1 \le i \le |Np|$ and $N_i \in Np$.
    4. The order of the parents doesn't matter.
    5. $\forall N_i \mid N \supset N_i$ meaning all parents of $N$
       are supertypes of $N$.
    6. $\forall \left( N_i, N_j \right) \mid i \ne j, N_i \nsupseteq N_j$
       meaning that no parent of $N$ is a supertype of any other parent of $N$.
    7. $N$ is unique such that any other node is either a subtype or supertype
       of $N$ but not equal to $N$.

2. $R$ is a node that is the root node in the tree:
    1. $R$ is an imaginary type that is considered a subtype of all types
       such that any parent of $R$ are not a supertype to any other type.
    2. $R$ is used to make a forest of inheritance into a single tree.
    3. When inserting a new node into the tree, the insertion starts
       comparing against $R_i \in Rp$.

3. Inserting a node, $X$, into a node, $Y$:
    1. Initial state: $|Xp| = 0$ and $Y \supset X$.
    2. Assign $A = \left\{ Y_i \mid Y_i \supset X \right\}$.
       If $|A| > 0$, then remove $A$ from $Yp$, add $A$ to $Xp$, and
       add $X$ as a parent of $Y$.
    3. For all $Y_i \supset X$, insert $X$ into $Y_i$.
    4. For all $\left( Y_i \cap X \right) \ne \emptyset$ but not used in
       prior steps, check the subtree for any node that is a supertype.
       Only follow branches that there is still an overlap and that overlap
       isn't already supertype of a parent already in $X$.
