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

    public IReadOnlyList<InterfaceDecl> Interfaces => this.inInterfaces.AsReadOnly();
    private readonly List<InterfaceDecl> inInterfaces = [];

    public IReadOnlyList<MethodDecl> Methods => this.inMethods.AsReadOnly();
    private readonly List<MethodDecl> inMethods = [];

    public IReadOnlyList<ObjectDecl> Objects => this.inObjects.AsReadOnly();
    private readonly List<ObjectDecl> inObjects = [];

    public IReadOnlyList<Value> Values => this.inValues.AsReadOnly();
    private readonly List<Value> inValues = [];

    void IInitializable.Initialize(Project project, Node node) {
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
            j2.Write(this.Methods, prefix: "\n\n", suffix: ";", separator: ";\n");
            j2.Write(this.Values, prefix: "\n\n", suffix: ";", separator: ";\n");
            j.Write("\n}");
        }
    }
}
