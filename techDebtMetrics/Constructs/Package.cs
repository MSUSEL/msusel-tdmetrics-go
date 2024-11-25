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
    private List<Package> inImports = [];

    public IReadOnlyList<InterfaceDecl> Interfaces => this.inInterfaces.AsReadOnly();
    private List<InterfaceDecl> inInterfaces = [];

    public IReadOnlyList<MethodDecl> Methods => this.inMethods.AsReadOnly();
    private List<MethodDecl> inMethods = [];

    public IReadOnlyList<ObjectDecl> Objects => this.inObjects.AsReadOnly();
    private List<ObjectDecl> inObjects = [];

    public IReadOnlyList<Value> Values => this.inValues.AsReadOnly();
    private List<Value> inValues = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Path = obj.ReadString("path");
        this.Name = obj.ReadString("name");
        this.inImports = obj.TryReadIndexList("imports", project.Packages);
        this.inInterfaces = obj.TryReadIndexList("interfaces", project.InterfaceDecls);
        this.inMethods = obj.TryReadIndexList("methods", project.MethodDecls);
        this.inObjects = obj.TryReadIndexList("objects", project.ObjectDecls);
        this.inValues = obj.TryReadIndexList("values", project.Values);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write("package ").Write(this.Name);
        if (j.Long) {
            j.WriteLine(" {");
            Journal j2 = j.Indent.AsShort;
            j2.Write("path: ").Write(this.Path).WriteLine(";");
            j2.Write("imports: ").Write(this.Imports).WriteLine(";");
            j2.WriteLine(this.Interfaces, prefix: "\n", suffix: "\n", separator: "\n\n");
            j2.WriteLine(this.Methods, prefix: "\n", suffix: "\n", separator: "\n\n");
            j2.WriteLine(this.Objects, prefix: "\n", suffix: "\n", separator: "\n\n");
            j2.WriteLine(this.Values, prefix: "\n", suffix: "\n", separator: "\n\n");
            j.Write("}");
        }
    }
}
