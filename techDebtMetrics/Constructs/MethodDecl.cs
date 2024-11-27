using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A declaration of a method or function.</summary>
/// <see cref="../../docs/genFeatureDef.md#method"/>
public class MethodDecl : IMethod, IDeclaration, IInitializable {

    /// <summary>The name of the method declaration.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The location the method was defined.</summary>
    public Location Location { get; private set; }

    /// <summary>The The method signature type description.</summary>
    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    /// <summary>Metrics gathered from the method's body, if it had a body.</summary>
    public Metrics? Metrics { get; private set; }

    /// <summary>The package the method was declared in.</summary>
    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    /// <summary>
    /// This the optional object this method is defined as a member of
    /// or null if the method is a function on it's own in the package.
    /// </summary>
    public ObjectDecl? Receiver { get; private set; }

    /// <summary>The type parameters for this method if the method is generic.</summary>
    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    /// <summary>The instances for this declaration if the method is generic.</summary>
    public IReadOnlyList<MethodInst> Instances => this.inInstances.AsReadOnly();
    private readonly List<MethodInst> inInstances = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.inSignature = obj.ReadIndex("signature", project.Signatures);
        this.Metrics = obj.TryReadIndex("metrics", project.Metrics);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        this.Receiver = obj.TryReadIndex("receiver", project.ObjectDecls);
        obj.TryReadIndexList("typeParams", this.inTypeParams, project.TypeParams);
        obj.TryReadIndexList("instances", this.inInstances, project.MethodInsts);
        this.Signature.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.Write(this.Name).Write(this.Signature);
}
