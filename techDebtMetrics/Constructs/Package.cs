using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>
/// A package is typically a folder of code containing zero or more classes (types).
/// A package may have imports, static methods, and global variables too.
/// </summary>
public class Package : IConstruct, IInitializable {
    public string Path { get; private set; } = "";
    public string Name { get; private set; } = "";

    public IReadOnlyList<Package> Imports => this.inImports.AsReadOnly();
    private readonly List<Package> inImports = [];

    public IReadOnlyList<TypeDef> Types => this.inTypes.AsReadOnly();
    private readonly List<TypeDef> inTypes = [];

    public IReadOnlyList<ValueDef> Values => this.inValues.AsReadOnly();
    private readonly List<ValueDef> inValues = [];

    public IReadOnlyList<Method> Methods => this.inMethods.AsReadOnly();
    private readonly List<Method> inMethods = [];

    void IInitializable.Initialize(Project project, Data.Node node) {
        Data.Object obj = node.AsObject();

        this.Path = obj.ReadString("path");
        this.Name = obj.ReadString("name");

        if (obj.Contains("imports")) {        
            Data.Array importArr = obj["imports"].AsArray();
            for (int i = 0; i < importArr.Count; i++) {
                uint pkgIndex = importArr[i].AsUint();
                this.inImports.Add(getter.GetPackageAtIndex(pkgIndex));
            }
        }

        obj.InitializeList(getter, "types", this.inTypes);
        obj.InitializeList(getter, "values", this.inValues);
        obj.InitializeList(getter, "methods", this.inMethods);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write("package ").Write(this.Name).WriteLine(" {");

        Journal j2 = j.Indent;
        j2.Write("path: ").Write(this.Path).WriteLine(";");

        foreach (Package import in this.Imports)
            j2.WriteLine().Write("import ").Write(import.Name).Write(" => ").Write(import.Path).WriteLine(";");

        foreach (TypeDef td in this.Types)
            j2.WriteLine().Write(td).WriteLine(";");

        foreach (ValueDef vd in this.Values)
            j2.WriteLine().Write(vd).WriteLine(";");
        
        foreach (Method m in this.Methods)
            j2.WriteLine().Write(m).WriteLine(";");

        j.Write("}");
    }
}
