namespace designRecovery.src.Extensions;

internal static class Extensions {
    public static string Join(this IEnumerable<string> source, string separator = ", ") =>
        string.Join(separator, source);
    
    public static IEnumerable<string> ToStrings<T>(this IEnumerable<T> source, string prefix = "", string suffix = "", string onNull = "<null>") =>
        source.Select(x => prefix+(x?.ToString() ?? onNull)+suffix);
}
