using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class MethodDecl : IMethod, IDeclaration, IInitializable {
    public string Name { get; private set; } = "";

    public Location Location { get; private set; }

    public IReadOnlyList<MethodInst> Instances => this.inInstances.AsReadOnly();
    private readonly List<MethodInst> inInstances = [];

    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    public Metrics? Metrics { get; private set; }

    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    public ObjectDecl? Receiver { get; private set; }

    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.inSignature = obj.ReadIndex("signature", project.Signatures);
        this.Metrics = obj.TryReadIndex("metrics", project.Metrics);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        this.Receiver = obj.TryReadIndex("receiver", project.ObjectDecls);
        obj.TryReadIndexList("instances", this.inInstances, project.MethodInsts);
        obj.TryReadIndexList("typeParams", this.inTypeParams, project.TypeParams);
        this.Signature.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.Write(this.Name).Write(this.Signature);
}
