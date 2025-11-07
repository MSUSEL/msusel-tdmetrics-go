using System.Collections.Generic;
using System.Linq;

namespace Commons.Extensions;

/// <summary>A collection of enumerator extensions.</summary>
internal static class GeneralExt {

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
}
