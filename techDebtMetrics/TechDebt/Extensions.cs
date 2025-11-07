using System;
using System.Linq;
using SCG = System.Collections.Generic;

namespace TechDebt;

public static class Extensions {
    public static void ForAll<T>(this SCG.IEnumerable<T> e, Action<T> handle) {
        foreach(T item in e) handle(item);
    }
    
    public static void ForAll<T, R>(this SCG.IEnumerable<T> e, Func<T, R> handle) {
        foreach(T item in e) handle(item);
    }

    public static SCG.IEnumerable<T> WhereNot<T>(this SCG.IEnumerable<T> e, Func<T, bool> predicate) => e.Where(v => !predicate(v));
}
