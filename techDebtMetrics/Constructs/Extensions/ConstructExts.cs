using System.Collections.Generic;

namespace Constructs.Extensions;

static public class ConstructExts {

    /// <summary>This will iterate across all the types reachable from the given construct.</summary>
    /// <remarks>
    /// This will iterate them types breadth first and will not interate duplicates.
    /// This will not output the given construct even if it was a type.
    /// </remarks>
    /// <param name="con">The construct to enumerate the types from.</param>
    /// <returns>The enumerator of all the reachable types.</returns>
    public static IEnumerator<ITypeDesc> AllSubtypes(this IConstruct con) {
        HashSet<IConstruct> touched = [con];
        Queue<IConstruct> pending = [];
        pending.Enqueue(con);
        while (pending.Count > 0) {
            IConstruct cur = pending.Dequeue();
            if (!touched.Add(cur)) continue;
            if (cur is ITypeDesc td) yield return td;
            foreach (IConstruct c in cur.SubConstructs) {
                if (touched.Contains(c)) continue;
                pending.Enqueue(c);
            }
        }
    }
}
