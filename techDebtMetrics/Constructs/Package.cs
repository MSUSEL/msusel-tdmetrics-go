using Commons.Data.Reader;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;

namespace Constructs;

/// <summary>
/// A package is typically a folder of code containing zero or more classes (types).
/// A package may have imports, static methods, and global variables too.
/// </summary>
/// <see cref="../../docs/genFeatureDef.md#package"/>
public class Package : IConstruct, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The import path for this package.</summary>
    public string Path { get; private set; } = "";

    /// <summary>The name of the package.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The list of packages that this package imports.</summary>
    public IReadOnlyList<Package> Imports => this.inImports.AsReadOnly();
    private readonly List<Package> inImports = [];

    /// <summary>The list of interfaces declared in this package.</summary>
    public IReadOnlyList<InterfaceDecl> Interfaces => this.inInterfaces.AsReadOnly();
    private readonly List<InterfaceDecl> inInterfaces = [];

    /// <summary>The list of methods declared in this package.</summary>
    public IReadOnlyList<MethodDecl> Methods => this.inMethods.AsReadOnly();
    private readonly List<MethodDecl> inMethods = [];

    /// <summary>The list of methods declared in this package and not part of an object.</summary>
    public IReadOnlyList<MethodDecl> StaticMethods =>
        [.. this.Methods.Where(m => m.Receiver is null)];

    /// <summary>The list of objects declared in this package.</summary>
    public IReadOnlyList<ObjectDecl> Objects => this.inObjects.AsReadOnly();
    private readonly List<ObjectDecl> inObjects = [];

    /// <summary>The package level values and constants declared in this package.</summary>
    public IReadOnlyList<Value> Values => this.inValues.AsReadOnly();
    private readonly List<Value> inValues = [];

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.Imports) yield return c;
            foreach (IConstruct c in this.Interfaces) yield return c;
            foreach (IConstruct c in this.Methods) yield return c;
            foreach (IConstruct c in this.StaticMethods) yield return c;
            foreach (IConstruct c in this.Objects) yield return c;
            foreach (IConstruct c in this.Values) yield return c;
        }
    }

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Path = obj.ReadString("path");
        this.Name = obj.ReadString("name");
        obj.TryReadIndexList("imports", this.inImports, project.Packages);
        obj.TryReadIndexList("interfaces", this.inInterfaces, project.InterfaceDecls);
        obj.TryReadIndexList("methods", this.inMethods, project.MethodDecls);
        obj.TryReadIndexList("objects", this.inObjects, project.ObjectDecls);
        obj.TryReadIndexList("values", this.inValues, project.Values);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write("package ").Write(this.Name);
        if (j.Long) {
            j.WriteLine(" {");
            Journal j2 = j.Indent.AsLong;
            j2.Write(this.Path, prefix: "path: ", suffix: ";");
            j2.AsShort.Write(this.Imports, prefix: "\nimports: ", suffix: ";");
            j2.Write(this.Interfaces, prefix: "\n\n", separator: "\n\n");
            j2.Write(this.Objects, prefix: "\n\n", separator: "\n\n");
            j2.Write(this.StaticMethods, prefix: "\n\n", suffix: ";", separator: ";\n");
            j2.Write(this.Values, prefix: "\n\n", suffix: ";", separator: ";\n");
            j.Write("\n}");
        }
    }
}
