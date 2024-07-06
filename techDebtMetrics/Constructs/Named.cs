using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Named : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");
    private ITypeDesc? inType;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }

    public override string ToString() => this.Name + ":" + this.Type;

    public string ToStub() => this.Name + " " + this.Type.ToStub();
}
