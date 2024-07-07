using System.Collections.Generic;
using System.Linq;

namespace Constructs.Extensions;

internal static class GeneralExt {

    public static string Join(this IEnumerable<string> source, string separator = ", ") =>
        string.Join(separator, source);

    public static IEnumerable<string> ToStrings<T>(this IEnumerable<T> source, string prefix = "", string suffix = "", string onNull = "<null>") =>
        source.Select(x => prefix + (x?.ToString() ?? onNull) + suffix);
}
