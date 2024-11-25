using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;
using System.IO;

namespace Constructs;

/// <summary>
/// A project represents all the packages to completely describe a program or a library.
/// </summary>
public class Project : IConstruct {
    public readonly string Language;
    public readonly Locations Locations;

    public IReadOnlyList<Abstract> Abstracts => this.inAbstracts.AsReadOnly();
    private readonly List<Abstract> inAbstracts = [];

    public IReadOnlyList<Argument> Arguments => this.inArguments.AsReadOnly();
    private readonly List<Argument> inArguments = [];

    public IReadOnlyList<Basic> Basics => this.inBasics.AsReadOnly();
    private readonly List<Basic> inBasics = [];

    public IReadOnlyList<Field> Fields => this.inFields.AsReadOnly();
    private readonly List<Field> inFields = [];

    public IReadOnlyList<InterfaceDecl> InterfaceDecls => this.inInterfaceDecls.AsReadOnly();
    private readonly List<InterfaceDecl> inInterfaceDecls = [];

    public IReadOnlyList<InterfaceDesc> InterfaceDescs => this.inInterfaceDescs.AsReadOnly();
    private readonly List<InterfaceDesc> inInterfaceDescs = [];
    
    public IReadOnlyList<InterfaceInst> InterfaceInsts => this.inInterfaceInsts.AsReadOnly();
    private readonly List<InterfaceInst> inInterfaceInsts = [];
    
    public IReadOnlyList<MethodDecl> MethodDecls => this.inMethodDecls.AsReadOnly();
    private readonly List<MethodDecl> inMethodDecls = [];

    public IReadOnlyList<MethodInst> MethodInsts => this.inMethodInsts.AsReadOnly();
    private readonly List<MethodInst> inMethodInsts = [];

    public IReadOnlyList<Metrics> Metrics => this.inMetrics.AsReadOnly();
    private readonly List<Metrics> inMetrics = [];

    public IReadOnlyList<ObjectDecl> ObjectDecls => this.inObjectDecls.AsReadOnly();
    private readonly List<ObjectDecl> inObjectDecls = [];

    public IReadOnlyList<ObjectInst> ObjectInsts => this.inObjectInsts.AsReadOnly();
    private readonly List<ObjectInst> inObjectInsts = [];

    public IReadOnlyList<Package> Packages => this.inPackages.AsReadOnly();
    private readonly List<Package> inPackages = [];

    public IReadOnlyList<Selection> Selections => this.inSelections.AsReadOnly();
    private readonly List<Selection> inSelections = [];

    public IReadOnlyList<Signature> Signatures => this.inSignatures.AsReadOnly();
    private readonly List<Signature> inSignatures = [];

    public IReadOnlyList<StructDesc> StructDescs => this.inStructDescs.AsReadOnly();
    private readonly List<StructDesc> inStructDescs = [];

    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    public IReadOnlyList<Value> Values => this.inValues.AsReadOnly();
    private readonly List<Value> inValues = [];

    /// <summary>Loads a project from a YAML file.</summary>
    /// <param name="path">The file path to read from.</param>
    /// <returns>The project that was read from the YAML.</returns>
    static public Project FromFile(string path) {
        using StreamReader r = new(path);
        return FromText(r.ReadToEnd());
    }

    /// <summary>Loads a project from a string containing YAML data.</summary>
    /// <param name="text">The string containing YAML data to read from.</param>
    /// <returns>The project that was read from the YAML.</returns>
    static public Project FromText(string text) =>
        new(Node.Parse(text));

    /// <summary>Creates a new project.</summary>
    /// <param name="root">The YAML root node to load.</param>
    internal Project(Node root) {
        Object obj = root.AsObject();
        this.Language = obj.ReadString("language");
        this.Locations = new(obj.TryReadNode("locs"));

        obj.PreallocateList("abstracts", this.inAbstracts);
        obj.PreallocateList("arguments", this.inArguments);
        obj.PreallocateList("basics", this.inBasics);
        obj.PreallocateList("fields", this.inFields);
        obj.PreallocateList("interfaceDecls", this.inInterfaceDecls);
        obj.PreallocateList("interfaceDescs", this.inInterfaceDescs);
        obj.PreallocateList("interfaceInsts", this.inInterfaceInsts);
        obj.PreallocateList("methods", this.inMethodDecls);
        obj.PreallocateList("methodInsts", this.inMethodInsts);
        obj.PreallocateList("metrics", this.inMetrics);
        obj.PreallocateList("objects", this.inObjectDecls);
        obj.PreallocateList("objectInsts", this.inObjectInsts);
        obj.PreallocateList("packages", this.inPackages);
        obj.PreallocateList("selections", this.inSelections);
        obj.PreallocateList("signatures", this.inSignatures);
        obj.PreallocateList("structDescs", this.inStructDescs);
        obj.PreallocateList("typeParams", this.inTypeParams);
        obj.PreallocateList("values", this.inValues);

        obj.InitializeList(this, "abstracts", this.inAbstracts);
        obj.InitializeList(this, "arguments", this.inArguments);
        obj.InitializeList(this, "basics", this.inBasics);
        obj.InitializeList(this, "fields", this.inFields);
        obj.InitializeList(this, "interfaceDecls", this.inInterfaceDecls);
        obj.InitializeList(this, "interfaceDescs", this.inInterfaceDescs);
        obj.InitializeList(this, "interfaceInsts", this.inInterfaceInsts);
        obj.InitializeList(this, "methods", this.inMethodDecls);
        obj.InitializeList(this, "methodInsts", this.inMethodInsts);
        obj.InitializeList(this, "metrics", this.inMetrics);
        obj.InitializeList(this, "objects", this.inObjectDecls);
        obj.InitializeList(this, "objectInsts", this.inObjectInsts);
        obj.InitializeList(this, "packages", this.inPackages);
        obj.InitializeList(this, "selections", this.inSelections);
        obj.InitializeList(this, "signatures", this.inSignatures);
        obj.InitializeList(this, "structDescs", this.inStructDescs);
        obj.InitializeList(this, "typeParams", this.inTypeParams);
        obj.InitializeList(this, "values", this.inValues);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsLong.Write(this.Packages, separator: "\n\n");
}
