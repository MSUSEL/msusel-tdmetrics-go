using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace Constructs;

/// <summary>
/// A package is typically a folder of code containing zero or more classes (types).
/// A package may have imports, static methods, and global variables too.
/// </summary>
public class Package : IConstruct, IInitializable {
    public string Path { get; private set; } = "";
    public string Name { get; private set; } = "";

    public IReadOnlyList<Package> Imports => this.inImport.AsReadOnly();
    private readonly List<Package> inImport = [];

    public IReadOnlyList<TypeDef> Types => this.inTypes.AsReadOnly();
    private readonly List<TypeDef> inTypes = [];

    public IReadOnlyList<ValueDef> Values => this.inValues.AsReadOnly();
    private readonly List<ValueDef> inValues = [];

    public IReadOnlyList<Method> Methods => this.inMethods.AsReadOnly();
    private readonly List<Method> inMethods = [];

    public Package(JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.PreallocateList("types", this.inTypes);
        obj.PreallocateList("values", this.inValues);
        obj.PreallocateList("methods", this.inMethods);
    }

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();

        this.Path = obj.ReadValue<string>("path");
        this.Name = obj.ReadValue<string>("name");

        JsonArray? importArr = obj["imports"]?.AsArray();
        if (importArr is not null) {
            for (int i = 0; i < importArr.Count; i++) {
                uint pkgIndex = importArr[i]?.GetValue<uint>() ??
                    throw new MissingDataException("import[" + i + "]");
                this.inImport.Add(getter.GetPackageAtIndex(pkgIndex));
            }
        }

        obj.InitializeList(getter, "types", this.inTypes);
        obj.InitializeList(getter, "values", this.inValues);
        obj.InitializeList(getter, "methods", this.inMethods);
    }

    public string ToStub() => throw new System.NotImplementedException(); // TODO: Implement
}
