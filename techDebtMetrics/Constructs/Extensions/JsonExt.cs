using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace Constructs.Extensions;

internal static class JsonExt {

    static public void PreallocateList<T>(this JsonObject obj, string name, List<T> list)
        where T : new() {
        int count = obj[name]?.AsArray()?.Count ?? 0;
        for (int i = 0; i < count; i++)
            list[i] = new T();
    }

    static public void InitializeList<T>(this JsonObject obj, TypeGetter getter, string name, List<T> list)
        where T : IInitializable {
        JsonArray? listArr = obj[name]?.AsArray();
        if (listArr is not null) {
            for (int i = 0; i < listArr.Count; i++) {
                JsonNode item = listArr[i] ??
                    throw new MissingDataException(name + "[" + i + "]");
                list[i].Initialize(getter, item);
            }
        }
    }

    static public T ReadValue<T>(this JsonObject obj, string name) {
        JsonNode n = obj[name] ??
            throw new MissingDataException(name);
        return n.GetValue<T>();
    }

    static public T ReadIndexType<T>(this JsonObject obj, string name, TypeGetter getter)
        where T : ITypeDesc {
        uint typeIndex = obj[name]?.GetValue<uint>() ??
            throw new MissingDataException(name);
        return getter.GetTypeAtIndex<T>(typeIndex);
    }

    static public void ReadIndexTypeList<T>(this JsonObject obj, string name, TypeGetter getter, List<T> list)
        where T : ITypeDesc {
        JsonArray? exactArr = obj[name]?.AsArray();
        if (exactArr is not null) {
            for (int i = 0; i < exactArr.Count; i++) {
                uint typeIndex = exactArr[i]?.GetValue<uint>() ??
                    throw new MissingDataException(name + "[" + i + "]");
                list.Add(getter.GetTypeAtIndex<T>(typeIndex));
            }
        }
    }
}
