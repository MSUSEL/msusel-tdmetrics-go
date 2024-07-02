using designRecovery.src.Constructs;
using System.Text.Json.Nodes;

namespace designRecovery.src.Extensions;

internal static class Extensions {
    public static string Join(this IEnumerable<string> source, string separator = ", ") =>
        string.Join(separator, source);
    
    public static IEnumerable<string> ToStrings<T>(this IEnumerable<T> source, string prefix = "", string suffix = "", string onNull = "<null>") =>
        source.Select(x => prefix+(x?.ToString() ?? onNull)+suffix);

    static public T ReadValue<T>(this JsonObject obj, string name) {
        JsonNode n = obj[name] ?? throw new MissingDataException(name);
        return n.GetValue<T>();
    }

    static public T ReadIndexType<T>(this JsonObject obj, string name, TypeGetter getter)
        where T: ITypeDesc {
        uint typeIndex = obj[name]?.GetValue<uint>() ??
            throw new MissingDataException(name);
        return getter.GetTypeAtIndex<T>(typeIndex);
    }
    
    static public void ReadIndexTypeList<T>(this JsonObject obj, string name, TypeGetter getter, List<T> list)
        where T: ITypeDesc {
        JsonArray? exactArr = obj[name]?.AsArray();
        if (exactArr is not null) {
            for (int i = 0; i < exactArr.Count; i++) {
                uint typeIndex = exactArr[i]?.GetValue<uint>() ??
                    throw new MissingDataException(name+"["+i+"]");
                list.Add(getter.GetTypeAtIndex<T>(typeIndex));
            }
        }
    }
}
