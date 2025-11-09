using Commons.Data.Locations;
using Commons.Data.Yaml;
using Constructs.Tooling;
using System;
using System.Collections.Generic;
using System.IO;

namespace Constructs;

/// <summary>A project represents all the packages to completely describe a program or a library.</summary>
/// <see cref="../../docs/genFeatureDef.md#project"/>
public class Project : IConstruct, IKeyResolver {

    /// <summary>Returns 1 since there is only one project.</summary>
    public int Index => 1;

    /// <summary>The name of the langauge that this project was created from.</summary>
    public readonly string Language;

    /// <summary>The file paths and line numbers for the source code that created this package.</summary>
    public readonly Reader Locations;

    /// <summary>The collection of all abstracts in this project.</summary>
    public IReadOnlyList<Abstract> Abstracts => this.inAbstracts.AsReadOnly();
    private readonly List<Abstract> inAbstracts = [];

    /// <summary>The collection of all arguments in this project.</summary>
    public IReadOnlyList<Argument> Arguments => this.inArguments.AsReadOnly();
    private readonly List<Argument> inArguments = [];

    /// <summary>The collection of all basic in this project.</summary>
    public IReadOnlyList<Basic> Basics => this.inBasics.AsReadOnly();
    private readonly List<Basic> inBasics = [];

    /// <summary>The collection of all fields in this project.</summary>
    public IReadOnlyList<Field> Fields => this.inFields.AsReadOnly();
    private readonly List<Field> inFields = [];

    /// <summary>The collection of all interface declarations in this project.</summary>
    public IReadOnlyList<InterfaceDecl> InterfaceDecls => this.inInterfaceDecls.AsReadOnly();
    private readonly List<InterfaceDecl> inInterfaceDecls = [];

    /// <summary>The collection of all interface type descriptions in this project.</summary>
    public IReadOnlyList<InterfaceDesc> InterfaceDescs => this.inInterfaceDescs.AsReadOnly();
    private readonly List<InterfaceDesc> inInterfaceDescs = [];

    /// <summary>The collection of all interface instances in this project.</summary>
    public IReadOnlyList<InterfaceInst> InterfaceInsts => this.inInterfaceInsts.AsReadOnly();
    private readonly List<InterfaceInst> inInterfaceInsts = [];

    /// <summary>The collection of all method declarations in this project.</summary>
    public IReadOnlyList<MethodDecl> MethodDecls => this.inMethodDecls.AsReadOnly();
    private readonly List<MethodDecl> inMethodDecls = [];

    /// <summary>The collection of all method instances in this project.</summary>
    public IReadOnlyList<MethodInst> MethodInsts => this.inMethodInsts.AsReadOnly();
    private readonly List<MethodInst> inMethodInsts = [];

    /// <summary>The collection of all metrics in this project.</summary>
    public IReadOnlyList<Metrics> Metrics => this.inMetrics.AsReadOnly();
    private readonly List<Metrics> inMetrics = [];

    /// <summary>The collection of all metrics in this project.</summary>
    public IReadOnlyList<ObjectDecl> ObjectDecls => this.inObjectDecls.AsReadOnly();
    private readonly List<ObjectDecl> inObjectDecls = [];

    /// <summary>The collection of all object instances in this project.</summary>
    public IReadOnlyList<ObjectInst> ObjectInsts => this.inObjectInsts.AsReadOnly();
    private readonly List<ObjectInst> inObjectInsts = [];

    /// <summary>The collection of all packages in this project.</summary>
    public IReadOnlyList<Package> Packages => this.inPackages.AsReadOnly();
    private readonly List<Package> inPackages = [];

    /// <summary>The collection of all selections in this project.</summary>
    public IReadOnlyList<Selection> Selections => this.inSelections.AsReadOnly();
    private readonly List<Selection> inSelections = [];

    /// <summary>The collection of all signature type descriptions in this project.</summary>
    public IReadOnlyList<Signature> Signatures => this.inSignatures.AsReadOnly();
    private readonly List<Signature> inSignatures = [];

    /// <summary>The collection of all structure type descriptions in this project.</summary>
    public IReadOnlyList<StructDesc> StructDescs => this.inStructDescs.AsReadOnly();
    private readonly List<StructDesc> inStructDescs = [];

    /// <summary>The collection of all type parameters in this project.</summary>
    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    /// <summary>The collection of all values and constants in this project.</summary>
    public IReadOnlyList<Value> Values => this.inValues.AsReadOnly();
    private readonly List<Value> inValues = [];

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.Abstracts) yield return c;
            foreach (IConstruct c in this.Arguments) yield return c;
            foreach (IConstruct c in this.Basics) yield return c;
            foreach (IConstruct c in this.Fields) yield return c;
            foreach (IConstruct c in this.InterfaceDecls) yield return c;
            foreach (IConstruct c in this.InterfaceDescs) yield return c;
            foreach (IConstruct c in this.InterfaceInsts) yield return c;
            foreach (IConstruct c in this.MethodDecls) yield return c;
            foreach (IConstruct c in this.MethodInsts) yield return c;
            foreach (IConstruct c in this.Metrics) yield return c;
            foreach (IConstruct c in this.ObjectDecls) yield return c;
            foreach (IConstruct c in this.ObjectInsts) yield return c;
            foreach (IConstruct c in this.Packages) yield return c;
            foreach (IConstruct c in this.Selections) yield return c;
            foreach (IConstruct c in this.Signatures) yield return c;
            foreach (IConstruct c in this.StructDescs) yield return c;
            foreach (IConstruct c in this.TypeParams) yield return c;
            foreach (IConstruct c in this.Values) yield return c;
        }
    }

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
        Commons.Data.Yaml.Object obj = root.AsObject();
        this.Language = obj.ReadString("language");
        this.Locations = Reader.Read(obj.TryReadNode("locs"));

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

    /// <summary>Reads the given index from the given source as part of reading the given key.</summary>
    /// <typeparam name="T">The type of the list to read from.</typeparam>
    /// <param name="key">The key that is being processed.</param>
    /// <param name="index">The index from the key used to read a value from the given list.</param>
    /// <param name="source">The list get an item at the given index from.</param>
    /// <returns>The item from the given list at the given index.</returns>
    static private T readKeyIndex<T>(string key, int index, IReadOnlyList<T> source) {
        if (index < 0 || index >= source.Count)
            throw new Exception("Key " + key + " out of range [0.." + source.Count + "): " + index);
        return source[index];
    }

    /// <summary>Reads a single key from the given project.</summary>
    /// <see cref="docs/genFeatureDef.md#keys"/>
    /// <param name="key">The key of the value to read.</param>
    /// <param name="project">The project to read a key from.</param>
    /// <returns>The read key from the project.</returns>
    public object FindData(string name, int index) {
        return name switch {
            "abstract" => readKeyIndex(name, index, this.Abstracts),
            "argument" => readKeyIndex(name, index, this.Arguments),
            "basic" => readKeyIndex(name, index, this.Basics),
            "field" => readKeyIndex(name, index, this.Fields),
            "interfaceDecl" => readKeyIndex(name, index, this.InterfaceDecls),
            "interfaceDesc" => readKeyIndex(name, index, this.InterfaceDescs),
            "interfaceInst" => readKeyIndex(name, index, this.InterfaceInsts),
            "method" => readKeyIndex(name, index, this.MethodDecls),
            "methodInst" => readKeyIndex(name, index, this.MethodInsts),
            "metrics" => readKeyIndex(name, index, this.Metrics),
            "object" => readKeyIndex(name, index, this.ObjectDecls),
            "objectInst" => readKeyIndex(name, index, this.ObjectInsts),
            "package" => readKeyIndex(name, index, this.Packages),
            "selection" => readKeyIndex(name, index, this.Selections),
            "signature" => readKeyIndex(name, index, this.Signatures),
            "structDesc" => readKeyIndex(name, index, this.StructDescs),
            "typeParam" => readKeyIndex(name, index, this.TypeParams),
            "value" => readKeyIndex(name, index, this.Values),
            _ => throw new InvalidDataException(name)
        };
    }
}
