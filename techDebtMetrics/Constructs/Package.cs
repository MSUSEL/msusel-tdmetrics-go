using Constructs.Tooling;
using System.Collections.Generic;
using System.Text;

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

    internal Package(Data.Node node) {
        Data.Object obj = node.AsObject();
        obj.PreallocateList("types", this.inTypes);
        obj.PreallocateList("values", this.inValues);
        obj.PreallocateList("methods", this.inMethods);
    }

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
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

    public string ToStub() {
        StringBuilder sb = new();
        sb.Append("package ");
        sb.Append(this.Name);
        sb.AppendLine(" {");

        sb.Append("   path: ");
        sb.Append(this.Path);
        sb.AppendLine(";");

        foreach (Package import in this.Imports) {
            sb.Append("   import: ");
            sb.Append(import.Name);
            sb.Append(" => ");
            sb.Append(import.Path);
            sb.AppendLine(";");
        }

        foreach (TypeDef td in this.Types) {
            sb.Append("   ");
            sb.Append(td.ToStub().Replace("\n", "\n   "));
            sb.AppendLine(";");
        }

        foreach (ValueDef vd in this.Values) {
            sb.Append("   ");
            sb.Append(vd.ToStub().Replace("\n", "\n   "));
            sb.AppendLine(";");
        }

        foreach (Method m in this.Methods) {
            sb.Append("   ");
            sb.Append(m.ToStub().Replace("\n", "\n   "));
            sb.AppendLine(";");
        }

        sb.Append('}');
        return sb.ToString();
    }
}
