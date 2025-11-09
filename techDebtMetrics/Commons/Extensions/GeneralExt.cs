using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;

namespace Commons.Extensions;

/// <summary>A collection of enumerator extensions.</summary>
public static class GeneralExt {

    /// <summary>Joins the given values together into a string.</summary>
    /// <typeparam name="T">The type of values to join.</typeparam>
    /// <param name="source">The source of the values to join as strings.</param>
    /// <param name="separator">The separator to put between the strings.</param>
    /// <returns>The string of the joined values.</returns>
    public static string Join<T>(this IEnumerable<T> source, string separator = ", ") =>
        string.Join(separator, source);

    /// <summary>Converts the enumerator of values into an enumerator of strings.</summary>
    /// <typeparam name="T">The type of values to create strings for.</typeparam>
    /// <param name="source">The source of the values to convert to strings.</param>
    /// <param name="prefix">The prefix to add to the front of each item when creating a string.</param>
    /// <param name="suffix">The suffix to add to the end of each item when creating a string.</param>
    /// <param name="onNull">The string to use in-place of the item when the item is null.</param>
    /// <returns>The enumerator of strings for the given values.</returns>
    public static IEnumerable<string> ToStrings<T>(this IEnumerable<T> source, string prefix = "", string suffix = "", string onNull = "<null>") =>
        source.Select(x => prefix + (x?.ToString() ?? onNull) + suffix);

    /// <summary>Performs the given function on all the values in the enumerable.</summary>
    /// <typeparam name="T">The type of the values.</typeparam>
    /// <param name="e">The enumerable to apply the given function to.</param>
    /// <param name="handle">The function to run with all of the values.</param>
    public static void ForAll<T>(this IEnumerable<T> e, Action<T> handle) {
        foreach (T item in e) handle(item);
    }

    /// <summary>Performs the given function on all the values in the enumerable.</summary>
    /// <typeparam name="T">The type of the values.</typeparam>
    /// <typeparam name="R">The type of the return type to ignore.</typeparam>
    /// <param name="e">The enumerable to apply the given function to.</param>
    /// <param name="handle">
    /// The function to run with all of the values.
    /// The function may have a return type that will be ignored.
    /// </param>
    public static void ForAll<T, R>(this IEnumerable<T> e, Func<T, R> handle) {
        foreach (T item in e) handle(item);
    }

    /// <summary>Filters the enumerable values with the given predicate.</summary>
    /// <typeparam name="T">The type of the values.</typeparam>
    /// <param name="e">The enumberable to filter.</param>
    /// <param name="predicate">
    /// The predicate to filter with.
    /// Any value that this returns false for will be emitted into the output enumeration.
    /// </param>
    /// <returns>The filtered enumberator.</returns>
    public static IEnumerable<T> WhereNot<T>(this IEnumerable<T> e, Func<T, bool> predicate) => e.Where(v => !predicate(v));

    /// <summary>BinarySearch of a sorted list.</summary>
    /// <typeparam name="T">The type of values in the list.</typeparam>
    /// <param name="list">The list to perform the search in.</param>
    /// <param name="value">The value to search for in the list.</param>
    /// <returns>
    /// The index of the location to insert the target location
    /// and if the target was found (true) or not (false).
    /// </returns>
    public static (int, bool) BinarySearch<T>(this IList<T> list, T value)
        where T : IComparable<T> =>
        list.BinarySearch(value.CompareTo);

    /// <summary>BinarySearch of a sorted list.</summary>
    /// <typeparam name="T">The type of values in the list.</typeparam>
    /// <param name="list">The list to perform the search in.</param>
    /// <param name="comparer">
    /// The comparer used to find a value via binary search.
    /// The list must be sorted in the same order as this comparer would sort te list.
    /// The comparer should return:
    /// - negative if the target value is less than the given value.
    /// - positive if the target value is higher than the given value.
    /// - zero if the target value matches the given value.
    /// </param>
    /// <returns>
    /// The index of the location to insert the target location
    /// and if the target was found (true) or not (false).
    /// </returns>
    public static (int, bool) BinarySearch<T>(this IList<T> list, Func<T, int> comparer) {
        int low = 0, high = list.Count - 1, floor = -1, mid, cmp;
        while (low <= high) {
            mid = low + (high - low) / 2;
            cmp = comparer(list[mid]); 
            if (cmp == 0) return (mid, true);
            if (cmp > 0) {
                floor = mid; 
                low = mid + 1; 
            }
            else high = mid - 1; 
        }
        return (floor, false);
    }
}
