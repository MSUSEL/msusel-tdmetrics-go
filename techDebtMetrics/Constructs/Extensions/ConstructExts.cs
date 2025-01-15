﻿using System;
using System.Collections.Generic;
using System.Linq;

namespace Constructs.Extensions;

/// <summary>A collection of enumerator extensions.</summary>
static public class ConstructExts {

    /// <summary>This will iterate across all the types reachable from the given construct.</summary>
    /// <remarks>
    /// This will iterate them types depth first and will not interate duplicates.
    /// This will not output the given construct even if it was a type.
    /// </remarks>
    /// <param name="con">The construct to enumerate the types from.</param>
    /// <returns>The enumerator of all the reachable types.</returns>
    public static IEnumerable<ITypeDesc> AllSubTypeDecs(this IConstruct con, bool show = false) {
        HashSet<IConstruct> touched = [];
        Stack<IConstruct> pending = [];
        Stack<string> indents = [];

        //======================================================================== TODO: REMOVE
        if (show) Console.WriteLine("----[ START ]--------------------------------------");
        
        // TODO: Debug why Set<string, int, Map<string, Pointer<int>>> was not being outputted.

        static void print(IConstruct c, string indent) {
            string text = "(" + (c is ITypeDesc) + ", " + (c is TypeParam) + ") " + c.GetType().FullName + "\n" + c;
            Console.WriteLine(indent + "|-- " + text.Replace("\n", "\n"+indent+"|   "));
        }
        if (show) print(con, ">> ");
        //========================================================================

        foreach (IConstruct sub in con.SubConstructs) {
            if (touched.Contains(sub)) continue;
            pending.Push(sub);
            indents.Push(">> |  ");
        }

        while (pending.Count > 0) {
            IConstruct c = pending.Pop();
            string indent = indents.Pop();
            if (!touched.Add(c)) continue;

            //======================================================================== TODO: REMOVE
            if (show) print(c, indent);
            //========================================================================

            if (c is ITypeDesc td) yield return td;

            foreach (IConstruct sub in c.SubConstructs) {
                pending.Push(sub);
                indents.Push(indent + "|  ");
            }
        }
        
        //======================================================================== TODO: REMOVE
        if (show) Console.WriteLine("----[ END ]----------------------------------------");
        //========================================================================
    }

    /// <summary>
    /// Determines if the construct is concrete type without any type parameters,
    /// i.e. not generic and not containing any generic types.
    /// </summary>
    /// <param name="con">The construct to determine if concrete.</param>
    /// <returns>True if concrete, false if genereic or containing a generic type.</returns>
    public static bool IsConcrete(this IConstruct con, bool show = false) =>
        !con.AllSubTypeDecs(show).OfType<TypeParam>().Any();
}
