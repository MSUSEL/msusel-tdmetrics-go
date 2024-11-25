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
            Journal j2 = j.Indent.AsLong;
            j2.Write("path: " + this.Path + ";");
            j2.AsShort.Write(this.Imports, prefix: "\nimports: ", suffix: ";");
            j2.Write(this.Interfaces, prefix:"\n", suffix: ";", separator: ";\n");
            j2.Write(this.Methods, prefix: "\n", suffix: ";", separator: ";\n");
            j2.Write(this.Objects, prefix: "\n", suffix: ";", separator: ";\n");
            j2.Write(this.Values, prefix: "\n", suffix: ";", separator: ";\n");
            j.Write("\n}");
        }
    }
}
