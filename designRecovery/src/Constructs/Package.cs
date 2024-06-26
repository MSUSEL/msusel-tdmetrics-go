﻿using designRecovery.src.Extensions;
using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Package : IInitializer {
    public string Path { get; private set; } = "";
    public string Name { get; private set; } = "";

    private readonly List<Package> inImport = [];
    public IReadOnlyList<Package> Imports => this.inImport.AsReadOnly();

    private readonly List<TypeDef> inTypes = [];
    public IReadOnlyList<TypeDef> Types => this.inTypes.AsReadOnly();

    private readonly List<ValueDef> inValues = [];
    public IReadOnlyList<ValueDef> Values => this.inValues.AsReadOnly();

    private readonly List<Method> inMethods = [];
    public IReadOnlyList<Method> Methods => this.inMethods.AsReadOnly();

    static private void preallocate<T>(JsonObject obj, string name, List<T> list) where T: new() {
        int count = obj[name]?.AsArray()?.Count ?? 0;
        for (int i = 0; i<count; i++) list[i] = new T();
    }

    public Package(JsonNode node) {
        JsonObject obj = node.AsObject();
        preallocate(obj, "types", this.inTypes);
        preallocate(obj, "values", this.inValues);
        preallocate(obj, "methods", this.inMethods);
    }

    static private void initializeList<T>(JsonObject obj, TypeGetter getter, string name, List<T> list)
        where T : IInitializer {
        JsonArray? listArr = obj[name]?.AsArray();
        if (listArr is not null) {
            for (int i = 0; i < listArr.Count; i++) {
                JsonNode item = listArr[i] ??
                    throw new MissingDataException(name+"["+i+"]");
                list[i].Initialize(getter, item);
            }
        }
    }

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();

        this.Path = obj.ReadValue<string>("path");
        this.Name = obj.ReadValue<string>("name");

        JsonArray? importArr = obj["imports"]?.AsArray();
        if (importArr is not null) {
            for (int i = 0; i < importArr.Count; i++) {
                uint pkgIndex = importArr[i]?.GetValue<uint>() ??
                    throw new MissingDataException("import["+i+"]");
                inImport.Add(getter.GetPackageAtIndex(pkgIndex));
            }
        }

        initializeList(obj, getter, "types", this.inTypes);
        initializeList(obj, getter, "values", this.inValues);
        initializeList(obj, getter, "methods", this.inMethods);
    }
}
