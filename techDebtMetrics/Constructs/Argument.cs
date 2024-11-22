using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Field : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    public TypeDef Type => this.inType ??
        throw new UninitializedException("type");
    private TypeDef? inType;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inType = obj.ReadIndexType<TypeDef>("type", getter);
    }

    public override string ToString() => this.Name + ":" + this.inType;

    public void ToStub(Journal j) {
        j.Write(this.Name);
        if (!j.Short)
            j.Write(" ").AsShort.Write(this.Type);
    }
}
