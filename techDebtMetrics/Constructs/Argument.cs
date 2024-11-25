using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Argument : IConstruct, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name   = obj.TryReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
    }

    public override string ToString() => this.Name + ":" + this.inType;

    public void ToStub(Journal j) {
        j.Write(this.Name);
        if (!j.Short)
            j.Write(" ").AsShort.Write(this.Type);
    }
}