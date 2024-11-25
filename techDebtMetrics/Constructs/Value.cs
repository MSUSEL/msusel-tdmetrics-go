using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Value : IDeclaration, IInitializable {

    public string Name { get; private set; } = "";

    public Location Location { get; private set; }

    public bool Constant { get; private set; }

    public Metrics? Metrics { get; private set; }

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name      = obj.ReadString("name");
        this.Location  = obj.TryReadLocation("loc", project);
        this.Constant  = obj.TryReadBool("const");
        this.Metrics   = obj.TryReadIndex("metrics", project.Metrics);
        this.inType    = obj.ReadKey<ITypeDesc>("type", project);
        this.inPackage = obj.ReadIndex("package", project.Packages);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) { }
}
