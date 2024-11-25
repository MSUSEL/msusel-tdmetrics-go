using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class MethodDecl : IDeclaration, IInitializable {
    public string Name { get; private set; } = "";

    public Location Location { get; private set; }

    public IReadOnlyList<MethodInst> Instances => this.inInstances.AsReadOnly();
    private List<MethodInst> inInstances = [];

    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    public Metrics Metrics => this.inMetrics ??
        throw new UninitializedException(nameof(this.Metrics));
    private Metrics? inMetrics;

    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    public ObjectDecl Receiver => this.inReceiver ??
        throw new UninitializedException(nameof(this.Receiver));
    private ObjectDecl? inReceiver;

    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private List<TypeParam> inTypeParams = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name         = obj.ReadString("name");
        this.Location     = obj.TryReadLocation("loc", project);
        this.inInstances  = obj.TryReadIndexList("instances", project.MethodInsts);
        this.inSignature  = obj.ReadIndex("signature", project.Signatures);
        this.inMetrics    = obj.ReadIndex("metrics", project.Metrics);
        this.inPackage    = obj.ReadIndex("package", project.Packages);
        this.inReceiver   = obj.ReadIndex("receiver", project.ObjectDecls);
        this.inTypeParams = obj.TryReadIndexList("typeParams", project.TypeParams);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.Write(this.Name).AsShort.Write(this.Signature);
}
