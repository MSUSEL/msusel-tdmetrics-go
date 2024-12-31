using System.Collections.Generic;
using System.Linq;
using Constructs.Tooling;

namespace Constructs.Extensions;

/// <summary>A collection of enumerator extensions.</summary>
static public class ConstructExts {

    /// <summary>This will iterate across all the types reachable from the given construct.</summary>
    /// <remarks>
    /// This will iterate them types breadth first and will not interate duplicates.
    /// This will not output the given construct even if it was a type.
    /// </remarks>
    /// <param name="con">The construct to enumerate the types from.</param>
    /// <returns>The enumerator of all the reachable types.</returns>
    public static IEnumerable<ITypeDesc> AllSubTypeDecs(this IConstruct con) {
        HashSet<IConstruct> touched = [];
        Queue<IConstruct> pending = [];
        pending.Enqueue(con);
        while (pending.Count > 0) {
            foreach (IConstruct c in pending.Dequeue().SubConstructs) {
                if (touched.Contains(c)) continue;
                if (c is ITypeDesc td) yield return td;
                pending.Enqueue(c);
            }
        }
    }

    /// <summary>
    /// Determines if the construct is concrete type without any type parameters,
    /// i.e. not generic and not containing any generic types.
    /// </summary>
    /// <param name="con">The construct to determine if concrete.</param>
    /// <returns>True if concrete, false if genereic or containing a generic type.</returns>
    public static bool IsConcrete(this IConstruct con) =>
        !con.AllSubTypeDecs().OfType<TypeParam>().Any();
}
